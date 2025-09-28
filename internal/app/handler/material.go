// internal/app/handler/material.go
package handler

import (
	"net/http"
	"strconv"

	"Lab1/internal/app/ds"

	"github.com/gin-gonic/gin"
)

func (h *Handler) IndexHandler(c *gin.Context) {
	searchQuery := c.Query("material_search")
	var materials []*ds.Material
	var err error

	if searchQuery != "" {
		materials, err = h.Repository.SearchMaterials(searchQuery)
	} else {
		materials, err = h.Repository.GetMaterials()
	}

	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	application, err := h.Repository.GetUserDraft(1)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	var cartCount int64
	if application != nil && application.Status == "черновик" {
		cartCount = h.Repository.GetApplicationMaterialsCount(application.ID)
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Materials":          materials,
		"SearchQuery":        searchQuery,
		"CartCount":          cartCount,
		"HasDraft":           application != nil && application.Status == "черновик",
		"CurrentApplication": application,
	})
}

func (h *Handler) MaterialHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	material, err := h.Repository.GetMaterialByID(uint(id))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	application, err := h.Repository.GetUserDraft(1)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	var cartCount int64
	if application != nil && application.Status == "черновик" {
		cartCount = h.Repository.GetApplicationMaterialsCount(application.ID)
	}

	c.HTML(http.StatusOK, "material.html", gin.H{
		"Material":  material,
		"CartCount": cartCount,
		"HasDraft":  application != nil && application.Status == "черновик",
	})
}

func (h *Handler) AddMaterialToApplicationHandler(c *gin.Context) {
	userID := uint(1)

	application, err := h.Repository.GetUserDraft(userID)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	if application == nil {
		application, err = h.Repository.CreateDraft(userID)
		if err != nil {
			h.errorHandler(c, http.StatusInternalServerError, err)
			return
		}
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	areaStr := c.PostForm("area")
	var area float64

	if areaStr == "" {
		area = 0
	} else {
		area, err = strconv.ParseFloat(areaStr, 64)
		if err != nil {
			h.errorHandler(c, http.StatusBadRequest, err)
			return
		}
	}

	err = h.Repository.AddMaterialToApplication(application.ID, uint(id), area)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusFound, "/")
}
