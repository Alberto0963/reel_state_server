package middlewares

import (
	"net/http"
	// "os"
	"path/filepath"
	"strings"

	"reelState/utils/token"

	// "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := token.TokenValid(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}

// func JwtAuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
// 			c.Abort()
// 			return
// 		}

// 		tokenString := strings.Split(authHeader, "Bearer ")[1]
// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			// Replace "your-secret-key" with your actual secret key used for signing JWT tokens
// 			return []byte("your-secret-key"), nil
// 		})

// 		if err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
// 			c.Abort()
// 			return
// 		}

// 		claims, ok := token.Claims.(jwt.MapClaims)
// 		if !ok || !token.Valid {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
// 			c.Abort()
// 			return
// 		}

// 		// Extract the user ID or other relevant information from the token claims
// 		userID, ok := claims["user_id"].(string)
// 		if !ok {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
// 			c.Abort()
// 			return
// 		}

// 		// Attach the user ID or other relevant information to the request context or context-specific variables
// 		c.Set("user_id", userID)

// 		c.Next()
// 	}
// }


// BlockFolderAccessMiddleware is a middleware function to block folder access and only allow access to specific files.
func BlockFolderAccessMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the requested URL path
		requestedPath := c.Request.URL.Path

		// Check if the requested path starts with "/public/"
		if strings.HasPrefix(requestedPath, "/public/") {
			// Extract the file extension from the requested path
			fileExtension := filepath.Ext(requestedPath)

			// Check if the file extension is allowed (e.g., images and videos)
			allowedExtensions := []string{".jpg",".json", ".jpeg", ".png", ".gif", ".mp4", ".avi", ".mkv",".mp3"}
			isAllowedExtension := false
			for _, ext := range allowedExtensions {
				if strings.EqualFold(fileExtension, ext) {
					isAllowedExtension = true
					break
				}
			}

			// If the file extension is not allowed, return a 404 Not Found error
			if !isAllowedExtension {
				c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
				c.Abort()
				return
			}
		}

		// Continue to the next middleware or route handler
		c.Next()
	}
}

