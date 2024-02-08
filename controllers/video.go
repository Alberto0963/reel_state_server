package controllers

import (
	// "io"
	// "mime/multipart"
	// "go/token"
	// "fmt"
	"bytes"
	"encoding/json"

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
	"github.com/jrallison/go-workers"
	// "github.com/jrallison/go-workers"
	// "golang.org/x/crypto/nacl/auth"
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
	Type             float64 `json:"type"`
}

func checkVideoSize(filePath string, maxSize int64) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	fileSize := fileInfo.Size()
	if fileSize > maxSize {
		return fmt.Errorf("Error: File size exceeds the maximum limit of %d MB", maxSize/(1024*1024))
	}

	return nil
}

// ErrorResponse represents a JSON error response structure
type ErrorResponse struct {
	Error string `json:"error"`
}

func HandleVideoUpload(c *gin.Context) {

	// c.Request.Body = http.MaxBytesReader(c.Request.Response., c.Request.Body, 300*1024*1024)

	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	userID, _ := token.ExtractTokenID(c)
	user, _ := models.GetUserByIDWithVideos(userID)
	countuservideos := len(user.Videos)
	if user.Id_Membership == 1 && countuservideos >= 1 {
		// http.Error(w, "Error: Invalid file size", http.StatusBadRequest)
		errorResponse := ErrorResponse{Error: fmt.Sprintf("exceeds the maximum limit of Videos")}
		c.JSON(http.StatusPreconditionFailed, errorResponse)
		return
	}

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

	//add audio to video
	if audioFileName != "" {
		destAudioPath := filepath.Join(url, "/public/audio", audioFileName)
		fileName = models.GenerateRandomName()

		workers.Enqueue("myqueue", "Add", joinAudioWithVideo(destAudioPath, tempFilePath, fileName+filepath.Ext(file.Filename)))
		tempFilePath = filepath.Join(url, "/public/videos", fileName+filepath.Ext(file.Filename))

	}

	workers.Enqueue("myqueue", "Add", getFrame(tempFilePath, fileName+".jpg"))

	// Add a job to a queue
	workers.Enqueue("myqueue", "Add", compressVideo(tempFilePath, finalVideoPath))

	var input VideoInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	v := models.Video{}
	v.Video_url = filepath.Join("/public/videos", finalVideoName+filepath.Ext(file.Filename))
	v.Image_cover = "public/video_cover/" + fileName + ".jpg"
	v.Description = input.Description
	v.Location = input.Location
	v.Area = input.Area
	v.Property_number = input.Property_number
	v.Price = input.Price
	v.Id_user = userID
	v.Latitude = input.Latitude
	v.Longitude = input.Longitude
	v.Type = input.Type

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

	c.JSON(http.StatusOK, gin.H{"message": "Video uploaded successfully"})
}

// func HandleVideoUpload(c *gin.Context) {

// 	// c.Request.Body = http.MaxBytesReader(c.Request.Response., c.Request.Body, 300*1024*1024)

// 	file, err := c.FormFile("video")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

// 	userID, _ := token.ExtractTokenID(c)
// 	user, _ := models.GetUserByIDWithVideos(userID)
// 	countuservideos := len(user.Videos)
// 	if user.Id_Membership == 1 && countuservideos >= 1 {
// 		// http.Error(w, "Error: Invalid file size", http.StatusBadRequest)
// 		errorResponse := ErrorResponse{Error: fmt.Sprintf("exceeds the maximum limit of Videos")}
// 		c.JSON(http.StatusPreconditionFailed, errorResponse)
// 		return
// 	}

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

// 	err = getFrame(tempFilePath, fileName+".jpg")
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	//add audio to video
// 	if audioFileName != "" {
// 		destAudioPath := filepath.Join(url, "/public/audio", audioFileName)
// 		fileName = models.GenerateRandomName()
// 		saveVideoWithAudio := new(async.Future[error])

// 		go func() {
// 			// err = saveVideoFile(file, destPath)
// 			fmt.Println("/////////////// inicio ///////////////////")
// 			// time.Sleep(50 * time.Second)

// 			async.ResolveFuture(saveVideoWithAudio, joinAudioWithVideo(destAudioPath, tempFilePath, fileName+filepath.Ext(file.Filename)), nil)

// 		}()

// 		async.Await(saveVideoWithAudio)

// 		fmt.Println("/////////////// final ///////////////////")
// 		tempFilePath = filepath.Join(url, "/public/videos", fileName+filepath.Ext(file.Filename))

