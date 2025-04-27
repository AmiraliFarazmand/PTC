package auth

import (
	"fmt"
	"net/http"

	"github.com/AmiraliFarazmand/PTC_Task/internal/ports"
	"github.com/AmiraliFarazmand/PTC_Task/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const tokenExpireTime int = 72

type AuthHandler struct {
	UserService    ports.UserService
	ProcessManager ports.ZeebeProcessManager
}

func NewAuthHandler(userService ports.UserService, processManager ports.ZeebeProcessManager) *AuthHandler {
	return &AuthHandler{
		UserService:    userService,
		ProcessManager: processManager,
	}
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.ProcessManager.StartSignupProcess(body.Username, body.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}
func (h *AuthHandler) Login(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Start login process and wait for result
	result, err := h.ProcessManager.StartLoginProcess(body.Username, body.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if !result.LoginValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": result.Error})
		return
	}

	// Set cookie with the token from Zeebe process
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", result.Token, 3600*tokenExpireTime, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"token":    result.Token,
		"username": result.Username,
	})
}


func (h *AuthHandler) ValidateHnadler(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	userDTO, ok := user.(ports.UserDTO)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cast user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "User is authenticated",
		"user_id":  userDTO.ID,
		"username": userDTO.Username,
	})
}

func (h *AuthHandler) RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No authorization token"})
		return
	}

	secretKey, _ := utils.ReadEnv("SECRET_KEY")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}

	user, err := h.UserService.FindUserByID(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	c.Set("user", user)
	c.Next()
}
