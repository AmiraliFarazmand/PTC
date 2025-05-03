package http

import (
	"net/http"

	"github.com/AmiraliFarazmand/PTC_Task/internal/ports"
	"github.com/gin-gonic/gin"
)

func (s *GinServer) processPurchase(c *gin.Context) {
	var body struct {
		Amount  int    `json:"amount" binding:"required"`
		Address string `json:"address" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

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

	// Start purchase process through Zeebe
	result, err := s.ProcessManager.StartPurchaseProcess(userDTO.ID, body.Amount, body.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !result.IsValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(201, gin.H{
		"message":     "Purchase created successfully",
		"purchase_id": result.PurchaseID,
	})
}

func (s *GinServer) confirmPayment(c *gin.Context) {
	purchaseID := c.Param("id")
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

	// Verify purchase belongs to user and update status
	err := s.PurchaseService.ConfirmPayment(purchaseID, userDTO.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Payment confirmed successfully"})
}
