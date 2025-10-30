package handler

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterAPI(router *gin.Engine) {
	api := router.Group("/api")

	// Material routes
	materials := api.Group("/materials")
	{
		materials.GET("", h.GetMaterials)
		materials.GET("/:id", h.GetMaterial)
		materials.POST("", h.CreateMaterial)
		materials.PUT("/:id", h.UpdateMaterial)
		materials.DELETE("/:id", h.DeleteMaterial)
		materials.POST("/:id/image", h.UploadMaterialImage)
		materials.POST("/:id/add-to-draft", h.AddMaterialToDraft)
	}

	// Application routes
	applications := api.Group("/applications")
	{
		applications.GET("/cart", h.GetCartInfo)
		applications.GET("", h.GetApplications)
		applications.GET("/:id", h.GetApplication)
		applications.PUT("/:id", h.UpdateApplication)
		applications.PUT("/:id/submit", h.SubmitApplication)
		applications.PUT("/:id/complete", h.CompleteApplication)
		applications.PUT("/:id/reject", h.RejectApplication)
		applications.DELETE("/:id", h.DeleteApplication)

		// Application materials routes - ПЕРЕМЕЩАЕМ внутрь группы applications
		applicationMaterials := applications.Group("/:id/materials")
		{
			applicationMaterials.DELETE("/:materialId", h.RemoveMaterialFromApplication)
			applicationMaterials.PUT("/:materialId", h.UpdateApplicationMaterial)
		}
	}

	// User routes
	users := api.Group("/users")
	{
		users.POST("/register", h.Register)
		users.POST("/login", h.Login)
		users.POST("/logout", h.Logout)
		users.GET("/me", h.GetCurrentUser)
		users.PUT("/me", h.UpdateCurrentUser)
	}
}
