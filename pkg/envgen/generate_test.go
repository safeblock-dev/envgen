package envgen_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/safeblock-dev/envgen/pkg/envgen"
)

func TestEnvgen_Generate(t *testing.T) {
	t.Parallel()

	t.Run("simple", func(t *testing.T) {
		t.Parallel()
		// Create temporary directory for test files
		tmpDir := t.TempDir()

		// Create test files
		configPath := filepath.Join(tmpDir, "config.yaml")
		outputPath := filepath.Join(tmpDir, "output.go")
		templatePath := filepath.Join(tmpDir, "template.tmpl")

		// Create minimal config file with required group
		configContent := `groups:
  - name: App
    description: Application settings
    fields:
      - name: debug
        type: bool
        description: Enable debug mode`
		err := os.WriteFile(configPath, []byte(configContent), 0o600)
		require.NoError(t, err)

		// Create template file
		err = os.WriteFile(templatePath, []byte("package main"), 0o600)
		require.NoError(t, err)

		// Create a properly initialized Envgen instance
		envgenClient, err := envgen.New(t.Context(), envgen.Options{
			ConfigPath:   configPath,
			OutputPath:   outputPath,
			TemplatePath: templatePath,
		})
		require.NoError(t, err)

		err = envgenClient.Generate(t.Context())
		require.NoError(t, err)

		result, err := os.ReadFile(outputPath)
		require.NoError(t, err)
		require.Equal(t, "package main\n", string(result))
	})
}
