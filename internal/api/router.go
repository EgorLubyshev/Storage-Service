package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-username/storage-service/internal/api/handlers"
	"github.com/your-username/storage-service/internal/api/middleware"
	"github.com/your-username/storage-service/internal/storage/postgres"
)

func NewRouter(db *sql.DB) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), corsMiddleware())

	userRepo := postgres.NewUserRepo(db)
	userHandler := handlers.NewUserHandler(userRepo)

	v1 := router.Group("/api/v1")
	{
		v1.GET("/user/login", userHandler.LoginForm)
		v1.POST("/user/login", userHandler.Login)
		v1.GET("/user/register", userHandler.RegisterForm)
		v1.POST("/user/register", userHandler.Register)
		v1.GET("/lobby", middleware.RequireAuth(), handlers.Lobby)
		v1.POST("/lobby/data", middleware.RequireAuth(), handlers.LobbyData)
		v1.POST("/lobby/files", middleware.RequireAuth(), handlers.LobbyFiles)
	}

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return router
}

func corsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		ctx.Writer.Header().Set("Access-Control-Expose-Headers", "x-user-token")

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}
