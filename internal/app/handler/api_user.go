package handler

import (
	"net/http"
	"time"

	"Lab1/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// @Summary Register new user
// @Description Register a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param input body object true "User registration data"
// @Success 201 {object} object "User created"
// @Failure 400 {object} object "Bad request"
// @Router /api/users/register [post]
func (h *Handler) Register(c *gin.Context) {
	var request struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем, не существует ли уже пользователь с таким логином
	existingUser, err := h.Repository.GetUserByLogin(request.Login)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}
	if existingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
		return
	}

	user := &ds.User{
		Login:       request.Login,
		Password:    request.Password,
		IsModerator: false,
	}

	if err := h.Repository.CreateUser(user); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    user.ID,
		"login": user.Login,
	})
}

// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags Users
// @Accept json
// @Produce json
// @Param input body object true "Login credentials"
// @Success 200 {object} object "Login successful"
// @Failure 401 {object} object "Unauthorized"
// @Router /api/users/login [post]
func (h *Handler) Login(c *gin.Context) {
	var request struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Repository.GetUserByLogin(request.Login)
	if err != nil || user == nil || user.Password != request.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Создаем JWT токен
	expirationTime := time.Now().Add(time.Duration(h.Config.JWTExpiresHours) * time.Hour)
	claims := &ds.JWTClaims{
		UserID:      user.ID,
		Login:       user.Login,
		IsModerator: user.IsModerator,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "materials-app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.Config.JWTSecret))
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      tokenString,
		"expires_at": expirationTime,
		"user": gin.H{
			"id":           user.ID,
			"login":        user.Login,
			"is_moderator": user.IsModerator,
		},
	})
}

// Logout выполняет выход пользователя
// @Summary User logout
// @Description Logout user and invalidate token
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object "Logout successful"
// @Router /api/users/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	const jwtPrefix = "Bearer "

	if len(tokenString) > len(jwtPrefix) {
		tokenString = tokenString[len(jwtPrefix):]

		// Добавляем токен в blacklist на оставшееся время
		claims := &ds.JWTClaims{}
		_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(h.Config.JWTSecret), nil
		})

		if err == nil {
			expiresIn := time.Until(time.Unix(claims.ExpiresAt, 0))
			if expiresIn > 0 {
				h.Redis.SetJWTBlacklist(c.Request.Context(), tokenString, expiresIn)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "logged out"})
}

// @Summary Get current user
// @Description Get current user profile information
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object "User profile"
// @Router /api/users/me [get]
func (h *Handler) GetCurrentUser(c *gin.Context) {
	user := h.GetCurrentUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateCurrentUser обновляет профиль пользователя
// @Summary Update current user
// @Description Update current user profile information
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body object true "Profile update data" example({"login": "string"})
// @Success 200 {object} object "Profile updated"
// @Router /api/users/me [put]
func (h *Handler) UpdateCurrentUser(c *gin.Context) {
	user := h.GetCurrentUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	var request struct {
		Login string `json:"login"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateUser := ds.User{Login: request.Login}
	if err := h.Repository.UpdateUser(user.ID, &updateUser); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}
