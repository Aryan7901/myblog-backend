package controllers

import (
	"context"
	"time"

	"backend/config"
	"backend/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateBlog(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var blogReq models.BlogRequest
	if err := c.ShouldBindJSON(&blogReq); err != nil {
		c.JSON(422, gin.H{"message": "Invalid inputs passed, please check your data."})
		return
	}

	userData := c.MustGet("userData").(map[string]string)
	uid, _ := primitive.ObjectIDFromHex(userData["userId"])

	blog := models.Blog{
		Title:       blogReq.Title,
		Description: blogReq.Description,
		Article:     blogReq.Article,
		Author:      uid,
		Comments:    []primitive.ObjectID{},
	}

	result, err := config.BlogCollection.InsertOne(ctx, blog)
	if err != nil {
		c.JSON(500, gin.H{"message": "Creating new blog failed, please try again later."})
		return
	}

	_, err = config.UserCollection.UpdateOne(ctx, bson.M{"_id": uid}, bson.M{"$push": bson.M{"blogs": result.InsertedID}})
	if err != nil {
		c.JSON(500, gin.H{"message": "Creating new blog failed, please try again later."})
		return
	}

	c.JSON(201, gin.H{"createdBlog": blog})
}

func UpdateBlog(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	blogId := c.Param("bid")
	bid, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		c.JSON(500, gin.H{"message": "Invalid blog ID"})
		return
	}

	var blogReq models.BlogRequest
	if err := c.ShouldBindJSON(&blogReq); err != nil {
		c.JSON(422, gin.H{"message": "Invalid inputs passed, please check your data."})
		return
	}

	userData := c.MustGet("userData").(map[string]string)
	uid, _ := primitive.ObjectIDFromHex(userData["userId"])

	var blog models.Blog
	err = config.BlogCollection.FindOne(ctx, bson.M{"_id": bid}).Decode(&blog)
	if err != nil {
		c.JSON(500, gin.H{"message": "Updating blog failed, please try again later."})
		return
	}

	if blog.Author != uid {
		c.JSON(401, gin.H{"message": "Unauthorized!"})
		return
	}

	_, err = config.BlogCollection.UpdateOne(ctx, bson.M{"_id": bid}, bson.M{"$set": bson.M{
		"title":       blogReq.Title,
		"description": blogReq.Description,
		"article":     blogReq.Article,
	}})
	if err != nil {
		c.JSON(500, gin.H{"message": "Updating blog failed, please try again later."})
		return
	}

	c.JSON(200, gin.H{"message": "Blog updated!"})
}

func DeleteBlog(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	blogId := c.Param("bid")
	bid, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		c.JSON(500, gin.H{"message": "Invalid blog ID"})
		return
	}

	userData := c.MustGet("userData").(map[string]string)
	uid, _ := primitive.ObjectIDFromHex(userData["userId"])

	var blog models.Blog
	err = config.BlogCollection.FindOne(ctx, bson.M{"_id": bid}).Decode(&blog)
	if err != nil {
		c.JSON(500, gin.H{"message": "Deleting blog failed, please try again later."})
		return
	}

	if blog.Author != uid {
		c.JSON(401, gin.H{"message": "Unauthorized!"})
		return
	}

	_, err = config.BlogCollection.DeleteOne(ctx, bson.M{"_id": bid})
	if err != nil {
		c.JSON(500, gin.H{"message": "Deleting blog failed, please try again later."})
		return
	}

	_, err = config.UserCollection.UpdateOne(ctx, bson.M{"_id": uid}, bson.M{"$pull": bson.M{"blogs": bid}})
	if err != nil {
		c.JSON(500, gin.H{"message": "Deleting blog failed, please try again later."})
		return
	}

	c.JSON(200, gin.H{"message": "Blog deleted!"})
}

func GetUserBlogs(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userData := c.MustGet("userData").(map[string]string)
	uid, _ := primitive.ObjectIDFromHex(userData["userId"])
	firstName := userData["firstName"]
	lastName := userData["lastName"]

	cursor, err := config.BlogCollection.Find(ctx, bson.M{"author": uid})
	if err != nil {
		c.JSON(500, gin.H{"message": "failed to get user's blogs"})
		return
	}
	defer cursor.Close(ctx)

	var blogs []models.Blog = make([]models.Blog, 0)
	if err := cursor.All(ctx, &blogs); err != nil {
		c.JSON(500, gin.H{"message": "failed to get user's blogs"})
		return
	}

	c.JSON(200, models.UserBlogResponse{
		Blogs: blogs,
		Author: models.Author{
			FirstName: firstName,
			LastName:  lastName,
		},
	})
}
