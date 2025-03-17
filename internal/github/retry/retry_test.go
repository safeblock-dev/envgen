package retry_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/safeblock-dev/envgen/internal/github/retry"
)

// MockTransport implements http.RoundTripper for testing.
type MockTransport struct {
	responses      []*http.Response
	errors         []error
	requestsCount  int
	requestsCalled int
}

func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	m.requestsCalled++
	m.requestsCount++

	if len(m.responses) == 0 || len(m.errors) == 0 {
		return &http.Response{StatusCode: http.StatusInternalServerError}, nil
	}

	if m.requestsCalled > len(m.responses) || m.requestsCalled > len(m.errors) {
		// Return the last configured response
		lastResponse := m.responses[len(m.responses)-1]
		lastError := m.errors[len(m.errors)-1]

		return lastResponse, lastError
	}

	return m.responses[m.requestsCalled-1], m.errors[m.requestsCalled-1]
}

func newTestRequest(t *testing.T) *http.Request {
	t.Helper()
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://example.com", nil)
	require.NoError(t, err)

	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(strings.NewReader("")), nil
	}

	return req
}

func TestRetryRoundTripper_RoundTrip(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		setupMock      func() *MockTransport
		maxRetries     int
		expectedCalls  int
		expectedStatus int
		expectError    bool
	}{
		{
			name: "success_on_first_try",
			setupMock: func() *MockTransport {
				return &MockTransport{
					responses: []*http.Response{{StatusCode: http.StatusOK}},
					errors:    []error{nil},
				}
			},
			maxRetries:     3,
			expectedCalls:  1,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "success_after_one_retry",
			setupMock: func() *MockTransport {
				return &MockTransport{
					responses: []*http.Response{
						{StatusCode: http.StatusInternalServerError},
						{StatusCode: http.StatusOK},
					},
					errors: []error{nil, nil},
				}
			},
			maxRetries:     3,
			expectedCalls:  2,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "error_after_max_retries",
			setupMock: func() *MockTransport {
				return &MockTransport{
					responses: []*http.Response{
						{StatusCode: http.StatusInternalServerError},
						{StatusCode: http.StatusInternalServerError},
						{StatusCode: http.StatusInternalServerError},
						{StatusCode: http.StatusInternalServerError},
					},
					errors: []error{nil, nil, nil, nil},
				}
			},
			maxRetries:     3,
			expectedCalls:  4,
			expectedStatus: http.StatusInternalServerError,
			expectError:    false,
		},
		{
			name: "network_error_retry",
			setupMock: func() *MockTransport {
				return &MockTransport{
					responses: []*http.Response{nil, nil, {StatusCode: http.StatusOK}},
					errors: []error{
						errors.New("connection refused"),
						errors.New("connection refused"),
						nil,
					},
				}
			},
			maxRetries:     3,
			expectedCalls:  3,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "zero_max_retries",
			setupMock: func() *MockTransport {
				return &MockTransport{
					responses: []*http.Response{{StatusCode: http.StatusInternalServerError}},
					errors:    []error{nil},
				}
			},
			maxRetries:     0,
			expectedCalls:  1,
			expectedStatus: http.StatusInternalServerError,
			expectError:    false,
		},
		{
			name: "negative_max_retries",
			setupMock: func() *MockTransport {
				return &MockTransport{
					responses: []*http.Response{{StatusCode: http.StatusInternalServerError}},
					errors:    []error{nil},
				}
			},
			maxRetries:     -1,
			expectedCalls:  1,
			expectedStatus: http.StatusInternalServerError,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			mock := tt.setupMock()
			rt := retry.New(mock).WithMaxRetries(tt.maxRetries)
			req := newTestRequest(t)

			// Act
			resp, err := rt.RoundTrip(req)

			// Assert
			require.Equal(t, tt.expectedCalls, mock.requestsCount, "unexpected number of calls")

			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, tt.expectedStatus, resp.StatusCode)
			}

			if resp != nil && resp.Body != nil {
				resp.Body.Close()
			}
		})
	}
}

func TestRetryCallback(t *testing.T) {
	t.Parallel()

	// Arrange
	mock := &MockTransport{
		responses: []*http.Response{
			{StatusCode: http.StatusInternalServerError},
			{StatusCode: http.StatusInternalServerError},
			{StatusCode: http.StatusOK},
		},
		errors: []error{nil, nil, nil},
	}

	callbackCount := 0
	rt := retry.New(mock).
		WithMaxRetries(3).
		WithRetryCallback(func(attempt int, resp *http.Response, err error) {
			callbackCount++
			require.Equal(t, callbackCount, attempt, "unexpected callback attempt number")
		})

	// Act
	resp, err := rt.RoundTrip(newTestRequest(t))
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	// Assert
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, 2, callbackCount, "unexpected callback count")
}

func TestWithCustomRetryCondition(t *testing.T) {
	t.Parallel()

	// Arrange
	mock := &MockTransport{
		responses: []*http.Response{
			{StatusCode: http.StatusTooManyRequests},
			{StatusCode: http.StatusOK},
		},
		errors: []error{nil, nil},
	}

	rt := retry.New(mock).
		WithMaxRetries(3).
		WithRetryCondition(func(resp *http.Response, err error) bool {
			return err != nil || (resp != nil && (resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= http.StatusInternalServerError))
		})

	// Act
	resp, err := rt.RoundTrip(newTestRequest(t))
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	// Assert
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, 2, mock.requestsCount, "unexpected number of calls")
}

func TestBackoffDelay(t *testing.T) {
	t.Parallel()

	// Arrange
	mock := &MockTransport{
		responses: []*http.Response{
			{StatusCode: http.StatusInternalServerError},
			{StatusCode: http.StatusInternalServerError},
			{StatusCode: http.StatusOK},
		},
		errors: []error{nil, nil, nil},
	}

	rt := retry.New(mock).
		WithMaxRetries(2).
		WithBackoffConfig(5*time.Millisecond, 20*time.Millisecond, 2.0).
		WithJitter(false)

	start := time.Now()

	// Act
	resp, err := rt.RoundTrip(newTestRequest(t))
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	// Assert
	duration := time.Since(start)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.GreaterOrEqual(t, duration, 15*time.Millisecond, "backoff duration too short")
	require.LessOrEqual(t, duration, 100*time.Millisecond, "backoff duration too long")
}

func TestContextCancellation(t *testing.T) {
	t.Parallel()

	// Arrange
	mock := &MockTransport{
		responses: []*http.Response{
			{StatusCode: http.StatusInternalServerError},
			{StatusCode: http.StatusInternalServerError},
			{StatusCode: http.StatusOK},
		},
		errors: []error{nil, nil, nil},
	}

	rt := retry.New(mock).
		WithMaxRetries(2).
		WithBackoffConfig(100*time.Millisecond, 200*time.Millisecond, 2.0)

	ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.com", nil)
	require.NoError(t, err)

	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(strings.NewReader("")), nil
	}

	// Act
	resp, err := rt.RoundTrip(req)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Nil(t, resp)
}

func TestIntegrationWithServer(t *testing.T) {
	t.Parallel()

	// Arrange
	count := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
		if count < 3 {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	ctx := t.Context()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
	require.NoError(t, err)

	client := &http.Client{
		Transport: retry.New(http.DefaultTransport).WithMaxRetries(3),
	}

	// Act
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, 3, count, "unexpected number of server calls")
}
