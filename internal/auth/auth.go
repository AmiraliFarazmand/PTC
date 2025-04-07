package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AmiraliFarazmand/PTC_Task/internal/db"
	"github.com/AmiraliFarazmand/PTC_Task/internal/utils"
	"github.com/AmiraliFarazmand/PTC_Task/internal/validators"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

const tokenExpireTime int = 72 // expire time in hour

type User struct {
	ID       bson.ObjectID `bson:"_id,omitempty"` // MongoDB ObjectID
	Username string        `bson:"username"`
	Password string        `bson:"password"`
}
type authRequest struct {
	Username string
	Password string
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

func Signup(c *gin.Context) {
	var body authRequest
	// Validate format of request
	if c.ShouldBindJSON(&body) != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	// Validate Username and Password of request body
	if err := validators.ValidateUsernamePassword(body.Username, body.Password); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	user := User{
		Username: body.Username,
		Password: string(hash),
	}
	// db.InsertIntoCollection(db.UserCollection, user)
	_, err = db.UserCollection.InsertOne(c, user)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created"})
}

func Login(c *gin.Context) {
	var body authRequest
	// Validate format of request
	if c.ShouldBindJSON(&body) != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	// var user User

	// err := db.UserCollection.FindOne(c.Request.Context(), bson.M{"username": body.Username}).Decode(&user)
	// if err != nil {
	// 	fmt.Println("case2", err)
	// 	fmt.Println(user.Username, user.Password)
	// 	utils.RespondWithError(c, http.StatusBadRequest, "Invalid username or password")
	// 	return
	// }
	var raw bson.M
	err := db.UserCollection.FindOne(c.Request.Context(), bson.M{"username": body.Username}).Decode(&raw)
	if err != nil {
		fmt.Println("case2", err)
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid username or password")
		return
	}
	claimedPassword := raw["password"].(string)
	err = bcrypt.CompareHashAndPassword([]byte(claimedPassword), []byte(body.Password))
	if err != nil {
		fmt.Println("case3", err.Error(), "####\n", claimedPassword, "****\n", body.Password)
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid username or password")
		return
	}
	userID := raw["_id"].(bson.ObjectID).Hex()

	fmt.Println("userID", userID)
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
}

func ValidateIsAuthenticated(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "UnAuthorized user")
		return
	}

	username := user.(User).Username
	c.JSON(http.StatusOK, gin.H{
		"message":  "I am Authenticated",
		"username": username,
	})
}
