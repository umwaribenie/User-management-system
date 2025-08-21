package main

import (
	"context"
	"fmt"
	"log"
	"os"

	// Project imports (ensure these match your go.mod module name)
	"github.com/umwaribenie/final_user_management/controllers"
	"github.com/umwaribenie/final_user_management/docs" // Import generated docs for Swagger
	"github.com/umwaribenie/final_user_management/models"
	"github.com/umwaribenie/final_user_management/repositories"
	"github.com/umwaribenie/final_user_management/routes"
	"github.com/umwaribenie/final_user_management/services"
	"github.com/umwaribenie/final_user_management/utils"

	// External libraries
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var ctx = context.Background() // Global context for Redis

// @title User Management API
// @version 1.0
// @description This is the API for the user management system.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// 1. Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// 2. Read server configuration
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	// 3. Set JWT secret key
	utils.SetJWTSecret(os.Getenv("JWT_SECRET_KEY"))

	// 4. Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	// Ping Redis to check the connection
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Successfully connected to Redis!")

	// 5. Initialize database connection
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Successfully connected to the database!")

	// 6. Auto migrate the database models
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to auto migrate:", err)
	}
	log.Println("Database migration completed.")

	// 7. Initialize repositories
	userRepo := repositories.NewUserRepository(db)

	// 8. Initialize services
	// THIS IS THE FIX: Pass the redisClient to the auth service constructor
	authService := services.NewAuthService(userRepo, redisClient)
	userService := services.NewUserService(userRepo)

	// 9. Initialize controllers
	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(authService)

	// 10. Set up router and routes
	router := gin.Default()
	routes.SetupRouter(router, userController, authController)

	// 11. Setup Swagger
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 12. Start server
	log.Printf("Server running on port %s", serverPort)
	log.Fatal(router.Run(":" + serverPort))
}
