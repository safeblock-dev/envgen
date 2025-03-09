// Code generated by envgen. DO NOT EDIT.
// This file was automatically generated and should not be modified manually.

//go:generate go run github.com/safeblock-dev/envgen/cmd/envgen -c ../../ignore.yaml -o groups.generated -t ../../../../templates/go-env

package config
import (
	"time"
)

// AppConfig represents Application settings
type AppConfig struct {
	Env string `env:"ENV,required"` // Application environment (Possible values: development, staging, production)
	Timeout time.Duration `env:"TIMEOUT" envDefault:"30s"` // Operation timeout
}