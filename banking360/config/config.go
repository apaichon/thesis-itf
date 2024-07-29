package config

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config represents application configuration
type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	MaxConnection int
	MaxIdleConnection int
	ConnectionLifeTime int
	RestApiPort int
	GraphQLPort int
}

const (
	DB_HOST     = "DB_HOST"
	DB_PORT     = "DB_PORT"
	DB_USER     = "DB_USER"
	DB_PASSWORD = "DB_PASSWORD"
	MAX_CONNECTION= "MAX_CONNECTION"
	MAX_IDLE_CONNECTION= "MAX_IDLE_CONNECTION"
	CONNECTION_LIFE_TIME= "CONNECTION_LIFE_TIME"
	REST_API_PORT= "REST_API_PORT"
	GRAPHQL_PORT= "GRAPHQL_PORT"
)

var instance *Config
var once sync.Once

// LoadConfig loads the configuration from environment variables
func NewConfig() *Config {
	once.Do(func() {
		
		relativePath := "../../config/.env"
		// Get the absolute path
		absolutePath, err := filepath.Abs(relativePath)
		if err != nil {
			fmt.Println(err)
			return
		}
		// Load environment variables from .env file
		if err := godotenv.Load(absolutePath); err != nil {
			fmt.Println("Failed to load env variables:", err)
			return
		}

		viper.AutomaticEnv()

		// Create a Config instance and set values from Viper
		instance = &Config{
			DBHost:     viper.GetString(DB_HOST),
			DBPort:     viper.GetInt(DB_PORT),
			DBUser:     viper.GetString(DB_USER),
			DBPassword: viper.GetString(DB_PASSWORD),
			MaxConnection: viper.GetInt (MAX_CONNECTION),
			MaxIdleConnection: viper.GetInt (MAX_IDLE_CONNECTION),
			ConnectionLifeTime: viper.GetInt (CONNECTION_LIFE_TIME),
			RestApiPort: viper.GetInt (REST_API_PORT),
			GraphQLPort: viper.GetInt (GRAPHQL_PORT),
		}
	})
	return instance
}

func (cfg *Config) GetDns() string {
	dns := fmt.Sprintf("tcp://%s:%v?username=%s&password=%s", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword)
	return dns
}

// GetConfig returns the singleton configuration instance
func GetConfig() *Config {
	return instance
}
