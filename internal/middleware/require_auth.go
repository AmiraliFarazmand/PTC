package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AmiraliFarazmand/PTC_Task/internal/auth"
	"github.com/AmiraliFarazmand/PTC_Task/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func validateClaims(c *gin.Context, claims jwt.MapClaims) (auth.User, bool) {
	// check if the token is expired
	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		return auth.User{}, false
	}

	// Parse the user ID from the claims
	userID, err := bson.ObjectIDFromHex(claims["sub"].(string))
	if err != nil {
		return auth.User{}, false
	}

	// Find the user in the database
	var user auth.User
	err = db.UserCollection.FindOne(c.Request.Context(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return auth.User{}, false
	}
	return user, true
}

func RequireAuth(c *gin.Context) {
	// get the cookie off the request
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// TODO: move this to a env file
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
	user, ok := validateClaims(c, claims)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// set the user to the context
	c.Set("user", user)
	c.Next()
}
