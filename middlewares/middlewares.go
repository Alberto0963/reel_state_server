package middlewares

import (
	"net/http"
	// "os"
	"path/filepath"
	"strings"

	"reelState/utils/token"

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
			allowedExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".mp4", ".avi", ".mkv"}
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

