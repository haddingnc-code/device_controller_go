package main

import (
	"devices-api-go/config"
	_ "devices-api-go/docs" // Blank import to initialize Swagger documentation specs
	handler "devices-api-go/internal/handler"
	"devices-api-go/internal/middleware"
	"devices-api-go/internal/repository"
	"devices-api-go/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // Explicit alias for static assets
	ginSwagger "github.com/swaggo/gin-swagger" // Explicit alias for Gin wrapper
	"log"
)

// @title           Devices API
// @version         1.0
// @description     A REST API for persisting and managing device resources protected by AOP validation rules.
// @host            localhost:8080
// @BasePath        /api/v1
func main() {
	// 1. Initialize the PostgreSQL connection pool
	config.ConnectDatabase()

	// 2. Instantiate the concrete data access layer
	repo := repository.NewDeviceRepository()

	// 3. Instantiate the core business service
	coreSvc := service.NewDeviceService(repo)

	// 4. Weaving: Wrap the core service inside our AOP validation aspect decorator
	advisedSvc := service.NewDeviceServiceAspect(coreSvc, repo)

	// 5. Inject the advised service proxy into the HTTP Controller handler
	ctrl := handler.NewDeviceController(advisedSvc)

	// 6. Create a default Gin router engine instance
	router := gin.Default()

	// 7. Attach the Global Error Handler Middleware to safely serialize business errors
	router.Use(middleware.GlobalErrorHandler())

	// 8. Expose the Swagger UI Interactive Documentation Route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 9. Define HTTP Routes directly bound to the validation-protected controller
	api := router.Group("/api/v1/devices")
	{
		// Public read endpoints (GET)
		api.GET("", ctrl.GetAll)
		api.GET("/:id", ctrl.GetByID)
		api.GET("/search/brand", ctrl.GetByBrand)
		api.GET("/search/state", ctrl.GetByState)

		// Write endpoint (POST)
		api.POST("", ctrl.Create)

		// Mutation endpoints fully protected by the service-layer AOP aspect architecture
		api.PUT("/:id", ctrl.FullUpdate)
		api.PATCH("/:id", ctrl.PartialUpdate)
		api.DELETE("/:id", ctrl.Delete)
	}

	// 10. Start the web server on port 8080
	log.Println("Go Devices API with Pure AOP Interface Architecture is running on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start the web server: %v", err)
	}
}
