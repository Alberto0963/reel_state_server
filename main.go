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

	// SMS "reelState/utils"

	// "github.com/joho/godotenv"
	"github.com/jrallison/go-workers"
	// "golang.org/x/net/websocket"

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
	// Apply CORS middleware to all routes
	r.Use(middlewares.CORSMiddleware())
	// Specify the directory containing your public files
	publicDir := os.Getenv("MY_URL") + "./public"

	// Create a file server handler for the public directory
	fs := http.FileSystem(http.Dir(publicDir))

	// Serve static files from the "public" directory
	r.StaticFS("/public", fs)
	// r.LoadHTMLGlob("templates/*")


	// Manejador para la ruta raíz
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Bienvenido a ReelState")
	})

	// Ruta para servir el archivo JSON
	r.GET("/.well-known/assetlinks.json", func(c *gin.Context) {
		c.File(publicDir+"/.well-known/assetlinks.json") // Especifica la ruta del archivo JSON
	})

	public := r.Group("/api")
	webhooks := r.Group("/webhooks")
	// r.LoadHTMLGlob(os.Getenv("MY_URL") + "templates/*")
	r.GET("/video/:videoID", controllers.GetVideoFromLink)
	// public.GET("/getFeedVideos", controllers.HandleGetAllVideos)
	public.GET("/GetMemberShips", controllers.GetPublicMemberShips)

	public.GET("/getFeedVideos", controllers.HandleGetAllVideos)
	public.GET("/getSearchVideos", controllers.HandleSearchVideos)
	public.GET("/UserByIdHandler/", controllers.UserByIdHandler)
	public.GET("/getUserVideos", controllers.GetUserVideos)
	public.GET("/getReportsTypes", controllers.HandleGetTypeRepors)
	public.GET("/searchProfile", controllers.SearchProfile)
	public.POST("/loginWithGoogle", controllers.LoginWithGoogleHandler)
	public.POST("/validatePhone", controllers.ValidatePhone)

	public.POST("/sendVerificationCode", controllers.SendVerificationCode)
	public.POST("/CodeIsValid", controllers.ValidateVerificationCode)
	public.POST("/ValidateUserName", controllers.ValidateUserName)
	public.POST("/register", controllers.RegisterHandler)
	public.POST("/registerWithGoogle", controllers.RegisterHandlerWithGogle)
	public.POST("/registerWithFacebook", controllers.RegisterHandlerWithFacebook)
	public.POST("/loginWithFacebook", controllers.LoginHandlerWithFacebook)

	
	
	public.POST("/addUserSubscription", controllers.HandleWebhook)
	public.POST("/openpayWebhook", controllers.OpenpayWebhookHandler)
	webhooks.POST("/webhookWhatsapp", controllers.WebhookHandlerWhatsapp)

	public.POST("/login", controllers.LoginHandler)
	public.POST("/UpdatePasswordHandler", controllers.UpdatePasswordHandler)
	public.GET("/getAroundVideos", controllers.HandleGetAroundVideos)
	public.GET("/getReviewsByProfile", controllers.GetReviewsByProfile)

	protected := r.Group("/api/admin")
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.POST("/upload", controllers.HandleVideoUpload)
	protected.POST("/edit", controllers.HandleVideoEdit)
	protected.POST("/updateusername", controllers.UpdateUsernameHandler)
	protected.POST("/deleteUser", controllers.DeleteUser)
	protected.POST("/makePay_openPay", controllers.MakePayOpenPay)

	// protected.POST("/uploadVideoWithAudioUpload", controllers.HandleVideoWithAudioUpload)

	protected.POST("/UpdateProfileImageUserName", controllers.UpdateProfileImageUserName)
	protected.POST("/updateUserName", controllers.UpdateUsernameHandler)
	protected.POST("/updatePhoneNumber", controllers.UpdatePhoneNumberHandler)

	// protected.POST("/changePassword", controllers.ChangePassword)

	protected.POST("/UpdateCoverImageUserName", controllers.UpdateCoverImageUserName)
	protected.POST("/SetFavorite", controllers.SetFavorite)
	protected.POST("/DeleteUserVideo", controllers.DeleteUserVideo)
	protected.POST("/setProfileLike", controllers.Setlike)

	protected.POST("/createSubscription", controllers.CreateSubscription)
	protected.POST("/cancelSubscription", controllers.CancelSubscription)
	protected.POST("/updateLinkUserName", controllers.UpdateLinkUserName)
	protected.POST("/createOrUpdateReview", controllers.CreateOrUpdateReview)
	protected.POST("/addOrUpdateSales", controllers.AddOrUpdateSales)
	protected.POST("/updateDeviceToken", controllers.UpdateDeviceToken)
	protected.POST("/deleteReview", controllers.DeleteReview)
	protected.POST("/ai", controllers.GenerateResponse)


	protected.POST("/createSponsor", controllers.CreateSponsor)
	protected.POST("/editSponsor", controllers.EditSponsor)
	protected.POST("/deleteSponsor", controllers.DeleteSponsor)

	// protected.GET("/user", controllers.CurrentUserHandler)
	protected.GET("/getMyVideos", controllers.GetMyVideos)
	protected.GET("/getVideoSponsors", controllers.HandleGetVideosSponsors)

	protected.GET("/getMyFavoritesVideos", controllers.GetUserFavoritesVideos)
	protected.GET("/GetMCancelSubscriptionemberShips", controllers.GetMemberShips)
	protected.GET("/user", controllers.CurrentUserHandler)
	protected.GET("/getCategoriesAndTypes", controllers.HandleGetCategoriesAndTypes)
	protected.GET("/getsongs", controllers.HandleGetAllSongs)
	protected.GET("/status/:taskId", controllers.CheckStatus)
	protected.GET("/getUserSubscription", controllers.GetUserSubscription)


	// workers.Run()
	log.Fatal(http.ListenAndServe(":8080", r))

}

// func DoneAsync() int {
// 	fmt.Println("Warming up ...")
// 	time.Sleep(3 * time.Second)
// 	fmt.Println("Done ...")
// 	return 1
// }
