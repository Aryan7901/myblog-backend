package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName string               `json:"firstName" bson:"firstName" binding:"required"`
	LastName  string               `json:"lastName" bson:"lastName" binding:"required"`
	Email     string               `json:"email" bson:"email" binding:"required,email"`
	Password  string               `json:"password" bson:"password" binding:"required,min=8"`
	Blogs     []primitive.ObjectID `json:"blogs" bson:"blogs"`
}

type UserResponse struct {
	ID        string `json:"user"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Token     string `json:"token,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SignupRequest struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
}

type Comment struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	User    primitive.ObjectID `json:"user" bson:"user" binding:"required"`
	Content string             `json:"content" bson:"content" binding:"required"`
	Date    time.Time          `json:"date" bson:"date" binding:"required"`
	Blog    primitive.ObjectID `json:"blog" bson:"blog" binding:"required"`
}

type Blog struct {
	ID          primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string               `json:"title" bson:"title" binding:"required"`
	Author      primitive.ObjectID   `json:"author" bson:"author" binding:"required"`
	Description string               `json:"description" bson:"description" binding:"required"`
	Article     string               `json:"article" bson:"article" binding:"required,min=500"`
	Comments    []primitive.ObjectID `json:"comments" bson:"comments"`
}

type BlogRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Article     string `json:"article" binding:"required,min=500"`
}

type CommentRequest struct {
	Comment string `json:"comment" binding:"required"`
}

type Author struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
}

type BlogResponse struct {
	ID          string            `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string            `json:"title"`
	Author      Author            `json:"author"`
	Description string            `json:"description"`
	Article     string            `json:"article"`
	Comments    []CommentResponse `json:"comments"`
}

type CommentResponse struct {
	ID      string    `json:"_id,omitempty" bson:"_id,omitempty"`
	User    Author    `json:"user"`
	Content string    `json:"content"`
	Date    time.Time `json:"date"`
	Blog    string    `json:"blog,omitempty"`
}

type UserBlogResponse struct {
	Blogs  []Blog `json:"blogs"`
	Author Author `json:"author"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
