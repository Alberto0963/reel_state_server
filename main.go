package main

import (
	// "os/user"
	"reelState/controllers"
	"reelState/middlewares"
	"reelState/models"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)

	models.InitDB()

	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.Use(middlewares.BlockFolderAccessMiddleware())

	// Serve static files from the "public" directory
	r.Static("/public", "./public")

	public := r.Group("/api")

	public.POST("/register", controllers.RegisterHandler)
	public.POST("/login", controllers.LoginHandler)

	protected := r.Group("/api/admin")
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.GET("/user", controllers.CurrentUserHandler)
	protected.POST("/upload", controllers.HandleVideoUpload)

	r.Run(":8080")

}
