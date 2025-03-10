// Code generated by envgen. DO NOT EDIT.
// This file was automatically generated and should not be modified manually.

//go:generate envgen -c ../types.yaml -o types.generated -t ../../../templates/go-env

package types
import (
	"net"
	"net/url"
	"time"
)

// AppConfig represents Application settings
type AppConfig struct {
	Env string `env:"ENV,required"` // Application environment (Possible values: development, staging, production)
	ApiUrl *url.URL `env:"API_URL,required"` // API endpoint
	RequestTimeout time.Duration `env:"REQUEST_TIMEOUT" envDefault:"30s"` // API request timeout
	ResponseTimeout time.Duration `env:"RESPONSE_TIMEOUT" envDefault:"30s"` // API response timeout
	AllowedIps []net.IP `env:"ALLOWED_IPS"` // List of allowed IP addresses
}