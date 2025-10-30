package handler

import (
	"Lab1/internal/app/auth"
	"Lab1/internal/app/ds"
	"Lab1/internal/calculations"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GET /api/applications/cart
func (h *Handler) GetCartInfo(c *gin.Context) {
	user := auth.GetCurrentUser()

	application, err := h.Repository.GetUserDraft(user.ID)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	var count int64
	if application != nil {
		count = h.Repository.GetApplicationMaterialsCount(application.ID)
	}

	c.JSON(http.StatusOK, gin.H{
		"application_id": application.ID,
		"items_count":    count,
	})
}

// GET /api/applications
func (h *Handler) GetApplications(c *gin.Context) {
	status := c.Query("status")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate *time.Time

	if startDateStr != "" {
		if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &t
		}
	}

	if endDateStr != "" {
		if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = &t
		}
	}

	applications, err := h.Repository.GetApplications(status, startDate, endDate)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, applications)
}

// GET /api/applications/:id
func (h *Handler) GetApplication(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	application, err := h.Repository.GetApplicationByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	c.JSON(http.StatusOK, application)
}

// PUT /api/applications/:id
func (h *Handler) UpdateApplication(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var application ds.MaterialsApplication // ИСПРАВЛЕНО
	if err := c.ShouldBindJSON(&application); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repository.UpdateApplication(uint(id), &application); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// PUT /api/applications/:id/submit
func (h *Handler) SubmitApplication(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Проверяем что заявка принадлежит пользователю
	user := auth.GetCurrentUser()
	application, err := h.Repository.GetApplicationByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	if application.CreatorID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Проверяем обязательные поля
	if application.TotalArea == 0 || len(application.Materials) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Application must have materials and total area"})
		return
	}

	if err := h.Repository.SubmitApplication(uint(id)); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "submitted"})
}

// PUT /api/applications/:id/complete
func (h *Handler) CompleteApplication(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	user := auth.GetCurrentUser()
	if !user.IsModerator {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only moderators can complete applications"})
		return
	}

	// Рассчитываем экономию
	application, err := h.Repository.GetApplicationByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	totalSavings := calculations.CalculateTotalSavings(application, application.Materials)

	if err := h.Repository.CompleteApplication(uint(id), user.ID, totalSavings); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":        "completed",
		"total_savings": totalSavings,
	})
}

// PUT /api/applications/:id/reject
func (h *Handler) RejectApplication(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	user := auth.GetCurrentUser()
	if !user.IsModerator {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only moderators can reject applications"})
		return
	}

	if err := h.Repository.RejectApplication(uint(id), user.ID); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "rejected"})
}

// DELETE /api/applications/:id
func (h *Handler) DeleteApplication(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	user := auth.GetCurrentUser()
	application, err := h.Repository.GetApplicationByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	if application.CreatorID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := h.Repository.DeleteApplication(uint(id)); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
