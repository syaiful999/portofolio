package main

import (
	cf "moyo-master-service/config"
	"moyo-master-service/server"

	_ "github.com/asim/go-micro/plugins/registry/kubernetes/v4"
	_ "github.com/asim/go-micro/plugins/wrapper/monitoring/prometheus/v4"
	_ "github.com/asim/go-micro/plugins/wrapper/monitoring/victoriametrics/v4"
	"github.com/joho/godotenv"
	log "go-micro.dev/v4/logger"
)

var (
	conf = cf.Config{}
)

func main() {
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("./vault/secrets/.env"); err != nil {
			log.Fatal("error load .env")
		}
	}

	conf.Init()

	server.Init(conf)
}
