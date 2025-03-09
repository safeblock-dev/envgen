// Code generated by envgen. DO NOT EDIT.
// This file was automatically generated and should not be modified manually.

//go:generate go run github.com/safeblock-dev/envgen/cmd/envgen -c prefix.yaml -o prefix.generated -t ../../templates/go-env

package prefix

// AppConfig represents Application settings
type AppConfig struct {
	Debug bool `env:"APP_DEBUG" envDefault:"false"` // Enable debug mode
	Port int `env:"APP_PORT,required"` // Server port
}

// DatabaseConfig represents Database settings
type DatabaseConfig struct {
	Host string `env:"DB_HOST" envDefault:"localhost"` // Database host
	Port int `env:"DB_PORT" envDefault:"5432"` // Database port
}