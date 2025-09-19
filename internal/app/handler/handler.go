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

	cart, _ := h.repo.GetCart()

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
	cart, err := h.repo.GetCart()
	if err != nil {
		logrus.Error(err)
	}

	var cartItems []struct {
		Material models.Material
		Area     float64
		HeatLoss float64
		Savings  float64
	}

	// Статичные данные для корзины
	staticItems := []struct {
		MaterialID int
		Area       float64
	}{
		{MaterialID: 1, Area: 15.5},
		{MaterialID: 3, Area: 22.0},
	}

	// Параметры расчета
	indoorTemp := 22.0   // °C
	outdoorTemp := -15.0 // °C
	tempDiff := indoorTemp - outdoorTemp

	for _, staticItem := range staticItems {
		material, _ := h.repo.GetMaterialByID(staticItem.MaterialID)

		// Расчет теплопотерь: Q = (ΔT * Area) / (thickness / lambda)
		// Предполагаем толщину утеплителя 0.1 м (10 см)
		thickness := 0.1
		heatLoss := (tempDiff * staticItem.Area) / (thickness / material.Lambda)

		// Расчет экономии (упрощенный)
		energyPrice := 5.0 // руб/кВт·ч
		monthlySavings := (heatLoss / 1000) * 24 * 30 * energyPrice * 0.3

		cartItems = append(cartItems, struct {
			Material models.Material
			Area     float64
			HeatLoss float64
			Savings  float64
		}{
			Material: material,
			Area:     staticItem.Area,
			HeatLoss: heatLoss,
			Savings:  monthlySavings,
		})
	}

	c.HTML(http.StatusOK, "cart.html", gin.H{
		"Cart":      cart,
		"CartItems": cartItems,
	})
}

func (h *Handler) AddToCartHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		c.Redirect(http.StatusFound, "/")
		return
	}

	areaStr := c.PostForm("area")
	area, err := strconv.ParseFloat(areaStr, 64)
	if err != nil || area <= 0 {
		area = 10.0
	}

	err = h.repo.AddToCart(id, area)
	if err != nil {
		logrus.Error(err)
	}

	c.Redirect(http.StatusFound, "/cart")
}
