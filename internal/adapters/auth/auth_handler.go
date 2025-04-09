package auth

import (
    "net/http"

    "github.com/AmiraliFarazmand/PTC_Task/internal/app"
    "github.com/gin-gonic/gin"
)

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

    c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user_id": user.ID.Hex()})
}
