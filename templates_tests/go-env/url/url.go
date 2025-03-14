// Code generated by envgen. DO NOT EDIT.
// This file was automatically generated and should not be modified manually.

//go:generate envgen -c ../url.yaml -o url.generated -t https://raw.githubusercontent.com/safeblock-dev/envgen/main/templates/go-env

package urltest

// AppConfig represents Basic application settings
type AppConfig struct {
	Debug bool `env:"DEBUG" envDefault:"false"` // Enable debug mode
	Port int `env:"PORT,required"` // Server port
}