// 		d, err = os.Stat(tempFilePath)

// 		fmt.Println(d)

// 	}

// 	/// end add audio to video

// 	// compress video

// 	// // Register the worker function for a specific queue
// 	// workers.Process("myqueue", compressVideo(tempFilePath, finalVideoPath), 10)
// 	// Add a job to a queue
// 	// workers.Enqueue("myqueue", "Add", compressVideo(tempFilePath, finalVideoPath))
// 	// Enqueue a job onto the specified queue
// 	// workers.Enqueue("myqueue", "MyBackgroundTask", nil)
// 	go compressVideo(tempFilePath, finalVideoPath)

// 	// err = compressVideo(tempFilePath, finalVideoPath)
// 	// if err != nil {
// 	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 	// 	return
// 	// }

// 	// end compress video

// 	// var finalpath = ("public/videos/" + fileName + filepath.Ext(file.Filename))
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

// 	sale_type_id, err := strconv.ParseUint(input.Sale_type_id, 10, 32)
// 	if err != nil {
// 		// Handle the error if the conversion fails
// 		fmt.Println("Error converting string to uint:", err)
// 		return
// 	}
// 	v.Sale_type_id = int(sale_type_id)
// 	sale_category_id, err := strconv.ParseUint(input.Sale_category_id, 10, 32)
// 	if err != nil {
//
// Handle the error if the conversion fails
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

type RequestCompressVideo struct {
	Video_path       string `json:"video_path"`
	Final_video_path string `json:"final_video_path"`
}

func compressVideo(tempFilePath string, finalVideoPath string) error {

	url := os.Getenv("api_compress_video")

	data := RequestCompressVideo{
		Video_path:       tempFilePath,   //"/home/albert/Downloads/ssstik.io_1691458134586 (copy).mp4",
		Final_video_path: finalVideoPath, //"/home/albert/Downloads/dreams.mp3",
		// Final_video_name: finalVideoName,
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
	return nil
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

func HandleGetAllVideos(c *gin.Context) {
	userID, _ := token.ExtractTokenID(c)

	p := c.Query("page")
	page, err := strconv.ParseUint(p, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page"})
		return
	}

	sale := c.Query("sale")
	sale_id, err := strconv.ParseUint(sale, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sale type"})
		return
	}

	typeV := c.Query("type")
	TypeVideo, err := strconv.ParseUint(typeV, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type"})
		return
	}

	cat, err := models.FetchAllVideos(int(userID), int(sale_id), int(TypeVideo), int(page))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": cat})

}

func HandleGetAllSongs(c *gin.Context) {

	songs, err := models.GetAllSongs()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": songs})

}

func HandleGetAllCategoriesVideos(c *gin.Context) {
	userID, _ := token.ExtractTokenID(c)

	p := c.Query("page")
	page, err := strconv.ParseUint(p, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page"})
		return
	}

	sale := c.Query("sale")
	sale_id, err := strconv.ParseUint(sale, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sale type"})
		return
	}

	typeV := c.Query("type")
	typeVideo, err := strconv.ParseUint(typeV, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type"})
		return
	}

	cat := c.Query("category")
	category := 0
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category"})
	// 	return
	// }
	switch cat {
	case "Residencial":
		category = 1
	case "Comercial":
		category = 2
	case "Terreno":
		category = 3
	case "Corporativo":
		category = 4
	case "Industrial":
		category = 5

	case "Residential":
		category = 1
	case "Business":
		category = 2
	case "Land":
		category = 3
	case "corporate":
		category = 4
	case "industry":
		category = 5
	case "Vip":
		category = 6
	default:
		category = 0
	}

	data, err := models.FetchAllCategoryVideos(int(userID), int(sale_id), int(typeVideo), category, int(page))
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

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	// body, _ := ioutil.ReadAll(response.Body)
	// fmt.Println("response Body:", string(body))
	return nil
}

type RequestAudioVideoData struct {
	Video_path       string `json:"video_path"`
	Audio_path       string `json:"audio_path"`
	Final_video_name string `json:"final_video_name"`
}

func joinAudioWithVideo(audioPath string, videoPath string, finalVideoName string) error {

	url := os.Getenv("api_join_audio_video")

	data := RequestAudioVideoData{
		Video_path:       videoPath, //"/home/albert/Downloads/ssstik.io_1691458134586 (copy).mp4",
		Audio_path:       audioPath, //"/home/albert/Downloads/dreams.mp3",
		Final_video_name: finalVideoName,
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
