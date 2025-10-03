package handler

import (
	"Lab1/internal/app/repository"
	"Lab1/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	repo *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		repo: r,
	}
}

func (h *Handler) IndexHandler(c *gin.Context) {
	searchQuery := c.Query("search")
	var materials []models.Material
	var err error

	if searchQuery != "" {
		materials, err = h.repo.GetMaterialsByName(searchQuery)
	} else {
		materials, err = h.repo.GetAllMaterials()
	}

	if err != nil {
		logrus.Error(err)
	}

	cart, _ := h.repo.GetCart(1)

	cartItemCount := len(cart.Items)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Materials":     materials,
		"Cart":          cart,
		"Search":        searchQuery,
		"CartItemCount": cartItemCount,
	})
}

func (h *Handler) MaterialHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		c.Redirect(http.StatusFound, "/")
		return
	}

	material, err := h.repo.GetMaterialByID(id)
	if err != nil {
		logrus.Error(err)
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.HTML(http.StatusOK, "material.html", gin.H{
		"Material": material,
	})
}

func (h *Handler) CartHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		c.Redirect(http.StatusFound, "/")
		return
	}

	cart, err := h.repo.GetCart(id)
	if err != nil {
		logrus.Error(err)
	}

	var cartItems []struct {
		Material models.Material
		Area     float64
		HeatLoss float64
		Savings  float64
	}

	for _, item := range cart.Items {
		material, _ := h.repo.GetMaterialByID(item.MaterialID)

		heatLoss := 1250.0
		monthlySavings := 2450.0

		cartItems = append(cartItems, struct {
			Material models.Material
			Area     float64
			HeatLoss float64
			Savings  float64
		}{
			Material: material,
			Area:     item.Area,
			HeatLoss: heatLoss,
			Savings:  monthlySavings,
		})
	}

	c.HTML(http.StatusOK, "materials_applications.html", gin.H{
		"Cart":      cart,
		"CartItems": cartItems,
	})
}
