package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// TODO: разобраться с конфигом Postgres (и убрать storagePath)

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-default:"local" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
}

type HTTPServer struct {
	Address           string        `yaml:"address" env-default:"localhost:8080"`
	ReadTimeout       time.Duration `yaml:"read_timeout" env-default:"15ms"`
	ProcessingTimeout time.Duration `yaml:"processing_timeout" env-default:"20ms"`
	WriteTimeout      time.Duration `yaml:"write_timeout" env-default:"15ms"`
	IdleTimeout       time.Duration `yaml:"idle_timeout" env-default:"30s"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
