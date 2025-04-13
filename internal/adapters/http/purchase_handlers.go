package http

import "github.com/gin-gonic/gin"

func (s *GinServer)confirmPayment(c *gin.Context) {
	var body struct {
		PurchaseID string `json:"purchase_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	err := s.PurchaseService.ConfirmPayment(body.PurchaseID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Payment confirmed successfully"})
}


func (s *GinServer)processPurchase(c *gin.Context) {
	var body struct {
		UserID  string `json:"user_id" binding:"required"`
		Amount  int    `json:"amount" binding:"required"`
		Address string `json:"address" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	err := s.PurchaseService.CreatePurchase(body.UserID, body.Amount, body.Address)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "Purchase created successfully"})
}