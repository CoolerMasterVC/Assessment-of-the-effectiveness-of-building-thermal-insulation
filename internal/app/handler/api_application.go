package handler

import (
	"Lab1/internal/app/ds"
	"Lab1/internal/calculations"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetApplications возвращает список заявок
// @Summary Get applications list
// @Description Get applications with filtering. For moderators - all applications, for users - only their applications
// @Tags Applications
// @Accept json
// @Produce json
// @Param status query string false "Status filter"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Security BearerAuth
// @Success 200 {array} ds.MaterialsApplication
// @Failure 401 {object} object "Unauthorized"
// @Failure 500 {object} object "Internal server error"
// @Router /api/applications [get]
func (h *Handler) GetApplications(c *gin.Context) {
	user := h.GetCurrentUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

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

	var applications []ds.MaterialsApplication
	var err error

	if user.IsModerator {
		// Модератор - все заявки кроме удаленных
		applications, err = h.Repository.GetApplications(status, startDate, endDate)
	} else {
		// Обычный пользователь - только свои заявки кроме удаленных
		applications, err = h.Repository.GetUserApplications(user.ID, status, startDate, endDate)
	}

	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, applications)
}

// CompleteApplication завершает заявку (только для модераторов)
// @Summary Complete application
// @Description Complete application (moderator only)
// @Tags Applications
// @Accept json
// @Produce json
// @Param id path int true "Application ID"
// @Security BearerAuth
// @Success 200 {object} object "Application completed"
// @Failure 400 {object} object "Bad request"
// @Failure 403 {object} object "Forbidden"
// @Failure 404 {object} object "Not found"
// @Router /api/applications/{id}/complete [put]
func (h *Handler) CompleteApplication(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	user := h.GetCurrentUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

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

	appMaterials, err := h.Repository.GetApplicationMaterials(uint(id))
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	totalSavings := calculations.CalculateTotalSavings(application, appMaterials)

	if err := h.Repository.CompleteApplication(uint(id), user.ID, totalSavings); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":        "completed",
		"total_savings": totalSavings,
	})
}

// RejectApplication отклоняет заявку (только для модераторов)
// @Summary Reject application
// @Description Reject application (moderator only)
// @Tags Applications
// @Accept json
// @Produce json
// @Param id path int true "Application ID"
// @Security BearerAuth
// @Success 200 {object} object "Application rejected"
// @Failure 400 {object} object "Bad request"
// @Failure 403 {object} object "Forbidden"
// @Failure 404 {object} object "Not found"
// @Router /api/applications/{id}/reject [put]
func (h *Handler) RejectApplication(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	user := h.GetCurrentUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

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

// GetCartInfo возвращает информацию о корзине
// @Summary Get cart info
// @Description Get current user's draft application info
// @Tags Applications
// @Produce json
// @Success 200 {object} object "Cart information"
// @Router /api/applications/cart [get]
func (h *Handler) GetCartInfo(c *gin.Context) {
	user := h.GetCurrentUserFromContext(c)

	var application *ds.MaterialsApplication
	var count int64 = 0

	if user != nil {
		application, _ = h.Repository.GetUserDraft(user.ID)
		if application != nil {
			count = h.Repository.GetApplicationMaterialsCount(application.ID)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"application_id": func() interface{} {
			if application != nil {
				return application.ID
			}
			return nil
		}(),
		"items_count": count,
	})
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

	var application ds.MaterialsApplication
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
	user := h.GetCurrentUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

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

// DELETE /api/applications/:id
func (h *Handler) DeleteApplication(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	user := h.GetCurrentUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

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
