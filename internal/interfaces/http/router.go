package http

import (
	"github.com/Mluciano1401/go-api/internal/interfaces/http/handlers"
	"github.com/Mluciano1401/go-api/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	jwtSecret string,
) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	users := api.Group("/users")
	users.Use(middleware.JWTAuth(jwtSecret))
	{
		users.GET("", userHandler.GetAll)
		users.GET("/:id", userHandler.GetByID)
		users.PUT("/:id", userHandler.Update)
		users.DELETE("/:id", userHandler.Delete)
	}

	return r
}
