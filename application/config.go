package application

import (
	"strconv"
	"os"
)

// Config represents the configuration for the application.
type Config struct {
	RedisAddress string  // Address of the Redis server
	ServerPort   uint16  // Port on which the server should listen
}

// LoadConfig loads the configuration from environment variables.
func LoadConfig() Config {
	// Default configuration values
	cfg := Config{
		RedisAddress: "localhost:6379",
		ServerPort:   3000,
	}

	// Check if the REDIS_ADDR environment variable is set
	if redisAddr, exists := os.LookupEnv("REDIS_ADDR"); exists {
		// Override the default Redis address with the value from the environment variable
		cfg.RedisAddress = redisAddr
	}

	// Check if the SERVER_PORT environment variable is set
	if serverPort, exists := os.LookupEnv("SERVER_PORT"); exists {
		// Convert the environment variable value to a uint16
		if port, err := strconv.ParseUint(serverPort, 10, 16); err == nil {
			// Override the default server port with the parsed value
			cfg.ServerPort = uint16(port)
		}
	}

	// Return the loaded configuration
	return cfg
}