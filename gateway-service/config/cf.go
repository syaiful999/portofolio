package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Service  Service  `yaml:"service"`
	Hosts    Hosts    `yaml:"hosts"`
	Services Services `yaml:"services"`
}

type Service struct {
	Name      string `yaml:"name" envconfig:"SERVICE_NAME"`
	Version   string `yaml:"version" envconfig:"SERVICE_VERSION"`
	Address   string `yaml:"address" envconfig:"SERVICE_ADDRESS"`
	Port      int    `yaml:"port" envconfig:"SERVICE_PORT"`
	PortFiber int    `yaml:"port_fiber" envconfig:"SERVICE_PORT_FIBER"`
	JwtSecret string `yaml:"jwt_secret" envconfig:"SERVICE_JWT_SECRET"`
}

type Services struct {
	MasterDataURL string `yaml:"master_data_url" envconfig:"MASTER_DATA_SERVICE_HOST"`
}

type Hosts struct {
	Cache Cache `yaml:"cache"`
}

type Cache struct {
	Address string `yaml:"address" envconfig:"CACHE_ADDRESS"`
	Port    int    `yaml:"port" envconfig:"CACHE_PORT"`
	Driver  string `yaml:"driver" envconfig:"CACHE_DRIVER"`
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
