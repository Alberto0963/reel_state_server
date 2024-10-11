package controllers

import (
	// "io"
	// "mime/multipart"
	// "go/token"
	// "fmt"
	"bytes"
	"encoding/json"
	"sync"

	// "os/exec"

	// "time"

	// "errors"
	// "strings"

	// "encoding/json"
	// "encoding/json"
	// "io/ioutil"

	// "encoding/json"
	"fmt"

	"io"

	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	// SMS "reelState/utils"
	"reelState/utils/token"
	"strconv"

	// "os"

	// "reelState/auth"
	"reelState/models"

	"github.com/gin-gonic/gin"
	"github.com/ianlopshire/go-async"
	// "github.com/jrallison/go-workers"
	// "github.com/jrallison/go-workers"
	// "github.com/jrallison/go-workers"
	// "golang.org/x/crypto/nacl/auth"
)

// Define a global map to track task statuses
var (
	taskStatusMap = make(map[string]string)
	mutex         = &sync.Mutex{}
)

type RegisterFavInput struct {
	// Id_user int `json:"id_user" binding:"required"`
	Id_video int `json:"id_video" binding:"required"`
	// Phone    string `json:"phone" binding:"required"`
	// Code     string `json:"phone" binding:"required"`
}

type VideoInput struct {
	// Username string `json:"username" binding:"required"`
	// Password string `json:"password" binding:"required"`
	// ID uint `json:"id" binding:"required" `
	// Video_url string `json:"video_url" binding:"required"`
	Description     string `json:"description" binding:"required"`
	Location        string `json:"location" binding:"required"`
	Area            string `json:"area" binding:"required"`
	Property_number string `json:"property_number" binding:"required"`
	Price           string `json:"price" binding:"required"`
	// Id_user string `json:"id_user" binding:"required"`
	Sale_type_id     string  `json:"sale_type_id" binding:"required"`
	Sale_category_id string  `json:"sale_category_id" binding:"required"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	Type             int     `json:"type"`
}

func GetVideoFromLink(c *gin.Context) {
	userID, _ := token.ExtractTokenID(c)
	var data models.FeedVideo
	videoID := c.Param("videoID")
	id_vid, err := strconv.ParseUint(videoID, 10, 64)

	data, err = models.GetVideo(int(id_vid), int(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	videoURL := os.Getenv("App_URL")
	// videoURL := "https://api.reelstate.mx/" + data.Video_url
	// img := "https://api.reelstate.mx/" + data.Image_cover
	c.HTML(200, "newPlayVideo.html", gin.H{
		"title":       "Video Showcase",
		"VideoURL":    videoURL + data.Video_url,
		"User":        data.User.Username,
		"Description": data.Description,
		"Img":         videoURL + data.Image_cover,
	})
}

func checkVideoSize(filePath string, maxSize int64) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	fileSize := fileInfo.Size()
	if fileSize > maxSize {
		return fmt.Errorf("Error: File size exceeds the maximum limit of %d MB ", maxSize/(1024*1024))
	}

	return nil
}

// ErrorResponse represents a JSON error response structure
type ErrorResponse struct {
	Error string `json:"error"`
}

// func HandleVideoUpload(c *gin.Context) {

// 	// c.Request.Body = http.MaxBytesReader(c.Request.Response., c.Request.Body, 300*1024*1024)

// 	file, err := c.FormFile("video")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	userID, _ := token.ExtractTokenID(c)

// 	err = models.ValidateUserType(int(userID))
// 	if err != nil {
// 		c.JSON(http.StatusForbidden, gin.H{"error": "no puedes publicar mas videos"})
// 		return
// 	}
// 	audioFileName := c.PostForm("audio")
// 	url := os.Getenv("MY_URL")

// 	// Generate a random file name
// 	fileName := models.GenerateRandomName()
// 	tempFilePath := filepath.Join(url, "/public/videos", fileName+filepath.Ext(file.Filename))
// 	finalVideoName := models.GenerateRandomName()
// 	finalVideoPath := filepath.Join(url, "/public/videos", finalVideoName+filepath.Ext(file.Filename))

// 	// Create the destination file
// 	//destPath := filepath.Join("", fileName)
// 	// baseDir, err := os.Getwd() // Get the current working directory

// 	// if url != nil {

// 	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 	// 	return
// 	// }

// 	// Get the file size from the request
// 	fileSize := c.Request.ContentLength
// 	if fileSize <= 0 || fileSize > (300*1024*1024) {
// 		// http.Error(w, "Error: Invalid file size", http.StatusBadRequest)
// 		errorResponse := ErrorResponse{Error: fmt.Sprintf("File size exceeds the maximum limit of %d MB", (300*1024*1024)/(1024*1024))}
// 		c.JSON(http.StatusRequestEntityTooLarge, errorResponse)

// 		return
// 	}

// 	// user, _ := models.GetUserByIDWithVideos(userID)
// 	// countuservideos := len(user.Videos)
// 	// if user.Id_Membership == 1 && countuservideos >= 1 {
// 	// 	// http.Error(w, "Error: Invalid file size", http.StatusBadRequest)
// 	// 	errorResponse := ErrorResponse{Error: fmt.Sprintf("exceeds the maximum limit of Videos")}
// 	// 	c.JSON(http.StatusPreconditionFailed, errorResponse)
// 	// 	return
// 	// }

// 	saveVideo := new(async.Future[error])
// 	// destPath := filepath.Join(url, "/public/videos", fileName+filepath.Ext(file.Filename))

// 	// save video.

// 	go func() {
// 		// err = saveVideoFile(file, destPath)
// 		fmt.Println("/////////////// inicio ///////////////////")
// 		// time.Sleep(50 * time.Second)

// 		async.ResolveFuture(saveVideo, saveVideoFile(file, tempFilePath), nil)

// 	}()

// 	async.Await(saveVideo)

// 	fmt.Println("/////////////// final ///////////////////")
// 	d, err := os.Stat(tempFilePath)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	fmt.Println(d)
// 	//// end save video

// 	//add audio to video
// 	if audioFileName != "" {
// 		destAudioPath := filepath.Join(url, "/public/audio", audioFileName)
// 		fileName = models.GenerateRandomName()

// 		workers.Enqueue("myqueue", "Add", joinAudioWithVideo(destAudioPath, tempFilePath, fileName+filepath.Ext(file.Filename)))
// 		tempFilePath = filepath.Join(url, "/public/videos", fileName+filepath.Ext(file.Filename))

// 	}

// 	workers.Enqueue("myqueue", "Add", getFrame(tempFilePath, fileName+".jpg"))

// 	// Add a job to a queue
// 	workers.Enqueue("myqueue", "Add", compressVideo(tempFilePath, finalVideoPath))

// 	var input VideoInput

// 	if err := c.ShouldBind(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	v := models.Video{}
// 	v.Video_url = filepath.Join("/public/videos", finalVideoName+filepath.Ext(file.Filename))
// 	v.Image_cover = "public/video_cover/" + fileName + ".jpg"
// 	v.Description = input.Description
// 	v.Location = input.Location
// 	v.Area = input.Area
// 	v.Property_number = input.Property_number
// 	v.Price = input.Price
// 	v.Id_user = userID
// 	v.Latitude = input.Latitude
// 	v.Longitude = input.Longitude
// 	v.Type = input.Type

// 	sale_type_id, err := strconv.ParseUint(input.Sale_type_id, 10, 32)
// 	if err != nil {
// 		// Handle the error if the conversion fails
// 		fmt.Println("Error converting string to uint:", err)
// 		return
// 	}
// 	v.Sale_type_id = int(sale_type_id)
// 	sale_category_id, err := strconv.ParseUint(input.Sale_category_id, 10, 32)
// 	if err != nil {
// 		// Handle the error if the conversion fails
// 		fmt.Println("Error converting string to uint:", err)
// 		return
// 	}
// 	v.Sale_category_id = int(sale_category_id)
// 	_, err = v.SaveVideo()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Video uploaded successfully"})
// }

func HandleVideoUpload(c *gin.Context) {
	file, _ := c.FormFile("video")
	// Assuming the file is saved successfully in your server's file system

	////////////////////////////
	// c.Request.Body = http.MaxBytesReader(c.Request.Response., c.Request.Body, 300*1024*1024)

	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _ := token.ExtractTokenID(c)

	err = models.ValidateUserType(int(userID))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "no puedes publicar mas videos"})
		return
	}
	audioFileName := c.PostForm("audio")
	url := os.Getenv("MY_URL")

	// Generate a random file name
	fileName := models.GenerateRandomName()
	tempFilePath := filepath.Join(url, "/public/videos", fileName+filepath.Ext(file.Filename))

	finalVideoName := models.GenerateRandomName()
	finalVideoPath := filepath.Join(url, "/public/videos", finalVideoName+filepath.Ext(file.Filename))

	// Create the destination file
	//destPath := filepath.Join("", fileName)
	// baseDir, err := os.Getwd() // Get the current working directory

	// if url != nil {

	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	// Get the file size from the request
	fileSize := c.Request.ContentLength
	if fileSize <= 0 || fileSize > (300*1024*1024) {
		// http.Error(w, "Error: Invalid file size", http.StatusBadRequest)
		errorResponse := ErrorResponse{Error: fmt.Sprintf("File size exceeds the maximum limit of %d MB", (300*1024*1024)/(1024*1024))}
		c.JSON(http.StatusRequestEntityTooLarge, errorResponse)

		return
	}

	// user, _ := models.GetUserByIDWithVideos(userID)
	// countuservideos := len(user.Videos)
	// if user.Id_Membership == 1 && countuservideos >= 1 {
	// 	// http.Error(w, "Error: Invalid file size", http.StatusBadRequest)
	// 	errorResponse := ErrorResponse{Error: fmt.Sprintf("exceeds the maximum limit of Videos")}
	// 	c.JSON(http.StatusPreconditionFailed, errorResponse)
	// 	return
	// }

	saveVideo := new(async.Future[error])
	// destPath := filepath.Join(url, "/public/videos", fileName+filepath.Ext(file.Filename))

	// save video.

	go func() {
		// err = saveVideoFile(file, destPath)
		fmt.Println("/////////////// inicio ///////////////////")
		// time.Sleep(50 * time.Second)

		async.ResolveFuture(saveVideo, saveVideoFile(file, tempFilePath), nil)

	}()

	async.Await(saveVideo)

	fmt.Println("/////////////// final ///////////////////")

	d, err := os.Stat(tempFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(d)
	//// end save video

	var input VideoInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	v := models.Video{}
	v.Video_url = filepath.Join("/public/videos", finalVideoName+filepath.Ext(file.Filename))
	v.Image_cover = "public/video_cover/" + finalVideoName + ".jpg"
	v.Description = input.Description
	v.Location = input.Location
	v.Area = input.Area
	v.Property_number = input.Property_number
	v.Price = input.Price
	v.Id_user = userID
	v.Latitude = input.Latitude
	v.Longitude = input.Longitude
	v.Type = 4

	sale_type_id, err := strconv.ParseUint(input.Sale_type_id, 10, 32)
	if err != nil {
		// Handle the error if the conversion fails
		fmt.Println("Error converting string to uint:", err)
		return
	}
	v.Sale_type_id = int(sale_type_id)
	sale_category_id, err := strconv.ParseUint(input.Sale_category_id, 10, 32)
	if err != nil {
		// Handle the error if the conversion fails
		fmt.Println("Error converting string to uint:", err)
		return
	}
	v.Sale_category_id = int(sale_category_id)
	_, err = v.SaveVideo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// c.JSON(http.StatusOK, gin.H{"message": "Video uploaded successfully"})
	/////////////////////////////
	// Generate a unique task ID for tracking (for simplicity, using the filename, ensure uniqueness in your implementation)
	taskId := filepath.Base(finalVideoPath)
	mutex.Lock()
	taskStatusMap[taskId] = "Started"
	mutex.Unlock()

	c.JSON(200, gin.H{
		"message": "Upload received, processing started.",
		"taskId":  taskId,
	})

	go getFrame(tempFilePath, finalVideoName+".jpg")

	fmt.Println("/////////////// audio file $///////////////////", audioFileName)

	if audioFileName != "" {
		destAudioPath := filepath.Join(url, "/public/audio", audioFileName)
		// fileName = models.GenerateRandomName()

		// joinAudioVideo := new(async.Future[error])
		// // destPath := filepath.Join(url, "/public/videos", fileName+filepath.Ext(file.Filename))

		// // save video.

		// go func() {
		// 	// err = saveVideoFile(file, destPath)
		// 	fmt.Println("/////////////// inicio Join Audio video///////////////////")
		// 	// time.Sleep(50 * time.Second)

		// 	async.ResolveFuture(joinAudioVideo, joinAudioWithVideo(destAudioPath, tempFilePath, fileName+filepath.Ext(file.Filename)), nil)

		// }()

		// async.Await(joinAudioVideo)

		// fmt.Println("/////////////// final join audio Video///////////////////")

		// tempFilePath = filepath.Join(url, "/public/videos", fileName+filepath.Ext(file.Filename))
		go compressVideos(v.Id, input.Type, taskId, tempFilePath, finalVideoPath, destAudioPath, true)

		// go joinAudioWithVideo(destAudioPath, tempFilePath, fileName+filepath.Ext(file.Filename), finalVideoPath, v.Id, input.Type)
	} else {
		go compressVideos(v.Id, input.Type, taskId, tempFilePath, finalVideoPath, "", false)

	}

	//end add audio to video//////

	// Start the compression in a new goroutine
}

func compressVideos(idvideo int, typeV int, taskId, tempFilePath string, finalVideoPath string, audio_path string, has_audio bool) error {
	defer func() {
		mutex.Lock()
		taskStatusMap[taskId] = "Completed"
		mutex.Unlock()
	}()
	url := os.Getenv("api_compress_video")

	data := RequestCompressVideo{
		Video_path:       tempFilePath,   //"/home/albert/Downloads/ssstik.io_1691458134586 (copy).mp4",
		Final_video_path: finalVideoPath, //"/home/albert/Downloads/dreams.mp3",
		Audio_path:       audio_path,
		Has_audio:        has_audio,
		// Final_video_name: finalVideoName,
	}
	// // Open the file to be sent
	// file, err := os.Open(filePath)
	// if err != nil {
	//     mutex.Lock()
	//     taskStatusMap[taskId] = "Failed to open file"
	//     mutex.Unlock()
	//     return
	// }
	// defer file.Close()
	jsonData, err := json.Marshal(data)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal JSON"})
		return err
	}
	fmt.Println("HTTP JSON POST URL:", url)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal JSON"})
		return err
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	// Perform the request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		mutex.Lock()
		taskStatusMap[taskId] = "Failed to send request"
		mutex.Unlock()
		return err
	}
	defer response.Body.Close()
	deleteTemporalVideo(tempFilePath)
	fmt.Println("response Status compress:", response.Status)
	fmt.Println("response Headers compress:", response.Header)
	// Here, you can check the response status, body, etc.
	// For simplicity, this example doesn't do that.
	models.SetAvailable(idvideo, typeV)

	return nil
}

type RequestCompressVideo struct {
	Video_path       string `json:"video_path"`
	Final_video_path string `json:"final_video_path"`
	Audio_path       string `json:"audio_path"`
	Has_audio        bool   `json:"has_audio"`
}

// func compressVideo(tempFilePath string, finalVideoPath string) error {

// 	url := os.Getenv("api_compress_video")

// 	data := RequestCompressVideo{
// 		Video_path:       tempFilePath,   //"/home/albert/Downloads/ssstik.io_1691458134586 (copy).mp4",
// 		Final_video_path: finalVideoPath, //"/home/albert/Downloads/dreams.mp3",
// 		// Final_video_name: finalVideoName,
// 	}

// 	jsonData, err := json.Marshal(data)
// 	if err != nil {
// 		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal JSON"})
// 		return err
// 	}

// 	fmt.Println("HTTP JSON POST URL:", url)

// 	request, error := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
// 	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

// 	client := &http.Client{}
// 	response, error := client.Do(request)
// 	if error != nil {
// 		panic(error)
// 	}
// 	defer response.Body.Close()
// 	deleteTemporalVideo(tempFilePath)
// 	fmt.Println("response Status:", response.Status)
// 	fmt.Println("response Headers:", response.Header)
// 	// body, _ := ioutil.ReadAll(response.Body)
// 	// fmt.Println("response Body:", string(body))
// 	return nil
// }

func CheckStatus(c *gin.Context) {
	taskId := c.Param("taskId")
	mutex.Lock()
	status, exists := taskStatusMap[taskId]
	mutex.Unlock()

	if !exists {
		c.JSON(404, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(200, gin.H{"status": status})
}

func HandleVideoEdit(c *gin.Context) {

	var input VideoInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p := c.Query("idVideo")
	idVideo, err := strconv.ParseUint(p, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid idVideo"})
		return
	}

	v := models.Video{}
	v.Id = int(idVideo)
	v.Description = input.Description
	v.Location = input.Location
	v.Area = input.Area
	v.Property_number = input.Property_number
	v.Price = input.Price
	userID, _ := token.ExtractTokenID(c)
	v.Id_user = userID
	v.Latitude = input.Latitude
	v.Longitude = input.Longitude

	sale_type_id, err := strconv.ParseUint(input.Sale_type_id, 10, 32)
	if err != nil {
		// Handle the error if the conversion fails
		fmt.Println("Error converting string to uint:", err)
		return
	}
	v.Sale_type_id = int(sale_type_id)
	sale_category_id, err := strconv.ParseUint(input.Sale_category_id, 10, 32)
	if err != nil {
		// Handle the error if the conversion fails
		fmt.Println("Error converting string to uint:", err)
		return
	}
	v.Sale_category_id = int(sale_category_id)

	_, err = v.EditVideo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// err = getFrame(destPath, fileName+".jpg")
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{"message": "Video Edited successfully"})
}

func saveVideoFile(file *multipart.FileHeader, destination string) error {

	src, err := file.Open()
	if err != nil {
		// return err
		panic(err)

	}
	defer src.Close()

	dst, err := os.Create(destination)
	if err != nil {
		panic(err)

		// return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		panic(err)

		// return err
	}

	// ch <- true
	fmt.Println("/////////////// upload ///////////////////")

	return nil

}

func HandleGetCategoriesAndTypes(c *gin.Context) {

	cat, err := models.GetCategory()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	types, err := models.GetTypes()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "categories": cat, "types": types})
}

func HandleGetVideosSponsors(c *gin.Context) {

	region_code := c.Query("code")

	sponsors, err := models.GetSponsors(region_code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": sponsors})
}

func HandleGetAroundVideos(c *gin.Context) {

	lat := c.Query("latitude")
	long := c.Query("longitude")
	dist := c.Query("distance")

	// p := c.Query("page")

	// page, err := strconv.ParseUint(p, 10, 64)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page"})
	// 	return
	// }

	latitude, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
		return
	}

	longitude, err := strconv.ParseFloat(long, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
		return
	}

	distance, err := strconv.ParseFloat(dist, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
		return
	}

	userID, _ := token.ExtractTokenID(c)

	vid, err := models.GetPlacesAroundLocation(latitude, longitude, distance, int(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// types, err := models.GetTypes()
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }23

	c.JSON(http.StatusOK, gin.H{"message": "success", "videos": vid})
}

// func HandleGetAllVideos(c *gin.Context) {
// 	userID, _ := token.ExtractTokenID(c)

// 	p := c.Query("page")
// 	page, err := strconv.ParseUint(p, 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page"})
// 		return
// 	}

// 	sale := c.Query("sale")
// 	sale_id, err := strconv.ParseUint(sale, 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sale type"})
// 		return
// 	}

// 	typeV := c.Query("type")
// 	TypeVideo, err := strconv.ParseUint(typeV, 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type"})
// 		return
// 	}

// 	cat, err := models.FetchAllVideos(int(userID), int(sale_id), int(TypeVideo), int(page))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "success", "data": cat})

// }

func HandleGetAllSongs(c *gin.Context) {

	songs, err := models.GetAllSongs()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": songs})

}

// func HandleGetAllVideos(c *gin.Context) {
// 	userID, _ := token.ExtractTokenID(c)

// 	p := c.Query("page")
// 	page, err := strconv.ParseUint(p, 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page"})
// 		return
// 	}

// 	sale := c.Query("sale")
// 	sale_id, err := strconv.ParseUint(sale, 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sale type"})
// 		return
// 	}

// 	typeV := c.Query("type")
// 	typeVideo, err := strconv.ParseUint(typeV, 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type"})
// 		return
// 	}

// 	cat := c.Query("category")
// 	category := 0
// 	// if err != nil {
// 	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category"})
// 	// 	return
// 	// }
// 	switch cat {
// 	case "Residencial":
// 		category = 1
// 	case "Comercial":
// 		category = 2
// 	case "Terreno":
// 		category = 3
// 	case "Corporativo":
// 		category = 4
// 	case "Industrial":
// 		category = 5

// 	case "Residential":
// 		category = 1
// 	case "Business":
// 		category = 2
// 	case "Land":
// 		category = 3
// 	case "corporate":
// 		category = 4
// 	case "industry":
// 		category = 5
// 	case "Vip":
// 		category = 6
// 	default:
// 		category = 0
// 	}

// 	var data []models.FeedVideo

// 	if category > 0 {
// 		data, err = models.FetchAllCategoryVideos(int(userID), int(sale_id), int(typeVideo), category, int(page))
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
// 			return
// 		}
// 	} else {
// 		data, err = models.FetchAllVideos(int(userID), int(sale_id), int(typeVideo), int(page))
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
// 			return
// 		}
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "success", "data": data})

// }

func HandleGetAllVideos(c *gin.Context) {
	userID, _ := token.ExtractTokenID(c)

	// Parsear la página
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page"})
		return
	}

	// Parsear sale_id
	saleID, err := strconv.Atoi(c.Query("sale"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sale type"})
		return
	}

	// Parsear typeVideo
	typeVideo, err := strconv.Atoi(c.Query("type"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type"})
		return
	}

	// Parsear idVideo (opcional)
	var idVideo *int
	if idVideoStr := c.Query("idVideo"); idVideoStr != "" {
		id, err := strconv.Atoi(idVideoStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid idVideo"})
			return
		}
		idVideo = &id
	}

	// Determinar categoría
	categoryMap := map[string]int{
		"All":         0,
		"Todos":       0,
		"Residencial": 1,
		"Comercial":   2,
		"Terreno":     3,
		"Corporativo": 4,
		"Industrial":  5,
		"Exclusivo":   6,
		"Residential": 1,
		"Business":    2,
		"Land":        3,
		"corporate":   4,
		"industry":    5,
		"Vip":         6,
	}

	// Obtener la categoría
	categoryStr := c.Query("category")
	category, err := strconv.Atoi(categoryStr)
	if err != nil {
		// Si no es un número, intenta mapearlo
		var exists bool
		category, exists = categoryMap[categoryStr]
		if !exists && categoryStr != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category"})
			return
		}
	}

	var data []models.FeedVideo

	// Obtener los videos según la categoría
	if category > 0 {
		data, err = models.FetchAllCategoryVideos(int(userID), saleID, typeVideo, category, page, idVideo)
	} else {
		data, err = models.FetchAllVideos(int(userID), saleID, typeVideo, page, idVideo)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": data})
}

type RequestData struct {
	Path_video string `json:"path_video"`
	Image_name string `json:"image_name"`
}

func getFrame(filePath string, fileName string) error {

	url := os.Getenv("api_frame")

	// jsonStr := []byte(`{"path_video":"/home/albert/Downloads/ssstik.io_1691458134586.mp4","image_name":"kk.jpg"}`)
	// jsonStrign := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
	// Create a map to hold the data

	data := RequestData{
		Path_video: filePath,
		Image_name: fileName,
	}

	// Convert the data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal JSON"})
		return err
	}

	fmt.Println("HTTP JSON POST URL:", url)

	// var jsonData = []byte(`{
	// 	"name": "morpheus",
	// 	"job": "leader"
	// }`)
	request, error := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	fmt.Println("response Status get Frame:", response.Status)
	fmt.Println("response Headers get Frame:", response.Header)
	// body, _ := ioutil.ReadAll(response.Body)
	// fmt.Println("response Body:", string(body))
	return nil
}

type RequestAudioVideoData struct {
	Video_path       string `json:"video_path"`
	Audio_path       string `json:"audio_path"`
	Final_video_name string `json:"final_video_name"`
}

func joinAudioWithVideo(audioPath string, inputPath string, outputPath string) error {

	url := os.Getenv("api_join_audio_video")

	data := RequestAudioVideoData{
		Video_path:       inputPath, //"/home/albert/Downloads/ssstik.io_1691458134586 (copy).mp4",
		Audio_path:       audioPath, //"/home/albert/Downloads/dreams.mp3",
		Final_video_name: outputPath,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal JSON"})
		return err
	}

	fmt.Println("HTTP JSON POST URL:", url)

	request, error := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	// body, _ := ioutil.ReadAll(response.Body)
	// fmt.Println("response Body:", string(body))
	// Start the compression in a new goroutine
	// tempFilePath := filepath.Join(url, "/public/videos", outputPath)
	// taskId := filepath.Base(finalVideoPath)
	// mutex.Lock()
	// taskStatusMap[taskId] = "Started"
	// mutex.Unlock()

	// go compressVideos(idvideo, typeV, taskId, tempFilePath, finalVideoPath)
	return nil

}

func SetFavorite(c *gin.Context) {

	var input RegisterFavInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	model := models.Favorites{}

	userID, _ := token.ExtractTokenID(c)
	model.Id_user = int(userID)
	model.Id_video = input.Id_video
	// value := c.Query("id_user")
	// id_user, err := strconv.ParseUint(value, 10, 64)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id_user"})
	// 	return
	// }

	// value = c.Query("id_video")
	// sale_id, err := strconv.ParseUint(value, 10, 64)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id_video"})
	// 	return
	// }

	err := models.IsVideoFavorite(model.Id_user, model.Id_video)
	if err != nil {

		//
		fav, err := models.SetVideoFavorite(model)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error to set favorite"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "success", "data": fav})
		return
	} else {
		err := models.DeleteFavoritetByID(model.Id_user, model.Id_video)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error to delete"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "success delete"})
		return
		// c.JSON(http.StatusOK, gin.H{"message": "success"})

	}

}

func HandleSearchVideos(c *gin.Context) {
	userID, _ := token.ExtractTokenID(c)

	p := c.Query("page")
	page, err := strconv.ParseUint(p, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page"})
		return
	}

	search_text := c.Query("search")
	// search, err := strconv.p(search_text, 10, 64)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sale type"})
	// 	return
	// }

	vid, err := models.SearchVideos(search_text, int(page), int(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": vid})

}

func HandleGetTypeRepors(c *gin.Context) {
	rcode := c.Query("code")
	// regioncode, err := strconv.ParseUint(rcode, 10, 64)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page"})
	// 	return
	// }

	rep, err := models.GetReportsTypes(rcode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": rep})

}

func deleteTemporalVideo(filePath string) {
	// Verificar si el archivo existe
	if _, err := os.Stat(filePath); err == nil {
		// El archivo existe, proceder a eliminarlo
		err := os.Remove(filePath)
		if err != nil {
			// Manejar el error en caso de que el archivo no pueda ser eliminado
			fmt.Print(err)
		} else {
			// Confirmar la eliminación del archivo
			fmt.Printf("El archivo %s ha sido eliminado exitosamente.\n", filePath)
		}
	} else if os.IsNotExist(err) {
		// El archivo no existe, manejar según sea necesario
		fmt.Printf("El archivo %s no existe.\n", filePath)
	} else {
		// Hubo un problema al verificar el archivo (permisos, etc.)
		fmt.Print(err)
	}

}
