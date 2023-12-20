package main

import (
	// "os/user"
	// "fmt"
	"fmt"
	"log"
	"net/http"
	"os"
	"reelState/controllers"
	"reelState/middlewares"
	"reelState/models"

	"github.com/jrallison/go-workers"

	// "time"
	"github.com/gin-gonic/gin"
	// "github.com/ianlopshire/go-async"
)

const sampleRate = 44100
const seconds = 2

func MyBackgroundTask(msg *workers.Msg) {
	// Perform background task
	fmt.Println("Background task is running...")
}

type myMiddleware struct{}

func (r *myMiddleware) Call(queue string, message *workers.Msg, next func() bool) (acknowledge bool) {
	// do something before each message is processed
	fmt.Println("Procesando Video")
	acknowledge = next()
	// do something after each message is processed

	fmt.Println("Video Procesado")
	return
}

func main() {

	// mode := os.Getenv("GIN_MODE")

	models.InitDB()

	// Start the worker process
	workers.Configure(map[string]string{
		// location of redis instance
		"server": "localhost:6379",
		// instance of the database
		"database": "0",
		// number of connections to keep open with redis
		"pool": "30",
		// unique process id for this instance of workers (for proper recovery of inprogress jobs on crash)
		"process": "1",
	})

	workers.Middleware.Append(&myMiddleware{})

	r := gin.Default()

	r.SetTrustedProxies(nil)
	r.Use(middlewares.BlockFolderAccessMiddleware())
	// Specify the directory containing your public files
	publicDir := os.Getenv("MY_URL") + "./public"

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
	// protected.POST("/uploadVideoWithAudioUpload", controllers.HandleVideoWithAudioUpload)

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

	// // pull messages from "myqueue" with concurrency of 10
	// workers.Process("myqueue", myJob, 10)

	// // pull messages from "myqueue2" with concurrency of 20
	// workers.Process("myqueue2", myJob, 20)
	// stats will be available at http://localhost:8080/stats
	// pull messages from "myqueue" with concurrency of 10
	// workers.Process("myqueue", myJob, 10)
	// Enqueue the background task
	    // Register the background task
    // workers.Process("myqueue", MyBackgroundTask,10)
	// workers.Enqueue("myqueue", "MyBackgroundTask", nil)
	// // Register the background task
	workers.Process("myqueue", MyBackgroundTask,1)
	go workers.StatsServer(8081)
	// Blocks until process is told to exit via unix signal

	// workers.Run()
	log.Fatal(http.ListenAndServe(":8080", r))

}

// func DoneAsync() int {
// 	fmt.Println("Warming up ...")
// 	time.Sleep(3 * time.Second)
// 	fmt.Println("Done ...")
// 	return 1
// }
