// Code generated by envgen. DO NOT EDIT.
// This file was automatically generated and should not be modified manually.

// is URL? true

package url

// App Basic application settings
type App struct {
	Debug bool `env:"DEBUG" envDefault:"false"` // Enable debug mode
	Port int `env:"PORT,required"` // Server port
}
