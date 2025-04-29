package http

import (
	"github.com/AmiraliFarazmand/PTC_Task/internal/adapters/auth"
	"github.com/AmiraliFarazmand/PTC_Task/internal/ports"
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/zbc"
	"github.com/gin-gonic/gin"
)

type GinServer struct {
	PurchaseService ports.PurchaseService
	UserService     ports.UserService
	ProcessManager  ports.ZeebeProcessManager
	ZeebeClient     zbc.Client
	Router          *gin.Engine
}

func NewGinServer(purchaseService ports.PurchaseService, userService ports.UserService, processManager ports.ZeebeProcessManager, zeebeClient zbc.Client) *GinServer {
	server := &GinServer{
		PurchaseService: purchaseService,
		UserService:     userService,
		ProcessManager:  processManager,
		ZeebeClient:     zeebeClient,
		Router:          gin.Default(),
	}

	// User routes
	authHandler := auth.NewAuthHandler(server.UserService, server.ProcessManager)

	// Authentication routes
	server.Router.POST("/signup", authHandler.Signup)
	server.Router.POST("/login", authHandler.Login)
	server.Router.GET("/validate", authHandler.RequireAuth, authHandler.ValidateHnadler)

	// Purchase routes
	server.Router.POST("/purchase", authHandler.RequireAuth, server.processPurchase)
	server.Router.PUT("/purchase/pay", authHandler.RequireAuth, server.confirmPayment)

	return server
}

func (s *GinServer) Run() error {
	return s.Router.Run(":8313")
}
