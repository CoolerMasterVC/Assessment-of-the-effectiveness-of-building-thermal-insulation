package handler

import (
	"Lab1/internal/app/config"
	"Lab1/internal/app/ds"
	"Lab1/internal/app/redis"
	"Lab1/internal/app/repository"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository  *repository.Repository
	Config      *config.Config
	MinioClient *minio.Client
	Redis       *redis.Client
}

func NewHandler(r *repository.Repository, cfg *config.Config, redisClient *redis.Client, minioClient *minio.Client) *Handler {
	return &Handler{
		Repository:  r,
		Config:      cfg,
		Redis:       redisClient,
		MinioClient: minioClient, // Добавьте эту строку
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	// Глобальный middleware для аутентификации
	router.Use(h.AuthMiddleware())

	h.RegisterTemplates(router)
	h.RegisterStatic(router)

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.GET("/test", func(c *gin.Context) {
		c.HTML(200, "test.html", gin.H{
			"Message": "Test page works!",
		})
	})

	router.Use(func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next()
	})

	api := router.Group("/api")

	// Публичные маршруты (доступны всем)
	public := api.Group("")
	{
		// Аутентификация (доступна только гостям)
		public.POST("/users/register", h.Register)
		public.POST("/users/login", h.Login)
		// Просмотр материалов (доступен всем)
		public.GET("/materials", h.GetMaterials)
		public.GET("/materials/:id", h.GetMaterial)
		// Корзина (доступна всем, но для гостя возвращает 0,0)
		public.GET("/mat_applics/cart", h.GetCartInfo)
	}

	// Защищенные маршруты (требуют аутентификации)
	protected := api.Group("")
	protected.Use(h.RequireAuth())
	{
		// Пользовательские маршруты
		users := protected.Group("/users")
		{
			users.GET("/me", h.GetCurrentUser)
			users.PUT("/me", h.UpdateCurrentUser)
			users.POST("/logout", h.Logout)
		}

		// Маршруты заявок (для всех аутентифицированных)
		applications := protected.Group("/mat_applics")
		{
			applications.GET("", h.GetApplications)
			applications.GET("/:id", h.GetApplication)
			applications.PUT("/:id", h.UpdateApplication)
			applications.PUT("/:id/submit", h.SubmitApplication)
			applications.DELETE("/:id", h.DeleteApplication)

			applicationMaterials := applications.Group("/:id/materials")
			{
				applicationMaterials.DELETE("/:materialId", h.RemoveMaterialFromApplication)
				applicationMaterials.PUT("/:materialId", h.UpdateApplicationMaterial)
			}
		}

		// Маршруты материалов (чтение для всех, изменение для модераторов)
		materials := protected.Group("/materials")
		{
			materials.POST("", h.RequireModerator(), h.CreateMaterial)
			materials.PUT("/:id", h.RequireModerator(), h.UpdateMaterial)
			materials.DELETE("/:id", h.RequireModerator(), h.DeleteMaterial)
			materials.POST("/:id/add-to-draft", h.AddMaterialToDraft)
			materials.POST("/:id/image", h.RequireModerator(), h.UploadMaterialImage)
		}

		// Маршруты только для модераторов
		moderator := protected.Group("")
		moderator.Use(h.RequireModerator())
		{
			moderator.PUT("/mat_applics/:id/complete", h.CompleteApplication)
			moderator.PUT("/mat_applics/:id/reject", h.RejectApplication)
		}
	}

	// HTML маршруты
	router.GET("/", h.IndexHandler)
	router.GET("/material/:id", h.MaterialHandler)
	router.GET("/materials_aplication/:id", h.ApplicationHandler)
	router.POST("/materials_aplication/:id/delete", h.DeleteApplicationHandler)
	router.POST("/material/:id/add", h.AddMaterialToApplicationHandler)

	router.NoRoute(h.NotFoundHandler)
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.Static("/static", "./static")
}

func (h *Handler) RegisterTemplates(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
}

// AuthMiddleware проверяет JWT токен
func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.Next() // Пропускаем, если нет токена (гость)
			return
		}

		const jwtPrefix = "Bearer "
		if !strings.HasPrefix(tokenString, jwtPrefix) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString = tokenString[len(jwtPrefix):]

		// Проверяем blacklist
		isBlacklisted, err := h.Redis.IsJWTBlacklisted(c.Request.Context(), tokenString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			c.Abort()
			return
		}
		if isBlacklisted {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token revoked"})
			c.Abort()
			return
		}

		// Парсим токен
		claims := &ds.JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(h.Config.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Сохраняем пользователя в контекст
		c.Set("user", &ds.User{
			ID:          claims.UserID,
			Login:       claims.Login,
			IsModerator: claims.IsModerator,
		})

		c.Next()
	}
}

// RequireAuth middleware требует аутентификации
func (h *Handler) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists || user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireModerator middleware требует прав модератора
func (h *Handler) RequireModerator() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists || user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}

		currentUser := user.(*ds.User)
		if !currentUser.IsModerator {
			c.JSON(http.StatusForbidden, gin.H{"error": "moderator access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// GetCurrentUserFromContext получает пользователя из контекста
func (h *Handler) GetCurrentUserFromContext(c *gin.Context) *ds.User {
	user, exists := c.Get("user")
	if !exists || user == nil {
		return nil
	}
	return user.(*ds.User)
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
