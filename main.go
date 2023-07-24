package main

import (
	// "os/user"
	"log"
	"net/http"
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
	r.Static("/public", "home/reelstate/go/reel_state_server/public")

	public := r.Group("/api")

	public.POST("/register", controllers.RegisterHandler)
	public.POST("/login", controllers.LoginHandler)
	public.GET("/getFeedVideos", controllers.HandleGetAllVideos)

	protected := r.Group("/api/admin")
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.GET("/user", controllers.CurrentUserHandler)
	protected.POST("/upload", controllers.HandleVideoUpload)
	protected.GET("/getCategoriesAndTypes", controllers.HandleGetCategoriesAndTypes)

	log.Fatal(http.ListenAndServe(":8080", r))

}
