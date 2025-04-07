package auth

import (
	"net/http"
	"strconv"
	"time"

	"github.com/AmiraliFarazmand/PTC_Task/internal/db"
	"github.com/AmiraliFarazmand/PTC_Task/internal/utils"
	"github.com/AmiraliFarazmand/PTC_Task/internal/validators"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive" //TODO: primitive for previous version
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

const tokenExpireTime int = 72 // expire time in hour

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"` // MongoDB ObjectID
	Username string             `bson:"username"`
	Password string             `bson:"password"`
}
type authRequest struct {
	Username string
	Password string
}

func createToken(userID uint) (string, error) {

	// Create Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * time.Duration(tokenExpireTime)).Unix(),
	})
	// TODO: should be stored in an environment variable or a config file
	secretKey := "SomeRandomSecretKey"
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

	// Insert into DB
	// user := bson.M{"username": body.Username, "password": string(hash)}
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

	var user User
	if validators.CheckUniquenessUsername(body.Username) != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid username or password")
		return
	}
	err := db.UserCollection.FindOne(c, bson.M{"username": body.Username}).Decode(&user)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid username or password")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid username or password")
		return
	}

	idUint,_ := strconv.ParseUint(user.ID.Hex(), 16, 0)
	tokenString, err := createToken(uint(idUint))
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
