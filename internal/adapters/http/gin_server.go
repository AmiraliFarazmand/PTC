package http

import (
	"github.com/AmiraliFarazmand/PTC_Task/internal/adapters/auth"
	"github.com/AmiraliFarazmand/PTC_Task/internal/core/app"
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

	// User routes
	authHandler := auth.NewAuthHandler(&s.UserService)
	//authentication routes
	r.POST("/signup", authHandler.Signup)
	r.POST("/login", authHandler.Login)
	r.GET("/validate", authHandler.RequireAuth, authHandler.ValidateHnadler)
	// Purchase routes
	r.POST("/purchase", authHandler.RequireAuth, s.processPurchase)
	r.PUT("/purchase/pay", authHandler.RequireAuth, s.confirmPayment)

	// Start the server
	r.Run(":8090")
}
