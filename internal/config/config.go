package config

import (
	"goth-todo/internal/core/repository"
	"goth-todo/internal/core/services"
	"goth-todo/internal/handlers"
	"goth-todo/internal/middleware"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type App struct {
	Router *gin.Engine
	DB     *pgxpool.Pool
}

// for later
type Repositories struct {
	UserRepo repository.UserRepository
	TaskRepo repository.TaskRepository
}

type Services struct {
	UserService services.UserService
}

type Handlers struct {
	UserHandler handlers.UserHandler
}

func loadEnv() {
	// Optional: only load .env in local environments
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, assuming CI/production")
		}
	}
}

func RunApp(db *pgxpool.Pool) *App {

	loadEnv()

	r := gin.Default()
	r.Static("/public", "./public")
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(middleware.CSPMiddleware())

	// Initialize repositories
	taskRepo := repository.NewTaskRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	taskService := services.NewTaskService(*taskRepo)
	userService := services.NewUserService(*userRepo)

	// Initialize handlers
	taskHandlers := handlers.NewTaskHandler(taskService)
	contentHandlers := handlers.NewContentHandlers()
	userHandlers := handlers.NewUserHandler(userService)

	// might remove this section, just make this whole process one big function. Don't love the idea but we'll see.
	// Setup routes
	// router.SetupRoutes(r, taskHandlers, contentHandlers)

	r.GET("/", contentHandlers.GetDashboardPage)

	// User routes
	userRoutes := r.Group("/user")
	{
		userRoutes.POST("/login", userHandlers.Login)
	}

	// Task routes
	taskRoutes := r.Group("/task")
	{
		taskRoutes.GET("/get", taskHandlers.GetTasks)
		taskRoutes.POST("/add", taskHandlers.AddTask)
		// taskRoutes.POST("/toggle/:id", taskHandlers.ToggleTask)
		// taskRoutes.POST("/delete/:id", taskHandlers.DeleteTask)
	}
	taskRoutes.Use(middleware.JWTMiddleware())

	return &App{
		Router: r,
		DB:     db,
	}
}
