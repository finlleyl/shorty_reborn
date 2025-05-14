package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string     `yaml:"env" env-default:"local"`
	HTTPServer HTTPServer `yaml:"http_server"`
	Database   Database   `yaml:"database"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Database struct {
	Driver   string        `yaml:"driver" env:"DB_DRIVER" env-default:"postgres"`
	Host     string        `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     int           `yaml:"port" env:"DB_PORT" env-default:"5432"`
	User     string        `yaml:"user" env:"DB_USER" env-required:"true"`
	Password string        `yaml:"password" env:"DB_PASSWORD" env-required:"true"`
	Name     string        `yaml:"name" env:"DB_NAME" env-required:"true"`
	SSLMode  string        `yaml:"ssl_mode" env:"DB_SSL_MODE" env-default:"disable"`
	Timeout  time.Duration `yaml:"timeout" env:"DB_TIMEOUT" env-default:"5s"`
}

func MustLoad() *Config {
	configPath, exists := os.LookupEnv("CONFIG_PATH")
	if !exists {
		if wd, err := os.Getwd(); err != nil {
			fmt.Printf("Current working directory: %s\n", wd)
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Failed to read config: %s", err)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("Failed to read env: %s", err)
	}

	return &cfg
}
