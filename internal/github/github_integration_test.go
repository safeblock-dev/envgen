package github_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/safeblock-dev/envgen/internal/github"
)

const (
	testOwner      = "safeblock-dev"
	testRepo       = "envgen"
	testSHA1Length = 40
)

func skipShort(t *testing.T) {
	t.Helper()

	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
}

func setupIntegrationClient(t *testing.T, repo string) *github.Client {
	t.Helper()

	client, err := github.New(testOwner, repo)
	require.NoError(t, err)

	return client
}

func TestIntegration_GetLatestCommit(t *testing.T) {
	skipShort(t)
	t.Parallel()

	t.Run("success_valid_repo", func(t *testing.T) {
		t.Parallel()
		client := setupIntegrationClient(t, testRepo)

		sha, err := client.GetLatestCommit(t.Context())
		require.NoError(t, err)
		require.NotEmpty(t, sha)
		require.Len(t, sha, testSHA1Length)
	})

	t.Run("error_nonexistent_repo", func(t *testing.T) {
		t.Parallel()
		client := setupIntegrationClient(t, "nonexistent-repo-12345")

		sha, err := client.GetLatestCommit(t.Context())
		require.Error(t, err)
		require.Empty(t, sha)
	})
}

func TestIntegration_GetList(t *testing.T) {
	skipShort(t)
	t.Parallel()

	t.Run("success_templates_dir", func(t *testing.T) {
		t.Parallel()
		client := setupIntegrationClient(t, testRepo)

		files, err := client.GetList(t.Context(), "templates")
		require.NoError(t, err)
		require.NotEmpty(t, files)

		hasValidFile := false

		for _, file := range files {
			require.NotEmpty(t, file.Name)
			require.NotEmpty(t, file.Path)

			if file.Type == "file" && file.DownloadURL != "" {
				hasValidFile = true
			}
		}

		require.True(t, hasValidFile, "no valid file found in templates directory")
	})

	t.Run("error_nonexistent_dir", func(t *testing.T) {
		t.Parallel()
		client := setupIntegrationClient(t, testRepo)

		files, err := client.GetList(t.Context(), "nonexistent-dir-12345")
		require.Error(t, err)
		require.Nil(t, files)
	})
}

func TestIntegration_GetFile(t *testing.T) {
	skipShort(t)
	t.Parallel()

	t.Run("success_main_go", func(t *testing.T) {
		t.Parallel()
		client := setupIntegrationClient(t, testRepo)

		content, err := client.GetFile(t.Context(), "cmd/envgen/main.go")
		require.NoError(t, err)
		require.Contains(t, content, "package main")
		require.Contains(t, content, "func main()")
	})

	t.Run("error_nonexistent_file", func(t *testing.T) {
		t.Parallel()
		client := setupIntegrationClient(t, testRepo)

		content, err := client.GetFile(t.Context(), "nonexistent/file.go")
		require.Error(t, err)
		require.Empty(t, content)
	})
}

func TestIntegration_GetFileFromURL(t *testing.T) {
	skipShort(t)
	t.Parallel()

	t.Run("success_valid_url", func(t *testing.T) {
		t.Parallel()
		client := setupIntegrationClient(t, testRepo)

		files, err := client.GetList(t.Context(), "templates")
		require.NoError(t, err)
		require.NotEmpty(t, files)

		var fileURL string

		for _, file := range files {
			if file.Type == "file" && file.DownloadURL != "" {
				fileURL = file.DownloadURL

				break
			}
		}

		require.NotEmpty(t, fileURL, "no valid file URL found")

		content, err := client.GetFileFromURL(t.Context(), fileURL)
		require.NoError(t, err)
		require.NotEmpty(t, content)
	})

	t.Run("error_invalid_url", func(t *testing.T) {
		t.Parallel()
		client := setupIntegrationClient(t, testRepo)

		content, err := client.GetFileFromURL(t.Context(), "https://raw.githubusercontent.com/safeblock-dev/envgen/main/nonexistent.txt")
		require.Error(t, err)
		require.Empty(t, content)
	})
}

func TestIntegration_GetStandardTemplateURL(t *testing.T) {
	skipShort(t)
	t.Parallel()

	t.Run("success_template", func(t *testing.T) {
		t.Parallel()
		client := setupIntegrationClient(t, testRepo)

		url, err := client.GetStandardTemplateURL(t.Context(), "go-env")
		require.NoError(t, err)
		require.Contains(t, url, "templates/go-env")
		require.Contains(t, url, "raw.githubusercontent.com")
	})
}
