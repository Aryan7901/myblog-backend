package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusForbidden, gin.H{"message": "Authentication failed!"})
			c.Abort()
			return
		}

		tokenString := strings.Split(authHeader, " ")[1]
		if tokenString == "" {
			c.JSON(http.StatusForbidden, gin.H{"message": "Authentication failed!"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("TOKEN_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusForbidden, gin.H{"message": "Authentication failed!"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"message": "Authentication failed!"})
			c.Abort()
			return
		}

		userId, _ := claims["userId"].(string)
		email, _ := claims["email"].(string)
		firstName, _ := claims["firstName"].(string)
		lastName, _ := claims["lastName"].(string)

		c.Set("userData", map[string]string{
			"userId":    userId,
			"email":     email,
			"firstName": firstName,
			"lastName":  lastName,
		})

		c.Next()
	}
}
