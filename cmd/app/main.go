package main

import (
	"Project/internal/config"
	"Project/internal/repository/postgres"
	"Project/internal/server"
	"Project/pkg"
	"fmt"
	"log"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	gormDb := pkg.New(getConnString(conf))

	repo := postgres.New(gormDb)

	router := server.New(repo, conf)
	router.RunRouter(conf)
}

func getConnString(conf config.Config) string {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.DbHost,
		conf.DbPort,
		conf.DbUser,
		conf.DbPassword,
		conf.DbName,
	)

	return connectionString
}
