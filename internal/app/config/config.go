package config

import (
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceHost     string
	ServicePort     int
	JWTSecret       string
	JWTExpiresHours int
	RedisHost       string
	RedisPort       int

	MinioExternalEndpoint string

	// Minio настройки
	MinioEndpoint   string `yaml:"minio_endpoint" env:"MINIO_ENDPOINT" env-default:"localhost:9000"`
	MinioAccessKey  string `yaml:"minio_access_key" env:"MINIO_ACCESS_KEY" env-default:"minioadmin"`
	MinioSecretKey  string `yaml:"minio_secret_key" env:"MINIO_SECRET_KEY" env-default:"minioadmin"`
	MinioBucketName string `yaml:"minio_bucket_name" env:"MINIO_BUCKET_NAME" env-default:"images"`
	MinioUseSSL     bool   `yaml:"minio_use_ssl" env:"MINIO_USE_SSL" env-default:"false"`
	MinioRegion     string `yaml:"minio_region" env:"MINIO_REGION" env-default:"us-east-1"`
}

func NewConfig() (*Config, error) {
	var err error

	configName := "config"
	_ = godotenv.Load()
	if os.Getenv("CONFIG_NAME") != "" {
		configName = os.Getenv("CONFIG_NAME")
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")
	viper.WatchConfig()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	log.Info("config parsed")

	return cfg, nil
}
