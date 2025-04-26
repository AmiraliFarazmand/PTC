package http

import (
	"github.com/AmiraliFarazmand/PTC_Task/internal/adapters/auth"
	"github.com/AmiraliFarazmand/PTC_Task/internal/ports"
	"github.com/AmiraliFarazmand/PTC_Task/internal/utils"
	"github.com/gin-gonic/gin"
)

type GinServer struct {
	PurchaseService ports.PurchaseService
	UserService     ports.UserService
	ProcessManager  ports.ZeebeProcessManager
}

func NewGinServer(purchaseService ports.PurchaseService, userService ports.UserService, processManager ports.ZeebeProcessManager) *GinServer {
	return &GinServer{
		PurchaseService: purchaseService,
		UserService:     userService,
		ProcessManager:  processManager,
	}
}

func (s *GinServer) Start() {
	r := gin.Default()

	// User routes
	authHandler := auth.NewAuthHandler(s.UserService, s.ProcessManager)
	//authentication routes
	r.POST("/signup", authHandler.Signup)
	r.POST("/login", authHandler.Login)
	r.GET("/validate", authHandler.RequireAuth, authHandler.ValidateHnadler)

	// Purchase routes
	r.POST("/purchase", authHandler.RequireAuth, s.processPurchase)
	r.PUT("/purchase/pay", authHandler.RequireAuth, s.confirmPayment)

	// Start the server
	ginPort, _ := utils.ReadEnv("GIN_PORT")
	r.Run(":" + ginPort)
}
