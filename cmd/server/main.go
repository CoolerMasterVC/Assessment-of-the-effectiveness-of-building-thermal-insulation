package main

import (
	"Lab1/internal/app/config"
	"Lab1/internal/app/dsn"
	"Lab1/internal/app/handler"
	"Lab1/internal/app/redis"
	"Lab1/internal/app/repository"
	"Lab1/internal/pkg"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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

	handler.SetupCORS(router)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresString := dsn.FromEnv()
	fmt.Println("DSN:", postgresString)

	// Инициализация Minio клиента
	minioClient, err := minio.New("minio:9000", &minio.Options{ // ← ДОЛЖНО БЫТЬ minio:9000
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("Failed to initialize Minio client: %v", err)
	}

	// Проверка подключения к Minio
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Проверяем существует ли bucket, если нет - создаем
	exists, err := minioClient.BucketExists(ctx, conf.MinioBucketName)
	if err != nil {
		log.Printf("Warning: Cannot check Minio bucket: %v", err)
	} else if !exists {
		log.Printf("Creating Minio bucket: %s", conf.MinioBucketName)
		err = minioClient.MakeBucket(ctx, conf.MinioBucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Printf("Warning: Cannot create bucket: %v", err)
		}
	}

	log.Printf("Minio client initialized successfully. Endpoint: %s, Bucket: %s",
		conf.MinioEndpoint, conf.MinioBucketName)

	repo, err := repository.New(postgresString)
	if err != nil {
		logrus.Fatalf("error initializing repository: %v", err)
	}

	// Инициализация Redis
	redisClient, err := redis.New(conf.RedisHost, conf.RedisPort)
	if err != nil {
		logrus.Fatalf("error initializing redis: %v", err)
	}

	hand := handler.NewHandler(repo, conf, redisClient, minioClient)

	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}
