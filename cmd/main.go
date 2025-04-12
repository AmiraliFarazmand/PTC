package main

import (
	"github.com/AmiraliFarazmand/PTC_Task/internal/adapters/auth"
	"github.com/AmiraliFarazmand/PTC_Task/internal/adapters/db"
	"github.com/AmiraliFarazmand/PTC_Task/internal/app"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	// Connect to MongoDB

	// client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {

		panic(err)
	}
	userCollection := client.Database("ParsTasmimDB").Collection("Users")

	// Initialize repositories
	userRepo := &db.MongoUserRepository{Collection: userCollection}

	// Initialize services
	userService := &app.UserService{UserRepo: userRepo}

	// Initialize handlers
	authHandler := &auth.AuthHandler{UserService: userService}

	// Setup routes
	r := gin.Default()
	r.POST("/signup", authHandler.Signup)
	r.POST("/login", authHandler.Login)
	r.GET("/validate", authHandler.RequireAuth, authHandler.ValidateHnadler)

	// Start the server
	r.Run(":8080")
}
