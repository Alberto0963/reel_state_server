package main

import (
	// "os/user"
	// "fmt"
	"log"
	"net/http"
	"os"
	"reelState/controllers"
	"reelState/middlewares"
	"reelState/models"

	// "time"
	"github.com/gin-gonic/gin"
	// "github.com/ianlopshire/go-async"
)

const sampleRate = 44100
const seconds = 2

func main() {

	gin.SetMode(gin.DebugMode)

	models.InitDB()

	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.Use(middlewares.BlockFolderAccessMiddleware())
	// Specify the directory containing your public files
	publicDir :=  os.Getenv("MY_URL") +"./public"

	// Create a file server handler for the public directory
	fs := http.FileSystem(http.Dir(publicDir))

	// Serve static files from the "public" directory
	r.StaticFS("/public", fs)
	// Specify the directory containing your public files

	// Create a file server handler for the public directory

	// Register the file server handler with a specific URL path
	// Specify the directory containing your public files

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
	public.POST("/UpdatePasswordHandler", controllers.UpdatePasswordHandler)
	public.GET("/getAroundVideos", controllers.HandleGetAroundVideos)
	

	protected := r.Group("/api/admin")
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.POST("/upload", controllers.HandleVideoUpload)
	protected.POST("/edit", controllers.HandleVideoEdit)
	protected.POST("/updateusername", controllers.UpdateUsernameHandler)
	protected.POST("/uploadVideoWithAudioUpload", controllers.HandleVideoWithAudioUpload)

	protected.POST("/UpdateProfileImageUserName", controllers.UpdateProfileImageUserName)
	protected.POST("/UpdateCoverImageUserName", controllers.UpdateCoverImageUserName)
	protected.POST("/SetFavorite", controllers.SetFavorite)
	protected.POST("/DeleteUserVideo", controllers.DeleteUserVideo)

	// protected.GET("/user", controllers.CurrentUserHandler)
	protected.GET("/getMyVideos", controllers.GetMyVideos)
	protected.GET("/getMyFavoritesVideos", controllers.GetUserFavoritesVideos)
	protected.GET("/GetMemberShips", controllers.GetMemberShips)
	protected.GET("/user", controllers.CurrentUserHandler)
	protected.GET("/getCategoriesAndTypes", controllers.HandleGetCategoriesAndTypes)
	protected.GET("/getsongs", controllers.HandleGetAllSongs)

	// type User struct {
	// 	ID   int
	// 	Name string
	// }

	// output: {1 John Does} <nil>

	
	log.Fatal(http.ListenAndServe(":8080", r))

}

// func DoneAsync() int {
// 	fmt.Println("Warming up ...")
// 	time.Sleep(3 * time.Second)
// 	fmt.Println("Done ...")
// 	return 1
// }
