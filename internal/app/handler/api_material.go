package handler

import (
	"Lab1/internal/app/auth"
	"Lab1/internal/app/ds"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// GET /api/materials
func (h *Handler) GetMaterials(c *gin.Context) {
	filter := c.Query("filter")

	materials, err := h.Repository.GetMaterialsWithFilter(filter)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, materials)
}

// GET /api/materials/:id
func (h *Handler) GetMaterial(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	material, err := h.Repository.GetMaterialByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Material not found"})
		return
	}

	c.JSON(http.StatusOK, material)
}

// POST /api/materials
func (h *Handler) CreateMaterial(c *gin.Context) {
	var material ds.Material
	if err := c.ShouldBindJSON(&material); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repository.CreateMaterial(&material); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, material)
}

// PUT /api/materials/:id
func (h *Handler) UpdateMaterial(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var material ds.Material
	if err := c.ShouldBindJSON(&material); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repository.UpdateMaterial(uint(id), &material); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// DELETE /api/materials/:id
func (h *Handler) DeleteMaterial(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// УДАЛЯЕМ Minio вызовы - используем существующее хранилище
	if err := h.Repository.DeleteMaterial(uint(id)); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (h *Handler) UploadMaterialImage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Получаем URL из формы
	imageURL := c.PostForm("image_url")
	if imageURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image URL provided"})
		return
	}

	// Простая проверка что это похоже на URL
	if len(imageURL) < 10 || (!strings.HasPrefix(imageURL, "http://") && !strings.HasPrefix(imageURL, "https://")) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL format. Must start with http:// or https://"})
		return
	}

	// Сохраняем URL в базу данных
	if err := h.Repository.UpdateMaterialImage(uint(id), imageURL); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"image_url": imageURL})
}

// POST /api/materials/:id/add-to-draft
func (h *Handler) AddMaterialToDraft(c *gin.Context) {
	user := auth.GetCurrentUser()

	// Получаем или создаем черновик
	application, err := h.Repository.GetUserDraft(user.ID)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	if application == nil {
		application, err = h.Repository.CreateDraft(user.ID)
		if err != nil {
			h.errorHandler(c, http.StatusInternalServerError, err)
			return
		}
	}

	materialID, err := strconv.Atoi(c.Param("id"))
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

	if err := h.Repository.AddMaterialToApplication(application.ID, uint(materialID), request.Area); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"application_id": application.ID})
}
