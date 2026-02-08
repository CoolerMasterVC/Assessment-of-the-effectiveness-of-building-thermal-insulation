package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DELETE /api/applications/:id/materials/:materialId
func (h *Handler) RemoveMaterialFromApplication(c *gin.Context) {
	appID, err := strconv.Atoi(c.Param("id")) // Меняем :appId на :id
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID"})
		return
	}

	materialID, err := strconv.Atoi(c.Param("materialId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material ID"})
		return
	}

	user := h.GetCurrentUserFromContext(c)
	application, err := h.Repository.GetApplicationByID(uint(appID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	if application.CreatorID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := h.Repository.RemoveMaterialFromApplication(uint(appID), uint(materialID)); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "removed"})
}

// PUT /api/applications/:id/materials/:materialId
func (h *Handler) UpdateApplicationMaterial(c *gin.Context) {
	appID, err := strconv.Atoi(c.Param("id")) // Меняем :appId на :id
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID"})
		return
	}

	materialID, err := strconv.Atoi(c.Param("materialId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material ID"})
		return
	}

	var request struct {
		Area float64 `json:"area"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := h.GetCurrentUserFromContext(c)
	application, err := h.Repository.GetApplicationByID(uint(appID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	if application.CreatorID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := h.Repository.UpdateApplicationMaterial(uint(appID), uint(materialID), request.Area); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}
