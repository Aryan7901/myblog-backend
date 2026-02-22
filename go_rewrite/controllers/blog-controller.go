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

func GetAllBlogs(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "author",
				"foreignField": "_id",
				"as":           "authorData",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$authorData",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$project": bson.M{
				"_id":         1,
				"title":       1,
				"description": 1,
				"article":     1,
				"author": bson.M{
					"firstName": "$authorData.firstName",
					"lastName":  "$authorData.lastName",
				},
			},
		},
	}

	cursor, err := config.BlogCollection.Aggregate(ctx, pipeline)
	if err != nil {
		c.JSON(500, gin.H{"message": "Error Retrieving data, please try again later."})
		return
	}
	defer cursor.Close(ctx)

	var blogResponses []models.BlogResponse = make([]models.BlogResponse, 0)
	if err := cursor.All(ctx, &blogResponses); err != nil {
		c.JSON(500, gin.H{"message": "Error Retrieving data, please try again later."})
		return
	}

	c.JSON(200, gin.H{"blogs": blogResponses})
}

func GetBlogById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	blogId := c.Param("bid")
	oid, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid blog ID"})
		return
	}

	var blog models.Blog
	err = config.BlogCollection.FindOne(ctx, bson.M{"_id": oid}).Decode(&blog)
	if err != nil {
		c.JSON(500, gin.H{"message": "Error Retrieving blog, please try again later."})
		return
	}

	var author models.User
	_ = config.UserCollection.FindOne(ctx, bson.M{"_id": blog.Author}).Decode(&author)

	comments := make([]models.CommentResponse, 0)
	if len(blog.Comments) > 0 {
		commentCursor, err := config.CommentCollection.Find(ctx, bson.M{"_id": bson.M{"$in": blog.Comments}})
		if err == nil {
			var commentDocs []models.Comment
			commentCursor.All(ctx, &commentDocs)

			userIDs := make([]primitive.ObjectID, 0, len(commentDocs))
			for _, c := range commentDocs {
				userIDs = append(userIDs, c.User)
			}

			var users []models.User
			if len(userIDs) > 0 {
				userCursor, _ := config.UserCollection.Find(ctx, bson.M{"_id": bson.M{"$in": userIDs}})
				userCursor.All(ctx, &users)
			}

			userMap := make(map[primitive.ObjectID]models.User)
			for _, u := range users {
				userMap[u.ID] = u
			}

			for _, comment := range commentDocs {
				user := userMap[comment.User]
				comments = append(comments, models.CommentResponse{
					ID:      comment.ID.Hex(),
					Content: comment.Content,
					Date:    comment.Date,
					User:    models.Author{FirstName: user.FirstName, LastName: user.LastName},
				})
			}
		}
	}

	c.JSON(200, gin.H{"blog": models.BlogResponse{
		ID:          blog.ID.Hex(),
		Title:       blog.Title,
		Author:      models.Author{FirstName: author.FirstName, LastName: author.LastName},
		Description: blog.Description,
		Article:     blog.Article,
		Comments:    comments,
	}})
}

func MakeComment(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	blogId := c.Param("bid")
	bid, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid blog ID"})
		return
	}

	var commentReq models.CommentRequest
	if err := c.ShouldBindJSON(&commentReq); err != nil {
		c.JSON(422, gin.H{"message": "Invalid inputs passed, please check your data."})
		return
	}

	userData := c.MustGet("userData").(map[string]string)
	uid, _ := primitive.ObjectIDFromHex(userData["userId"])

	var blog models.Blog
	err = config.BlogCollection.FindOne(ctx, bson.M{"_id": bid}).Decode(&blog)
	if err != nil {
		c.JSON(500, gin.H{"message": "Could not find blog, please try again later."})
		return
	}

	comment := models.Comment{
		User:    uid,
		Content: commentReq.Comment,
		Date:    time.Now(),
		Blog:    bid,
	}

	result, err := config.CommentCollection.InsertOne(ctx, comment)
	if err != nil {
		c.JSON(500, gin.H{"message": "Adding comment failed, please try again later."})
		return
	}

	cid := result.InsertedID.(primitive.ObjectID)
	_, err = config.BlogCollection.UpdateOne(ctx, bson.M{"_id": bid}, bson.M{"$push": bson.M{"comments": cid}})
	if err != nil {
		c.JSON(500, gin.H{"message": "Adding comment failed, please try again later."})
		return
	}

	c.JSON(201, gin.H{"message": "Comment Created!"})
}

func UpdateComment(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	commentId := c.Param("cid")
	cid, err := primitive.ObjectIDFromHex(commentId)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid comment ID"})
		return
	}

	var commentReq models.CommentRequest
	if err := c.ShouldBindJSON(&commentReq); err != nil {
		c.JSON(422, gin.H{"message": "Invalid inputs passed, please check your data."})
		return
	}

	userData := c.MustGet("userData").(map[string]string)
	uid, _ := primitive.ObjectIDFromHex(userData["userId"])

	var comment models.Comment
	err = config.CommentCollection.FindOne(ctx, bson.M{"_id": cid}).Decode(&comment)
	if err != nil {
		c.JSON(500, gin.H{"message": "Updating comment failed, please try again later."})
		return
	}

	if comment.User != uid {
		c.JSON(401, gin.H{"message": "Unauthorized!"})
		return
	}

	_, err = config.CommentCollection.UpdateOne(ctx, bson.M{"_id": cid}, bson.M{"$set": bson.M{"content": commentReq.Comment}})
	if err != nil {
		c.JSON(500, gin.H{"message": "Updating comment failed, please try again later."})
		return
	}

	c.JSON(200, gin.H{"message": "Comment updated!"})
}

func DeleteComment(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	commentId := c.Param("cid")
	cid, err := primitive.ObjectIDFromHex(commentId)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid comment ID"})
		return
	}

	userData := c.MustGet("userData").(map[string]string)
	uid, _ := primitive.ObjectIDFromHex(userData["userId"])

	var comment models.Comment
	err = config.CommentCollection.FindOne(ctx, bson.M{"_id": cid}).Decode(&comment)
	if err != nil {
		c.JSON(500, gin.H{"message": "Deleting comment failed, please try again later."})
		return
	}

	if comment.User != uid {
		c.JSON(401, gin.H{"message": "Unauthorized!"})
		return
	}

	_, err = config.CommentCollection.DeleteOne(ctx, bson.M{"_id": cid})
	if err != nil {
		c.JSON(500, gin.H{"message": "Deleting comment failed, please try again later."})
		return
	}

	_, err = config.BlogCollection.UpdateOne(ctx, bson.M{"comments": cid}, bson.M{"$pull": bson.M{"comments": cid}})
	if err != nil {
		c.JSON(500, gin.H{"message": "Deleting comment failed, please try again later."})
		return
	}

	c.JSON(200, gin.H{"message": "Comment deleted!"})
}
