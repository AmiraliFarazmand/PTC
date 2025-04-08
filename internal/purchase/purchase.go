package purchase

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AmiraliFarazmand/PTC_Task/internal/utils"
	"github.com/AmiraliFarazmand/PTC_Task/internal/auth"
	"github.com/AmiraliFarazmand/PTC_Task/internal/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Purchase struct {
	UserID    bson.ObjectID `bson:"user_id"`
	Amount    int           `bson:"amount"`
	CreatedAt time.Time     `bson:"created_at"`
	Status    string        `bson:"status"`
	PaymentID string        `bson:"payment_id"`
	Address   string        `bson:"address"`
}

type PurchaseRequest struct {
	Amount  int    `json:"amount" binding:"required"`
	Address string `json:"address" binding:"required"`
}

func CreatePurchase(c *gin.Context) {
	user, _ := c.Get("user")
	userObj, ok := user.(auth.User)
    if !ok {
        utils.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve user")
        return
    }

	var req PurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid JSON body")
	}
	// fmt.Printf("****%+v\n", user)

    purchase := Purchase{
        UserID:    userObj.ID, // Use the ID field from the user object
        Amount:    req.Amount,
        CreatedAt: time.Now(),
        Status:    "pending",
        PaymentID: "",
        Address:   req.Address,
	}

    _, err := db.PurchaseCollection.InsertOne(c.Request.Context(), purchase)
    if err != nil {
        utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create purchase")
        return
    }
	c.JSON(http.StatusCreated, gin.H{"message": "Purchase created successfully"})
    fmt.Printf("%+v\n", purchase)

	// call a function to process the purchase in a limited time otherwise drop it

}

