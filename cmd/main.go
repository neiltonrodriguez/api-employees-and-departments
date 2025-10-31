package main

import (
	"fmt"
	"strconv"
	"time"

	"api-employees-and-departments/config"
	_ "api-employees-and-departments/docs"
	"api-employees-and-departments/internal/db"
	"api-employees-and-departments/internal/domain/department"
	"api-employees-and-departments/internal/domain/employee"
	infraCache "api-employees-and-departments/internal/infrastructure/cache"
	ginapi "api-employees-and-departments/internal/infrastructure/http/gin"
	"api-employees-and-departments/internal/infrastructure/logging"
	"api-employees-and-departments/internal/infrastructure/persistence"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
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
		logging.Fatal("Failed to load config", zap.Error(err))
	}

	if err := logging.InitLogger(cfg.AppEnv, cfg.LogLevel); err != nil {
		logging.Fatal("Failed to initialize logger", zap.Error(err))
	}
	defer logging.Sync()

	logging.Info("Application starting",
		zap.String("app_env", cfg.AppEnv),
		zap.String("log_level", cfg.LogLevel),
		zap.String("port", cfg.Port),
	)

	database, err := db.Connect(cfg)
	if err != nil {
		logging.Fatal("Failed to connect to database", zap.Error(err))
	}

	logging.Info("Database connected successfully",
		zap.String("host", cfg.DBHost),
		zap.String("database", cfg.DBName),
	)

	// Connect to Redis for caching
	redisClient, err := infraCache.NewRedisClient(cfg)
	if err != nil {
		logging.Fatal("Failed to connect to Redis", zap.Error(err))
	}

	logging.Info("Redis connected successfully",
		zap.String("host", cfg.RedisHost),
		zap.String("port", cfg.RedisPort),
	)

	// Create cache implementation
	cache := infraCache.NewRedisCache(redisClient)

	// Parse cache TTL
	cacheTTLSeconds, err := strconv.Atoi(cfg.CacheTTL)
	if err != nil {
		cacheTTLSeconds = 300 // Default 5 minutes
	}
	cacheTTL := time.Duration(cacheTTLSeconds) * time.Second

	logging.Info("Cache configured",
		zap.String("ttl", cacheTTL.String()),
	)

	// Note: Migrations are handled by Flyway before the application starts
	// See docker-compose.yml for Flyway configuration

	// Initialize repositories
	employeeRepo := persistence.NewEmployeeRepository(database)
	departmentRepo := persistence.NewDepartmentRepository(database)

	// Create adapter for employee repository
	employeeAdapter := persistence.NewEmployeeAdapter(employeeRepo.(*persistence.EmployeeRepository))

	// Create domain loggers with context for each service (DIP - Dependency Inversion Principle)
	employeeLogger := logging.NewZapLogger(logging.GetLogger().With(zap.String("service", "employee")))
	departmentLogger := logging.NewZapLogger(logging.GetLogger().With(zap.String("service", "department")))

	// Initialize services with logger and cache injection (DIP applied)
	employeeService := employee.NewService(employeeRepo, employeeLogger)
	departmentService := department.NewService(departmentRepo, employeeAdapter, departmentLogger, cache, cacheTTL)

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
	logging.Info("Starting HTTP server", zap.String("address", addr))

	if err := router.Run(addr); err != nil {
		logging.Fatal("Failed to start server", zap.Error(err))
	}
}
