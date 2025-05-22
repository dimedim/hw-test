package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.

type Config struct {
	Logger LoggerConf `yaml:"logger"`
	App    App        `yaml:"app"`
	// TODO
}

type App struct {
	Port string `yaml:"port" env:"APP_PORT" env-default:"8888"`
	Host string `yaml:"host" env:"APP_HOST" env-default:"localhost"`
}
type LoggerConf struct {
	Level string `yaml:"level"`
	// TODO
}

// TODO: функция то сделана, но используются аж 2 библы, лучше бы на змею перейти.
// Приоритеты:
// 1) env и .env файл соответственно;
// 2)Ямл файл
// 3)дефолт если не задано ничего. Если переменная пустая то останется пустой!
func LoadConfig(filepath string) Config {

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

	return cfg
}

// TODO
