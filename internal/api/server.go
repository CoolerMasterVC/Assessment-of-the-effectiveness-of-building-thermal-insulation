package api

import (
	"Lab1/internal/app/handler"
	"Lab1/internal/app/repository"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория")
	}

	handler := handler.NewHandler(repo)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	r.GET("/", handler.IndexHandler)
	r.GET("/material/:id", handler.MaterialHandler)
	r.GET("/materials_applications/:id", handler.CartMaterialHandler)
	//r.POST("/cart/add/:id", handler.AddToCartHandler)

	r.Run()
	log.Println("Server down")
}
