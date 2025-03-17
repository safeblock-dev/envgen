// Package retry provides HTTP RoundTripper implementations with retry functionality.
package retry

import (
	"context"
	"crypto/rand"
	"math/big"
	"net/http"
	"time"
)

const (
	defaultInitialBackoff   = 100 * time.Millisecond
	defaultMaxBackoff       = 10 * time.Second
	defaultBackoffFactor    = 2.0
	defaultJitterMultiplier = 2
	defaultMaxRetries       = 3
	defaultJitterDivisor    = 2
)

// RoundTripper is an implementation of http.RoundTripper that adds retry
// functionality to an HTTP client.
type RoundTripper struct {
	// Next is the underlying RoundTripper to be wrapped with retry logic.
	// If nil, http.DefaultTransport will be used.
	Next http.RoundTripper

	// MaxRetries is the maximum number of retry attempts.
	// Default is 3 if not specified.
	MaxRetries int

	// InitialBackoff is the initial delay before retrying.
	// Default is 100ms if not specified.
	InitialBackoff time.Duration

	// MaxBackoff is the maximum delay between retries.
	// Default is 10s if not specified.
	MaxBackoff time.Duration

	// BackoffFactor determines how quickly the delay increases.
	// The next backoff is calculated as current_backoff * BackoffFactor.
	// Default is 2.0 if not specified.
	BackoffFactor float64

	// Jitter adds random variation to the backoff to prevent synchronized
	// retries from multiple clients. Set to true to enable.
	// Default is true.
	Jitter bool

	// ShouldRetry determines whether a request should be retried based
	// on the response and/or error. Default behavior retries on network
	// errors and 5xx status codes.
	ShouldRetry func(*http.Response, error) bool

	// OnRetry is called each time a retry is attempted, with the attempt number
	// (starting from 1), response (may be nil) and error (may be nil).
	OnRetry func(attempt int, resp *http.Response, err error)
}

// New creates a new RoundTripper with default configuration.
func New(next http.RoundTripper) *RoundTripper {
	defaultRetryCondition := func(resp *http.Response, err error) bool {
		return err != nil || (resp != nil && resp.StatusCode >= http.StatusInternalServerError)
	}

	return &RoundTripper{
		Next:           next,
		MaxRetries:     defaultMaxRetries,
		InitialBackoff: defaultInitialBackoff,
		MaxBackoff:     defaultMaxBackoff,
		BackoffFactor:  defaultBackoffFactor,
		Jitter:         true,
		OnRetry:        nil,
		ShouldRetry:    defaultRetryCondition,
	}
}

// prepareRequest creates a copy of the request and ensures the body can be reused.
func (rt *RoundTripper) prepareRequest(req *http.Request, attempt int) (*http.Request, bool, error) {
	reqCopy := req.Clone(req.Context())
	bodyReusable := req.Body == nil || req.GetBody != nil

	if attempt > 0 && req.Body != nil {
		if !bodyReusable {
			return nil, false, nil
		}

		body, err := req.GetBody()
		if err != nil {
			return nil, false, err
		}

		reqCopy.Body = body
	}

	return reqCopy, bodyReusable, nil
}

// getTransport returns the transport to use.
func (rt *RoundTripper) getTransport() http.RoundTripper {
	if rt.Next != nil {
		return rt.Next
	}

	return http.DefaultTransport
}

// getRetryConfig returns the retry configuration with defaults applied.
func (rt *RoundTripper) getRetryConfig() (int, time.Duration, time.Duration, float64) {
	maxRetries := rt.MaxRetries
	if maxRetries < 0 {
		maxRetries = 0
	}

	initialBackoff := rt.InitialBackoff
	if initialBackoff <= 0 {
		initialBackoff = defaultInitialBackoff
	}

	maxBackoff := rt.MaxBackoff
	if maxBackoff <= 0 {
		maxBackoff = defaultMaxBackoff
	}

	backoffFactor := rt.BackoffFactor
	if backoffFactor < 1.0 {
		backoffFactor = defaultBackoffFactor
	}

	return maxRetries, initialBackoff, maxBackoff, backoffFactor
}

