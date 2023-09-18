package main

import (
	// "os/user"
	// "fmt"
	"log"
	"net/http"
	"reelState/controllers"
	"reelState/middlewares"
	"reelState/models"
	// "time"
	"github.com/gin-gonic/gin"
	// "github.com/ianlopshire/go-async"

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


	public.GET("/getFeedVideos", controllers.HandleGetAllVideos)
	public.GET("/getFeedCategoryVideos", controllers.HandleGetAllCategoriesVideos)
	public.GET("/getSearchVideos", controllers.HandleSearchVideos)
	public.GET("/UserByIdHandler/", controllers.UserByIdHandler)
	public.GET("/getUserVideos", controllers.GetUserVideos)

	public.POST("/sendVerificationCode", controllers.SendVerificationCode)
	public.POST("/CodeIsValid", controllers.ValidateVerificationCode)
	public.POST("/ValidateUserName", controllers.ValidateUserName)
	public.POST("/register", controllers.RegisterHandler)
	public.POST("/login", controllers.LoginHandler)
	
	
	
	protected := r.Group("/api/admin")
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.POST("/upload", controllers.HandleVideoUpload)
	protected.POST("/UpdateProfileImageUserName", controllers.UpdateProfileImageUserName)
	protected.POST("/SetFavorite", controllers.SetFavorite)

	// protected.GET("/user", controllers.CurrentUserHandler)
	protected.GET("/getMyVideos", controllers.GetMyVideos)
	protected.GET("/getMyFavoritesVideos", controllers.GetUserFavoritesVideos)

	protected.GET("/user", controllers.CurrentUserHandler)
	protected.GET("/getCategoriesAndTypes", controllers.HandleGetCategoriesAndTypes)
	protected.GET("/getAroundVideos", controllers.HandleGetAroundVideos)

	type User struct {
		ID   int
		Name string
	}
	
	
	// output: {1 John Does} <nil>

	log.Fatal(http.ListenAndServe(":8080", r))

}


// func DoneAsync() int {
// 	fmt.Println("Warming up ...")
// 	time.Sleep(3 * time.Second)
// 	fmt.Println("Done ...")
// 	return 1
// }
