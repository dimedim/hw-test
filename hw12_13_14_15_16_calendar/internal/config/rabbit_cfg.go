package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type RabbitConfig struct {
	Logger    LoggerConf `yaml:"logger"`
	DB        DB         `yaml:"db"`
	Rabbit    RabbitMQ   `yaml:"rabbitmq"`
	Scheduler Scheduler  `yaml:"scheduler"`
}

type DB struct {
	DSN  string `yaml:"dsn" env:"DB_DSN"`
	Type string `yaml:"type"`
}

type RabbitMQ struct {
	URL      string `yaml:"url"`
	ExchName string `yaml:"exch_name"`
	ExchType string `yaml:"exch_type"`
	Queue    string `yaml:"queue"`
	Key      string `yaml:"routing_key"`
}

type Scheduler struct {
	Interval   time.Duration `yaml:"interval"`
	DeleteDays int           `yaml:"delete_older_than_days"`
}

func LoadRabbitCfg(filepath string) *RabbitConfig {
	if err := godotenv.Load(); err != nil {
		panic("godotenv")
	}

	var cfg RabbitConfig
	err := cleanenv.ReadConfig(filepath, &cfg)
	if err != nil {
		panic("failed to read config from file: " + filepath + err.Error())
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("failed to load envs")
	}
	return &cfg
}
