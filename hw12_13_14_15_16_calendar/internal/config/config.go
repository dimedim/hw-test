package config

import (
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.

type Config struct {
	Logger LoggerConf `yaml:"logger"`
	App    App        `yaml:"app"`
	DB     PostgresDB `yaml:"database"`
	// TODO
}

type App struct {
	Port   string `yaml:"port" env:"APP_PORT" env-default:"8888"`
	Host   string `yaml:"host" env:"APP_HOST" env-default:"localhost"`
	DBType string `yaml:"db_type" env:"DB_TYPE"`
}
type LoggerConf struct {
	Level string `yaml:"level"`
	// TODO
}

type PostgresDB struct {
	Username string `yaml:"username" env:"DB_USER" envDefault:"postgres"`
	Password string `yaml:"password" env:"DB_PASSWORD" envDefault:"postgres"`
	Name     string `yaml:"db_name" env:"DB_NAME" envDefault:"postgres"`
	Host     string `yaml:"host" env:"DB_HOST" envDefault:"localhost"`
	Port     string `yaml:"port" env:"DB_PORT" envDefault:"5432"`
}

// TODO: функция то сделана, но используются аж 2 библы, лучше бы на змею перейти.
// Приоритеты:
// 1) env и .env файл соответственно;
// 2)Ямл файл
// 3)дефолт если не задано ничего. Если переменная пустая то останется пустой!
func LoadConfig(filepath string) *Config {

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

// return dsn like: <postgres://username:password@localhost:5432/database_name>
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
	// return "postgres://username:password@localhost:5432/database_name"
	return sb.String()
}

// TODO
