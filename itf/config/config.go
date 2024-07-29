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
	QueueDBHost          string
	QueueDBPort          int
	QueueDBUser          string
	QueueDBPassword      string
	LoadBalancePort int
	TempDBPath string
	TotalMemQueues int
	PersistencyBatchSize int
	QueueManagerInterval int
	OrchestrationManagerInterval int
	MaxConnection int
	MaxIdleConnection int
	ConnectionLifeTime int
}

const (
	QUEUE_DB_HOST          = "QUEUE_DB_HOST"
	QUEUE_DB_PORT          = "QUEUE_DB_PORT"
	QUEUE_DB_USER          = "QUEUE_DB_USER"
	QUEUE_DB_PASSWORD      = "QUEUE_DB_PASSWORD"
	LOAD_BALANCE_PORT = "LOAD_BALANCE_PORT"
	TEMP_DB_PATH = "TEMP_DB_PATH"
	TOTAL_MEM_QUEUES = "TOTAL_MEM_QUEUES"
	PERSISTENTCY_BATCH_SIZE = "PERSISTENTCY_BATCH_SIZE"
	QUEUE_MANAGER_INTERVAL = "QUEUE_MANAGER_INTERVAL"
	ORCHESTRATOR_MANAGER_INTERVAL ="ORCHESTRATOR_MANAGER_INTERVAL"
	MAX_CONNECTION= "MAX_CONNECTION"
	MAX_IDLE_CONNECTION= "MAX_IDLE_CONNECTION"
	CONNECTION_LIFE_TIME= "CONNECTION_LIFE_TIME"
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

		tempDBPath, err := filepath.Abs(viper.GetString(TEMP_DB_PATH))
		if err != nil {
			fmt.Println(err)
			tempDBPath =""
		}

		// Create a Config instance and set values from Viper
		instance = &Config{
			QueueDBHost:          viper.GetString(QUEUE_DB_HOST),
			QueueDBPort:          viper.GetInt(QUEUE_DB_PORT),
			QueueDBUser:          viper.GetString(QUEUE_DB_USER),
			QueueDBPassword:      viper.GetString(QUEUE_DB_PASSWORD),
			TempDBPath: tempDBPath,
			TotalMemQueues: viper.GetInt(TOTAL_MEM_QUEUES),
			PersistencyBatchSize: viper.GetInt(PERSISTENTCY_BATCH_SIZE),
			QueueManagerInterval: viper.GetInt(QUEUE_MANAGER_INTERVAL),
			OrchestrationManagerInterval: viper.GetInt(ORCHESTRATOR_MANAGER_INTERVAL),
			MaxConnection: viper.GetInt (MAX_CONNECTION),
			MaxIdleConnection: viper.GetInt (MAX_IDLE_CONNECTION),
			ConnectionLifeTime: viper.GetInt (CONNECTION_LIFE_TIME),
		}
	})
	return instance
}

func (cfg *Config) GetDns() string {
	dns:=fmt.Sprintf("tcp://%s:%v?username=%s&password=%s", cfg.QueueDBHost, cfg.QueueDBPort, cfg.QueueDBUser, cfg.QueueDBPassword)
	return dns
}

// GetConfig returns the singleton configuration instance
func GetConfig() *Config {
	return instance
}
