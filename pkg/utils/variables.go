package utils

// Configuration holds all the configuration parameters for the Litmus client
type Configuration struct {
	Endpoint string `envconfig:"LITMUS_ENDPOINT" default:"http://localhost:8080"`
	Username string `envconfig:"LITMUS_USERNAME" default:"admin"`
	Password string `envconfig:"LITMUS_PASSWORD" default:"litmus"`
}

var Config Configuration
