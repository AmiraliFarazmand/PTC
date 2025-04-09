package auth

import (
	"net/http"
	"time"

	"github.com/AmiraliFarazmand/PTC_Task/internal/app"
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
