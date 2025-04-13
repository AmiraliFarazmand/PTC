package http

import (
	"github.com/AmiraliFarazmand/PTC_Task/internal/adapters/auth"
	"github.com/AmiraliFarazmand/PTC_Task/internal/app"
	"github.com/gin-gonic/gin"
)

type GinServer struct {
	PurchaseService app.PurchaseServiceImpl
	UserService     app.UserServiceImpl
}

func NewGinServer(purchaseService app.PurchaseServiceImpl, userService app.UserServiceImpl) *GinServer {
	return &GinServer{
		PurchaseService: purchaseService,
		UserService:     userService,
	}
}

func (s *GinServer) Start() {
	r := gin.Default()

	// Purchase routes
	r.POST("/purchase", func(c *gin.Context) {
		var body struct {
			UserID  string `json:"user_id" binding:"required"`
			Amount  int    `json:"amount" binding:"required"`
			Address string `json:"address" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}
		err := s.PurchaseService.CreatePurchase(body.UserID, body.Amount, body.Address)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, gin.H{"message": "Purchase created successfully"})
	})

	r.POST("/purchase/pay", func(c *gin.Context) {
		var body struct {
			PurchaseID string `json:"purchase_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}
		err := s.PurchaseService.ConfirmPayment(body.PurchaseID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Payment confirmed successfully"})
	})

	// User routes
	authHandler := auth.AuthHandler{UserService: &s.UserService}
	r.POST("/signup", authHandler.Signup)
	r.POST("/login", authHandler.Login)
	r.GET("/validate", authHandler.RequireAuth, authHandler.ValidateHnadler)

	// Start the server
	r.Run(":8080")
}
