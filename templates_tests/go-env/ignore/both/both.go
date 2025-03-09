// Code generated by envgen. DO NOT EDIT.
// This file was automatically generated and should not be modified manually.

//go:generate envgen -c ../../ignore.yaml -o both.generated -t ../../../../templates/go-env

package config

// DatabaseConfig represents Database settings
type DatabaseConfig struct {
	Host string `env:"DB_HOST" envDefault:"localhost"` // Database host
	Port int `env:"DB_PORT" envDefault:"5432"` // Database port
}