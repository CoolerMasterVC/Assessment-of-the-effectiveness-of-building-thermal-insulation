package main

import (
	"Lab1/internal/app/config"
	"Lab1/internal/app/dsn"
	"Lab1/internal/app/handler"
	"Lab1/internal/app/redis"
	"Lab1/internal/app/repository"
	"Lab1/internal/pkg"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	_ "Lab1/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Materials App API
// @version 1.0
// @description API for heat insulation materials and applications

// @contact.name API Support
// @contact.url http://localhost:8080
// @contact.email support@materials-app.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT Authorization header using the Bearer scheme
func main() {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

	// Инициализация Redis
	redisClient, err := redis.New(conf.RedisHost, conf.RedisPort)
	if err != nil {
		logrus.Fatalf("error initializing redis: %v", err)
	}

	hand := handler.NewHandler(repo, conf, redisClient)

	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}
