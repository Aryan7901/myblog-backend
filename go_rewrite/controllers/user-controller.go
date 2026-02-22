package controllers

import (
	"context"
	"os"
	"time"

	"backend/config"
	"backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(422, gin.H{"message": "Invalid inputs passed, please check your data."})
		return
	}

	var existingUser models.User
	err := config.UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(422, gin.H{"message": "User exists already, please login instead."})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		c.JSON(500, gin.H{"message": "Could not create user, please try again."})
		return
	}
	user.Password = string(hashedPassword)
	user.Blogs = []primitive.ObjectID{}

	result, err := config.UserCollection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(500, gin.H{"message": "Signing up failed, please try again later."})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    result.InsertedID.(primitive.ObjectID).Hex(),
		"email":     user.Email,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"exp":       time.Now().Add(time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		c.JSON(500, gin.H{"message": "Signing up failed, please try again later."})
		return
	}

	c.JSON(201, models.UserResponse{
		ID:        result.InsertedID.(primitive.ObjectID).Hex(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Token:     tokenString,
	})
}

func Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var loginReq models.LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(422, gin.H{"message": "Invalid inputs passed, please check your data."})
		return
	}

	var user models.User
	err := config.UserCollection.FindOne(ctx, bson.M{"email": loginReq.Email}).Decode(&user)
	if err != nil {
		c.JSON(403, gin.H{"message": "Invalid credentials, could not log you in."})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		c.JSON(403, gin.H{"message": "Invalid credentials, could not log you in."})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    user.ID.Hex(),
		"email":     user.Email,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"exp":       time.Now().Add(time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		c.JSON(500, gin.H{"message": "Logging in failed, please try again later."})
		return
	}

	c.JSON(200, models.UserResponse{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Token:     tokenString,
	})
}
