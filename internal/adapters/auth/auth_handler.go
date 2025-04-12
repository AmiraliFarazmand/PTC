package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AmiraliFarazmand/PTC_Task/internal/app"
	"github.com/AmiraliFarazmand/PTC_Task/internal/domain"
	"github.com/AmiraliFarazmand/PTC_Task/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const tokenExpireTime int = 72
type AuthHandler struct {
    UserService *app.UserService
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

    err := h.UserService.Signup(body.Username, body.Password)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
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

    user, err := h.UserService.Login(body.Username, body.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    userID := user.ID.Hex()
	tokenString, err := createToken(userID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create token")
		return
	}
	// Set the token in a cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*tokenExpireTime, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
    c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user_id": user.ID.Hex()})
}

func createToken(userID string) (string, error) {

	// Create Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * time.Duration(tokenExpireTime)).Unix(),
	})
	// TODO: should be stored in an environment variable or a config file
	secretKey := "SomeRandomSecretKeyjsdfijsdfiojsjiofjsofhsidhfuiwhehwuifhwwiufhxciuv"
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
			"user_id":  userObj.ID.Hex(),
			"username": userObj.Username,
		})
	}


func validateClaims(c *gin.Context, claims jwt.MapClaims, userService *app.UserService) (domain.User, bool) {
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
        c.AbortWithStatus(http.StatusUnauthorized)
        return
    }

    // Parse the token
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }

        // TODO: move this to an env file
        secretKey := "SomeRandomSecretKeyjsdfijsdfiojsjiofjsofhsidhfuiwhehwuifhwwiufhxciuv"
        return []byte(secretKey), nil
    })
    if err != nil {
        c.AbortWithStatus(http.StatusUnauthorized)
        return
    }

    claims, claimsOk := token.Claims.(jwt.MapClaims)
    if !claimsOk {
        c.AbortWithStatus(http.StatusUnauthorized)
        return
    }

    // Validate claims and find the user
    user, ok := validateClaims(c, claims, h.UserService)
    if !ok {
        c.AbortWithStatus(http.StatusUnauthorized)
        return
    }

    // Set the user to the context
    c.Set("user", user)
    c.Next()
}