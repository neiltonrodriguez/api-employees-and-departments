package ginapi

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouterConfig struct {
	EmployeeHandler    *EmployeeHandler
	DepartmentHandler  *DepartmentHandler
	ManagerHandler     *ManagerHandler
}

// SetupRoutes configures all API routes
func SetupRoutes(router *gin.Engine, config *RouterConfig) {
	// Global middlewares
	router.Use(CORSMiddleware())
	router.Use(RequestID())
	router.Use(Logger())
	router.Use(Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"service": "api-employees-and-departments",
		})
	})

	// Swagger documentation
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Employee routes
		employees := v1.Group("/employees")
		{
			employees.GET("", config.EmployeeHandler.GetAll)
			employees.POST("/list", config.EmployeeHandler.List)
			employees.GET("/:id", config.EmployeeHandler.GetByID)
			employees.POST("", config.EmployeeHandler.Create)
			employees.PUT("/:id", config.EmployeeHandler.Update)
			employees.DELETE("/:id", config.EmployeeHandler.Delete)
		}

		// Department routes
		departments := v1.Group("/departments")
		{
			departments.GET("", config.DepartmentHandler.GetAll)
			departments.POST("/list", config.DepartmentHandler.List)
			departments.GET("/:id", config.DepartmentHandler.GetByID)
			departments.POST("", config.DepartmentHandler.Create)
			departments.PUT("/:id", config.DepartmentHandler.Update)
			departments.DELETE("/:id", config.DepartmentHandler.Delete)
		}

		// Manager routes
		managers := v1.Group("/managers")
		{
			managers.GET("/:id/employees", config.ManagerHandler.GetSubordinateEmployees)
		}
	}
}
