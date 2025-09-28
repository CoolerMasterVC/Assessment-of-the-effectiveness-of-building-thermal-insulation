// cmd/migrate/main.go
package main

import (
	"Lab1/internal/app/ds"
	"Lab1/internal/app/dsn"
	"log"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	err = db.AutoMigrate(
		&ds.User{},
		&ds.Material{},
		&ds.Application{},
		&ds.ApplicationMaterial{},
	)
	if err != nil {
		log.Fatal("cant migrate db: ", err)
	}
	log.Println("Migration completed successfully")
}