// getRetryCondition returns the retry condition function.
func (rt *RoundTripper) getRetryCondition() func(*http.Response, error) bool {
	if rt.ShouldRetry != nil {
		return rt.ShouldRetry
	}

	return func(resp *http.Response, err error) bool {
		return err != nil || (resp != nil && resp.StatusCode >= http.StatusInternalServerError)
	}
}

// getBackoffDuration calculates the backoff duration with jitter.
func (rt *RoundTripper) getBackoffDuration(backoff time.Duration) time.Duration {
	if !rt.Jitter {
		return backoff
	}

	jitterRange := float64(backoff) / defaultJitterDivisor

	jitterDelta, err := rand.Int(rand.Reader, big.NewInt(int64(jitterRange*defaultJitterMultiplier)))
	if err != nil {
		return backoff
	}

	return backoff - time.Duration(jitterRange) + time.Duration(jitterDelta.Int64())
}

// handleResponse processes the response and determines if a retry is needed.
func (rt *RoundTripper) handleResponse(
	resp *http.Response,
	err error,
	attempt int,
	maxRetries int,
	shouldRetry func(*http.Response, error) bool,
) bool {
	if attempt >= maxRetries || !shouldRetry(resp, err) {
		return false
	}

	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}

	if rt.OnRetry != nil {
		rt.OnRetry(attempt+1, resp, err)
	}

	return true
}

// calculateNextBackoff calculates the next backoff duration.
func (rt *RoundTripper) calculateNextBackoff(
	currentBackoff time.Duration,
	maxBackoff time.Duration,
	backoffFactor float64,
) time.Duration {
	nextBackoff := time.Duration(float64(currentBackoff) * backoffFactor)
	if nextBackoff > maxBackoff {
		return maxBackoff
	}

	return nextBackoff
}

// waitForNextRetry waits for the next retry attempt or returns if context is done.
func (rt *RoundTripper) waitForNextRetry(
	ctx context.Context,
	backoff time.Duration,
) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(backoff):
		return nil
	}
}

// RoundTrip implements the http.RoundTripper interface.
func (rt *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	transport := rt.getTransport()
	maxRetries, initialBackoff, maxBackoff, backoffFactor := rt.getRetryConfig()
	shouldRetry := rt.getRetryCondition()

	var (
		resp *http.Response
		err  error
	)

	backoff := initialBackoff

	for attempt := 0; attempt <= maxRetries; attempt++ {
		reqCopy, bodyReusable, prepareErr := rt.prepareRequest(req, attempt)
		if prepareErr != nil {
			return nil, prepareErr
		}

		if !bodyReusable && attempt > 0 {
			return resp, err
		}

		resp, err = transport.RoundTrip(reqCopy)

		if !rt.handleResponse(resp, err, attempt, maxRetries, shouldRetry) {
			return resp, err
		}

		nextBackoff := rt.getBackoffDuration(backoff)
		if waitErr := rt.waitForNextRetry(reqCopy.Context(), nextBackoff); waitErr != nil {
			return nil, waitErr
		}

		backoff = rt.calculateNextBackoff(backoff, maxBackoff, backoffFactor)
	}

	return resp, err
}

// WithMaxRetries sets the maximum number of retry attempts.
func (rt *RoundTripper) WithMaxRetries(retries int) *RoundTripper {
	rt.MaxRetries = retries

	return rt
}

// WithBackoffConfig configures the backoff parameters.
func (rt *RoundTripper) WithBackoffConfig(initial, maximum time.Duration, factor float64) *RoundTripper {
	rt.InitialBackoff = initial
	rt.MaxBackoff = maximum
	rt.BackoffFactor = factor

	return rt
}

// WithJitter enables or disables jitter in the backoff calculations.
func (rt *RoundTripper) WithJitter(enable bool) *RoundTripper {
	rt.Jitter = enable

	return rt
}

// WithRetryCondition sets the function that determines whether to retry.
func (rt *RoundTripper) WithRetryCondition(condition func(*http.Response, error) bool) *RoundTripper {
	rt.ShouldRetry = condition

	return rt
}

// WithRetryCallback sets the function to be called on each retry attempt.
func (rt *RoundTripper) WithRetryCallback(callback func(int, *http.Response, error)) *RoundTripper {
	rt.OnRetry = callback

	return rt
}
