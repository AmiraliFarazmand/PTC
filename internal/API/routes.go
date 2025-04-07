package api

import (
	"github.com/AmiraliFarazmand/PTC_Task/internal/auth"
	"github.com/AmiraliFarazmand/PTC_Task/internal/middleware"
	"github.com/AmiraliFarazmand/PTC_Task/internal/purchase"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Auth Routes
	r.POST("/signup", auth.Signup)
	r.POST("/login", auth.Login)
	r.GET("/validate", middleware.RequireAuth, auth.ValidateIsAuthenticated)
	r.POST("/purchase", middleware.RequireAuth, purchase.CreatePurchase)

	return r
}
