package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AmiraliFarazmand/PTC_Task/internal/core/domain"
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

	// Start the Zeebe signup process
	if err := h.ProcessManager.StartSignupProcess(body.Username, body.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start signup process"})
		return
	}

	// The actual signup will be handled by the Zeebe worker
	c.JSON(http.StatusAccepted, gin.H{"message": "Signup process started"})
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

	// Start the Zeebe login process
	if err := h.ProcessManager.StartLoginProcess(body.Username, body.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start login process"})
		return
	}

	user, err := h.UserService.Login(body.Username, body.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userID := user.ID
	tokenString, err := CreateToken(userID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create token")
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*tokenExpireTime, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user_id": user.ID})
}

func CreateToken(userID string) (string, error) {

	// Create Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * time.Duration(tokenExpireTime)).Unix(),
	})
	secretKey, _ := utils.ReadEnv("SECRET_KEY")
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (h *AuthHandler) ValidateHnadler(c *gin.Context) {
	// Get the user from the context (set by RequireAuth middleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	// Cast the user to the correct type
	userObj, ok := user.(domain.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cast user"})
		return
	}

	// Respond with the user information
	c.JSON(http.StatusOK, gin.H{
		"message":  "User is authenticated",
		"user_id":  userObj.ID,
		"username": userObj.Username,
	})
}

func validateClaims(claims jwt.MapClaims, userService ports.UserService) (domain.User, bool) {
	// Check if the token is expired
	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		return domain.User{}, false
	}

	// Parse the user ID from the claims
	userID, ok := claims["sub"].(string)
	if !ok {
		return domain.User{}, false
	}
	// Use the UserService to find the user
	user, err := userService.FindUserByID(userID)
	if err != nil {
		return domain.User{}, false
	}

	return user, true
}

func (h *AuthHandler) RequireAuth(c *gin.Context) {
	// Get the cookie off the request
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		secretKey, _ := utils.ReadEnv("SECRET_KEY")
		return []byte(secretKey), nil
	})
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	claims, claimsOk := token.Claims.(jwt.MapClaims)
	if !claimsOk {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	// Validate claims and find the user
	user, ok := validateClaims(claims, h.UserService)
	if !ok {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	// Set the user to the context
	c.Set("user", user)
	c.Next()
}
