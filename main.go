package main

import (
	"devices-api-go/config"
	"devices-api-go/controller"
	"devices-api-go/middleware"
	"devices-api-go/repository"
	"devices-api-go/service"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// 1. Initialize the PostgreSQL connection pool
	config.ConnectDatabase()

	// 2. Instantiate and wire dependencies manually (No reflection magic)
	repo := repository.NewDeviceRepository()
	svc := service.NewDeviceService(repo)
	ctrl := controller.NewDeviceController(svc)

	// Initialize Middlewares
	validMiddleware := middleware.NewValidationMiddleware(repo)

	// 3. Create a default Gin router engine instance
	router := gin.Default()

	// 4. Attach the Global Error Handler Middleware to all routes
	router.Use(middleware.GlobalErrorHandler())

	// 5. Define HTTP Routes and bind them to Controller Methods
	api := router.Group("/api/v1/devices")
	{
		// 🎯 HERE IS WHERE WE DEFINE THE HTTP METHODS EXPLICITLY!

		// Public read endpoints (GET)
		api.GET("", ctrl.GetAll)      // Maps GET /api/v1/devices -> ctrl.GetAll
		api.GET("/:id", ctrl.GetByID) // Maps GET /api/v1/devices/:id -> ctrl.GetByID
		api.GET("/search/brand", ctrl.GetByBrand)
		api.GET("/search/state", ctrl.GetByState)

		// Write endpoint (POST)
		api.POST("", ctrl.Create) // Maps POST /api/v1/devices -> ctrl.Create

		// Protected modification endpoints (PUT, PATCH, DELETE)
		// We pass the GuardDeviceRules middleware BEFORE the controller method executes
		api.PUT("/:id", validMiddleware.GuardDeviceRules(), ctrl.FullUpdate)
		api.PATCH("/:id", validMiddleware.GuardDeviceRules(), ctrl.PartialUpdate)
		api.DELETE("/:id", validMiddleware.GuardDeviceRules(), ctrl.Delete)
	}

	// 6. Start the high-performance web server on port 8080
	log.Println("Go Devices API is running and listening on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start the web server: %v", err)
	}
}
