package main

import (
	"log"
	"main/internal/api"
	"main/internal/repository"
	"main/internal/service"
	"main/pkg"
)

func main() {
	conf, err := pkg.GetConfigs()
	if err != nil {
		log.Fatalf("couldn't read configs: %v", err)
	}
	logger, err := pkg.GetLogger()
	if err != nil {
		log.Fatalf("couldn't init logger: %v", err)
	}
	db, err := pkg.GetDbConnection(conf.Db)
	if err != nil {
		log.Fatalf("couldn't get Database connect: %v", err)
	}
	repo := repository.GetRepository(db)

	srv := service.GetService(repo, logger, conf)

	newApi := api.NewApi(srv, logger)
	newApi.InitRoutes(conf)
}
