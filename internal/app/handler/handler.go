// internal/app/handler/handler.go
package handler

import (
	"Lab1/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.GET("/test", func(c *gin.Context) {
		c.HTML(200, "test.html", gin.H{
			"Message": "Test page works!",
		})
	})

	router.GET("/", h.IndexHandler)
	router.GET("/material/:id", h.MaterialHandler)
	router.GET("/materials_aplication/:id", h.ApplicationHandler)
	router.POST("/materials_aplication/:id/delete", h.DeleteApplicationHandler)
	router.POST("/material/:id/add", h.AddMaterialToApplicationHandler)
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
