package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App         App     `yaml:"app"`
	Retry       Retry   `yaml:"retry"`
	Service     Service `yaml:"service"`
	Kafka       Kafka   `yaml:"kafka"`
	DatabaseURL string
}

type App struct {
	Port            string        `yaml:"port"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type Retry struct {
	Backoff     string  `yaml:"backoff"`
	MaxAttempts int     `yaml:"max_attempts"`
	Jitter      float64 `yaml:"jitter"`
}

type Service struct {
	CacheSize int `yaml:"cache_size"`
}

type Kafka struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
}

func Load(yamlConfigFilePath string) (*Config, error) {
	cfg := &Config{}

	data, err := os.ReadFile(yamlConfigFilePath)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	cfg.DatabaseURL = os.Getenv("DATABASE_URL")

	if len(cfg.Kafka.Brokers) == 0 && os.Getenv("KAFKA_BROKER") != "" {
		cfg.Kafka.Brokers = []string{os.Getenv("KAFKA_BROKER")}
	}
	if cfg.Kafka.Topic == "" {
		cfg.Kafka.Topic = os.Getenv("KAFKA_TOPIC")
	}

	return cfg, nil
}
