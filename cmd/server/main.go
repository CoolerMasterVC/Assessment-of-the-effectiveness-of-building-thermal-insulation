// cmd/app/main.go
package main

import (
	"Lab1/internal/app/config"
	"Lab1/internal/app/dsn"
	"Lab1/internal/app/handler"
	"Lab1/internal/app/repository"
	"Lab1/internal/pkg"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	router := gin.Default()
	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresString := dsn.FromEnv()
	fmt.Println("DSN:", postgresString)

	repo, err := repository.New(postgresString)
	if err != nil {
		logrus.Fatalf("error initializing repository: %v", err)
	}

	hand := handler.NewHandler(repo)

	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}
