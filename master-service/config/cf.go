package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Service Service `yaml:"service"`
	Hosts   Hosts   `yaml:"hosts"`
}

type Service struct {
	Name       string `yaml:"name" envconfig:"SERVICE_NAME"`
	Version    string `yaml:"version" envconfig:"SERVICE_VERSION"`
	Address    string `yaml:"address" envconfig:"SERVICE_ADDRESS"`
	Port       int    `yaml:"port" envconfig:"SERVICE_PORT"`
	JwtSecret  string `yaml:"jwt_secret" envconfig:"SERVICE_JWT_SECRET"`
	StaticPath string `yaml:"static_path" envconfig:"STATIC_PATH"`
}

type Hosts struct {
	Database   Database   `yaml:"database"`
	Cache      Cache      `yaml:"cache"`
	Broker     Broker     `yaml:"broker"`
	Discovery  Discovery  `yaml:"discovery"`
	Sentry     Sentry     `yaml:"sentry"`
	Monitoring Monitoring `yaml:"monitoring"`
	Tracing    Tracing    `yaml:"tracing"`
	Services   Services   `yaml:"services"`
	Minio      Minio      `yaml:"minio"`
}

type Database struct {
	Name     string `yaml:"name" envconfig:"DATABASE_NAME"`
	Username string `yaml:"username" envconfig:"DATABASE_USERNAME"`
	Password string `yaml:"password" envconfig:"DATABASE_PASSWORD"`
	Address  string `yaml:"address" envconfig:"DATABASE_ADDRESS"`
	Port     int    `yaml:"port" envconfig:"DATABASE_PORT"`
	Driver   string `yaml:"driver" envconfig:"DATABASE_DRIVER"`
}
type Services struct {
	Web  string `yaml:"web" envconfig:"WEB_HOST"`
	Mail string `yaml:"mail" envconfig:"MAIL_SERVICE"`
}

type Cache struct {
	Address string `yaml:"address" envconfig:"CACHE_ADDRESS"`
	Port    int    `yaml:"port" envconfig:"CACHE_PORT"`
	Driver  string `yaml:"driver" envconfig:"CACHE_DRIVER"`
}

type Broker struct {
	Username string `yaml:"username" envconfig:"BROKER_USERNAME"`
	Password string `yaml:"password" envconfig:"BROKER_PASSWORD"`
	Address  string `yaml:"address" envconfig:"BROKER_ADDRESS"`
	Port     int    `yaml:"port" envconfig:"BROKER_PORT"`
	Driver   string `yaml:"driver" envconfig:"BROKER_DRIVER"`
}

func (b Broker) GetDriver() string {
	return b.Driver
}

func (b Broker) GetAddress() string {
	return b.Address
}

type Discovery struct {
	Username string `yaml:"username" envconfig:"DISCOVERY_USERNAME"`
	Password string `yaml:"password" envconfig:"DISCOVERY_PASSWORD"`
	Address  string `yaml:"address" envconfig:"DISCOVERY_ADDRESS"`
	Port     int    `yaml:"port" envconfig:"DISCOVERY_PORT"`
	Driver   string `yaml:"driver" envconfig:"DISCOVERY_DRIVER"`
}

type Monitoring struct {
	Username string `yaml:"username" envconfig:"MONITORING_USERNAME"`
	Password string `yaml:"password" envconfig:"MONITORING_PASSWORD"`
	Address  string `yaml:"address" envconfig:"MONITORING_ADDRESS"`
	Port     int    `yaml:"port" envconfig:"MONITORING_PORT"`
	Driver   string `yaml:"driver" envconfig:"MONITORING_DRIVER"`
}

type Tracing struct {
	Username string `yaml:"username" envconfig:"TRACING_USERNAME"`
	Password string `yaml:"password" envconfig:"TRACING_PASSWORD"`
	Address  string `yaml:"address" envconfig:"TRACING_ADDRESS"`
	Port     int    `yaml:"port" envconfig:"TRACING_PORT"`
	Driver   string `yaml:"driver" envconfig:"TRACING_DRIVER"`
}

type Sentry struct {
	Address string `yaml:"address" envconfig:"SENTRY_ADDRESS"`
	Env     string `yaml:"env" envconfig:"SENTRY_ENV"`
	Driver  string `yaml:"driver" envconfig:"SENTRY_DRIVER"`
}

type Minio struct {
	Address   string `yaml:"minio_address" envconfig:"MINIO_ADDRESS"`
	Console   string `yaml:"minio_console" envconfig:"MINIO_CONSOLE"`
	AccessKey string `yaml:"access_key" envconfig:"MINIO_ACCESS_KEY"`
	SecretKey string `yaml:"secret_key" envconfig:"MINIO_SECRET_KEY"`
	UseSSL    bool   `yaml:"use_ssl" envconfig:"MINIO_USE_SSL" default:"false"`
}

func (c *Config) Init() {
	var err error
	if err = godotenv.Load(".env"); err != nil {
		fmt.Println(err)
	}
	if err = envconfig.Process("", c); err != nil {
		fmt.Println(err)
	}
}
