package config

import (
	"context"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client
var UserCollection *mongo.Collection
var BlogCollection *mongo.Collection
var CommentCollection *mongo.Collection

func getDatabaseName(uri string) string {
	if idx := strings.LastIndex(uri, "/"); idx != -1 && idx+1 < len(uri) {
		remaining := uri[idx+1:]
		if !strings.HasPrefix(remaining, "?") && !strings.HasPrefix(remaining, "/") {
			if endIdx := strings.IndexAny(remaining, "?/"); endIdx != -1 {
				return remaining[:endIdx]
			}
			return remaining
		}
	}
	return "myblog"
}

func ConnectDB(uri string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	DB = client
	database := getDatabaseName(uri)
	UserCollection = DB.Database(database).Collection("users")
	BlogCollection = DB.Database(database).Collection("blogs")
	CommentCollection = DB.Database(database).Collection("comments")
	log.Printf("Connected to database: %s", database)
	return nil
}

func GetCollection(database, collection string) *mongo.Collection {
	return DB.Database(database).Collection(collection)
}
