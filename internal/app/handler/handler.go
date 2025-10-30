package handler

import (
	"Lab1/internal/app/repository"
	"net/http"
	"strings"

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
	h.RegisterAPI(router)

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

	h.RegisterStatic(router)

	router.NoRoute(h.NotFoundHandler)
}

// NotFoundHandler обрабатывает все несуществующие маршруты
func (h *Handler) NotFoundHandler(c *gin.Context) {
	// Логируем попытку доступа к несуществующему маршруту
	logrus.Warnf("404 Not Found: %s %s", c.Request.Method, c.Request.URL.Path)

	// Проверяем, это API запрос или HTML запрос
	if strings.HasPrefix(c.Request.URL.Path, "/api/") {
		// Для API возвращаем JSON ошибку
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Endpoint not found",
			"path":  c.Request.URL.Path,
		})
	} else {
		// Для HTML запросов редиректим на главную
		c.Redirect(http.StatusFound, "/")
	}
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	// УБИРАЕМ LoadHTMLGlob из статики - он может вызывать конфликты
	// router.LoadHTMLGlob("templates/*") // УДАЛИТЬ ЭТУ СТРОКУ

	// Регистрируем статику только один раз
	router.Static("/static", "./static")
}

// Отдельный метод для загрузки шаблонов
func (h *Handler) RegisterTemplates(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
