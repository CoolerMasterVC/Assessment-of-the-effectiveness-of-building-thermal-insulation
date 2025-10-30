package handler

import (
	"Lab1/internal/app/auth"
	"Lab1/internal/app/ds"
	"net/http"

	"github.com/gin-gonic/gin"
)

// POST /api/users/register
func (h *Handler) Register(c *gin.Context) {
	var request struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// БЕЗ шифрования пароля
	user := ds.User{
		Login:    request.Login,
		Password: request.Password, // Сохраняем пароль как есть
	}

	if err := h.Repository.CreateUser(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    user.ID,
		"login": user.Login,
	})
}

// POST /api/users/login
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
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// БЕЗ проверки хеша - сравниваем пароли напрямую
	if user.Password != request.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	auth.SetCurrentUser(user)

	c.JSON(http.StatusOK, gin.H{
		"id":           user.ID,
		"login":        user.Login,
		"is_moderator": user.IsModerator,
	})
}

// POST /api/users/logout
func (h *Handler) Logout(c *gin.Context) {
	auth.SetCurrentUser(nil)
	c.JSON(http.StatusOK, gin.H{"status": "logged out"})
}

// GET /api/users/me
func (h *Handler) GetCurrentUser(c *gin.Context) {
	user := auth.GetCurrentUser()
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// PUT /api/users/me
func (h *Handler) UpdateCurrentUser(c *gin.Context) {
	user := auth.GetCurrentUser()
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
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
