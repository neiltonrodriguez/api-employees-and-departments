package main

import (
	"fmt"
	"log"

	"api-employees-and-departments/config"
	_ "api-employees-and-departments/docs"
	"api-employees-and-departments/internal/db"
	"api-employees-and-departments/internal/domain/department"
	"api-employees-and-departments/internal/domain/employee"
	ginapi "api-employees-and-departments/internal/infrastructure/http/gin"
	"api-employees-and-departments/internal/infrastructure/migrations"
	"api-employees-and-departments/internal/infrastructure/persistence"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title API de Colaboradores e Departamentos
// @version 1.0
// @description API para gerenciamento de colaboradores e departamentos com hierarquia
// @termsOfService http://swagger.io/terms/

// @contact.name Suporte da API
// @contact.email suporte@api.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

func main() {
	// Load environment variables
	_ = godotenv.Load() // optional in production

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")

	// Run migrations
	if err := migrations.Run(database); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations executed successfully")

	// Initialize repositories
	employeeRepo := persistence.NewEmployeeRepository(database)
	departmentRepo := persistence.NewDepartmentRepository(database)

	// Create adapter for employee repository
	employeeAdapter := persistence.NewEmployeeAdapter(employeeRepo.(*persistence.EmployeeRepository))

	// Initialize services
	employeeService := employee.NewService(employeeRepo)
	departmentService := department.NewService(departmentRepo, employeeAdapter)

	// Initialize handlers
	employeeHandler := ginapi.NewEmployeeHandler(employeeService)
	departmentHandler := ginapi.NewDepartmentHandler(departmentService)
	managerHandler := ginapi.NewManagerHandler(departmentService, employeeService)

	// Setup Gin router (using New instead of Default to use custom middlewares)
	router := gin.New()

	// Setup routes (middlewares are configured inside SetupRoutes)
	ginapi.SetupRoutes(router, &ginapi.RouterConfig{
		EmployeeHandler:   employeeHandler,
		DepartmentHandler: departmentHandler,
		ManagerHandler:    managerHandler,
	})

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting server on %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
