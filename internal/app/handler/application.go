// internal/app/handler/application.go
package handler

import (
	"Lab1/internal/app/ds"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ApplicationHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	application, err := h.Repository.GetApplicationByID(uint(id))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if application.Status == "удалён" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	appMaterials, err := h.Repository.GetApplicationMaterials(uint(id))
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	// Рассчитываем данные для отображения
	type CartItem struct {
		Material ds.Material
		Area     float64
		HeatLoss float64
		Savings  float64
	}

	var cartItems []CartItem
	totalSavings := 0.0

	for _, appMaterial := range appMaterials {
		// Расчет теплопотерь и экономии (упрощенный)
		heatLoss := appMaterial.Area * (application.IndoorTemp - application.OutdoorTemp) / (appMaterial.Material.Thickness / appMaterial.Material.Lambda)
		savings := heatLoss * 0.1 * 24 * 30 // упрощенный расчет экономии

		cartItems = append(cartItems, CartItem{
			Material: appMaterial.Material,
			Area:     appMaterial.Area,
			HeatLoss: heatLoss,
			Savings:  savings,
		})
		totalSavings += savings
	}

	userApplication, err := h.Repository.GetUserDraft(1)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	var cartCount int64
	var hasDraft bool
	if userApplication != nil && userApplication.Status == "черновик" {
		cartCount = h.Repository.GetApplicationMaterialsCount(userApplication.ID)
		hasDraft = true
	}

	c.HTML(http.StatusOK, "application.html", gin.H{
		"Application":  application,
		"CartItems":    cartItems,
		"TotalSavings": totalSavings,
		"CartCount":    cartCount,
		"HasDraft":     hasDraft,
	})
}

func (h *Handler) DeleteApplicationHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	err = h.Repository.DeleteApplication(uint(id))
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusFound, "/")
}
