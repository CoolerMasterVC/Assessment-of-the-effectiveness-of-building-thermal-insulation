package handler

import (
	"Lab1/internal/app/auth"
	"Lab1/internal/app/ds"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
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

// POST /api/materials/:id/image
func (h *Handler) UploadMaterialImage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Получаем материал для проверки существования и получения старого изображения
	material, err := h.Repository.GetMaterialByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Material not found"})
		return
	}

	// ИСПРАВЛЕННЫЙ ПОДХОД - как в рабочем коде
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image file"})
		return
	}
	defer file.Close()

	// Проверяем что это изображение
	if !strings.Contains(header.Header.Get("Content-Type"), "image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is not an image"})
		return
	}

	// УДАЛЯЕМ старое изображение из Minio если оно есть
	if material.ImageURL != "" {
		h.deleteImageFromMinio(material.ImageURL)
	}

	// Генерируем название на латинице
	fileName := h.generateImageFileName(header.Filename, uint(id))

	// Загружаем в Minio - передаем file, fileName и fileSize
	imageURL, err := h.uploadImageToMinio(file, fileName, header.Size)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	// Обновляем URL в базе данных
	if err := h.Repository.UpdateMaterialImage(uint(id), imageURL); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"image_url": imageURL,
		"message":   "Image uploaded successfully",
	})
}

// uploadImageToMinio загружает изображение в Minio
func (h *Handler) uploadImageToMinio(file multipart.File, fileName string, fileSize int64) (string, error) {
	bucket := "images"
	ctx := context.Background()

	// Загружаем в Minio
	_, err := h.MinioClient.PutObject(ctx, bucket, fileName, file, fileSize, minio.PutObjectOptions{
		ContentType: "image/jpeg", // или определите тип из header
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to Minio: %v", err)
	}

	return fmt.Sprintf("http://localhost:9000/%s/%s", bucket, fileName), nil
}

// deleteImageFromMinio удаляет изображение из Minio
func (h *Handler) deleteImageFromMinio(imageURL string) {
	if h.MinioClient == nil || imageURL == "" {
		return
	}

	// Парсим URL чтобы получить bucket и object name
	parts := strings.Split(imageURL, "/")
	if len(parts) < 2 {
		return
	}

	bucket := parts[len(parts)-2]
	object := parts[len(parts)-1]

	ctx := context.Background()
	err := h.MinioClient.RemoveObject(ctx, bucket, object, minio.RemoveObjectOptions{})
	if err != nil {
		logrus.Warnf("Failed to delete image from Minio: %v", err)
	}
}

// generateImageFileName генерирует название файла на латинице
func (h *Handler) generateImageFileName(originalName string, materialID uint) string {
	ext := filepath.Ext(originalName)

	// Генерируем базовое имя на латинице
	baseName := fmt.Sprintf("material_%d_%d", materialID, time.Now().UnixNano())

	// Заменяем нелатинские символы если есть
	baseName = strings.ToLower(baseName)
	baseName = strings.ReplaceAll(baseName, " ", "_")

	return baseName + ext
}

// DELETE /api/materials/:id
func (h *Handler) DeleteMaterial(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Получаем материал чтобы получить URL изображения
	material, err := h.Repository.GetMaterialByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Material not found"})
		return
	}

	// УДАЛЯЕМ изображение из Minio если оно есть
	if material.ImageURL != "" {
		h.deleteImageFromMinio(material.ImageURL)
	}

	if err := h.Repository.DeleteMaterial(uint(id)); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
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
