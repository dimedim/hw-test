package config

import (
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Logger LoggerConf `yaml:"logger"`
	HTTP   HTTPServer `yaml:"http_server"`
	DB     PostgresDB `yaml:"database"`
	GRPC   GRPCServer `yaml:"grpc_server"`
}

type HTTPServer struct {
	Port        string        `yaml:"port" env:"APP_PORT" env-default:"8080"`
	Host        string        `yaml:"host" env:"APP_HOST" env-default:"localhost"`
	DBType      string        `yaml:"db_type" env:"DB_TYPE"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}
type LoggerConf struct {
	Level     string `yaml:"level"`
	LogFolder string `yaml:"log_folder"`
}

type PostgresDB struct {
	Username          string `yaml:"username" env:"DB_USER" env-default:"postgres"`
	Password          string `yaml:"password" env:"DB_PASSWORD" env-default:"postgres"`
	Name              string `yaml:"db_name" env:"DB_NAME" env-default:"postgres"`
	Host              string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port              string `yaml:"port" env:"DB_PORT" env-default:"5432"`
	MigrationFilepath string `yaml:"migrations_folder"`
}

type GRPCServer struct {
}

// Приоритеты:
// 1) env и .env файл соответственно;
// 2)Ямл файл
// 3)дефолт если не задано ничего. Если переменная пустая то останется пустой!
func MustLoad(filepath string) *Config {
	if err := godotenv.Load(); err != nil {
		panic("godotenv")
	}

	var cfg Config
	err := cleanenv.ReadConfig(filepath, &cfg)
	if err != nil {
		panic("failed to read config from file: " + filepath)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("failed to load envs")
	}

	return &cfg
}

// GetPostgresDSN return dsn like: <postgres://username:password@localhost:5432/database_name>?sslmode=disable
func (c *Config) GetPostgresDSN() string {
	sb := strings.Builder{}

	sb.WriteString(`postgres://`)
	sb.WriteString(c.DB.Username)
	sb.WriteString(":")
	sb.WriteString(c.DB.Password)
	sb.WriteString("@")
	sb.WriteString(c.DB.Host)
	sb.WriteString(":")
	sb.WriteString(c.DB.Port)
	sb.WriteString("/")
	sb.WriteString(c.DB.Name)
	sb.WriteString("?sslmode=disable")
	return sb.String()
}
