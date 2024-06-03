package main

import (
	"log"
	"main/internal/api"
	"main/internal/infrastructure"
	"main/internal/repository"
	"main/internal/service"
)

func main() {
	conf, err := infrastructure.GetConfigs()
	if err != nil {
		log.Fatalf("couldn't read configs: %v", err)
	}
	logger, err := infrastructure.InitLogger()
	if err != nil {
		log.Fatalf("couldn't init logger: %v", err)
	}
	db, err := infrastructure.GetDbConnection(conf)
	if err != nil {
		log.Fatalf("couldn't get Database connect: %v", err)
	}
	repo := repository.IntRepository(db)

	srv := service.GetService(repo, logger, conf)

	newApi := api.NewApi(srv, logger)
	newApi.InitRoutes(conf)
}
