package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env:"ENV" env-default:"local" env-required:"true"`
	HTTPServer `yaml:"http_server"`
	Postgres   PostgresConfig `yaml:"postgres"`
}

type HTTPServer struct {
	Address           string        `yaml:"address" env-default:"localhost:8080"`
	ReadTimeout       time.Duration `yaml:"read_timeout" env-default:"15ms"`
	ProcessingTimeout time.Duration `yaml:"processing_timeout" env-default:"20ms"`
	WriteTimeout      time.Duration `yaml:"write_timeout" env-default:"15ms"`
	IdleTimeout       time.Duration `yaml:"idle_timeout" env-default:"30s"`
}

type PostgresConfig struct {
	DBName   string `yaml:"db_name" env:"PG_DATABASE_NAME" env-required:"true"`
	User     string `yaml:"user" env:"PG_USER" env-required:"true"`
	Password string `yaml:"password" env:"PG_PASSWORD" env-required:"true"`
	Port     string `yaml:"port" env:"PG_PORT" env-required:"true"`
	Host     string `yaml:"host" env:"PG_HOST" env-required:"true" env-default:"localhost"`
}

func MustLoad() *Config {
	env := flag.String("env", "local", "which config to use: local, prod, dev")
	flag.Parse()

	configFilePath := fmt.Sprintf("./config/%s.yaml", *env)

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configFilePath)
	}

	var cfg Config

	log.Printf("Reading config from: %s", configFilePath)

	if err := cleanenv.ReadConfig(configFilePath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	log.Printf("Config loaded: %+v\n", cfg)

	return &cfg
}
