package controllers

import (
	// "io"
	// "mime/multipart"
	// "go/token"
	// "fmt"
	"bytes"
	"encoding/json"
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
	// "golang.org/x/crypto/nacl/auth"
)

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
	Sale_type_id     string `json:"sale_type_id" binding:"required"`
	Sale_category_id string `json:"sale_category_id" binding:"required"`
}

func HandleVideoUpload(c *gin.Context) {

	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a random file name
	fileName := models.GenerateRandomName() 

	// Create the destination file
	//destPath := filepath.Join("", fileName)
	// baseDir, err := os.Getwd() // Get the current working directory
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	url := os.Getenv("MY_URL")
	// if url != nil {

	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	destPath := filepath.Join(url, "/public/videos", fileName + filepath.Ext(file.Filename))
	err = saveVideoFile(file, destPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// coverUrl,err := SMS.GenerateImageFromVideo(destPath)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var input VideoInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	v := models.Video{}
	v.Video_url = ("public/videos/" + fileName)
	v.Description = input.Description
	v.Location = input.Location
	v.Area = input.Area
	v.Property_number = input.Property_number
	v.Price = input.Price
	userID, _ := token.ExtractTokenID(c)
	v.Id_user = userID
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

	err = getFrame(destPath, fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Video uploaded successfully"})
}

func saveVideoFile(file *multipart.FileHeader, destination string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

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

func HandleGetAllVideos(c *gin.Context) {

	cat, err := models.FetchAllVideos()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": cat})

}

func getFrame(filePath string, fileName string) error {

	url := os.Getenv("api_frame")

    jsonStr := []byte(`{"path_video":"/home/albert/Downloads/ssstik.io_1691458134586.mp4","image_name":"kk.jpg"}`)
	jjson := []byte(`{"path_video":"/home/albert/Downloads/ssstik.io_1691458134586.mp4","image_name":"kk.jpg"}`)

	d ,_:= json.Marshal(jjson)
	fmt.Println("Response:", d)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
		return err
    }
    defer resp.Body.Close()
	var result bytes.Buffer
	_, err = io.Copy(&result, resp.Body)
	if err != nil {
		// fmt.Println("Error reading response:", err)
		return err
	}

	fmt.Println("Response:", result.String())
	return nil
	// return nil
	
	// url := os.Getenv("api_frame")

	// // Create a new multipart writer to write the file as part of the request body
	// body := &bytes.Buffer{}
	// writer := multipart.NewWriter(body)

	
    // jsonStr := []byte(`{"path_video":"/home/albert/Downloads/ssstik.io_1691458134586.mp4","image_name":fileName + ".jpg"}`)

	

	// request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	// if err != nil {
	// 	// fmt.Println("Error creating request:", err)
	// 	return err
	// }

	// // Set the content type for the request
	// // request.Header.Set("Content-Type", writer.FormDataContentType())
	// request.Header.Set("Content-Type", "application/json")
	// // Make the request and get the response
	// client := http.Client{}
	// response, err := client.Do(request)
	// if err != nil {
	// 	// fmt.Println("Error making request:", err)
	// 	return err
	// }
	// defer response.Body.Close()

	// // Check the response status code
	// if response.StatusCode != http.StatusOK {
	// 	// fmt.Println("Request failed with status code:", response.StatusCode)
	// 	return err
	// }

	// // Read the response body
	// var result bytes.Buffer
	// _, err = io.Copy(&result, response.Body)
	// if err != nil {
	// 	// fmt.Println("Error reading response:", err)
	// 	return err
	// }

	// fmt.Println("Response:", result.String())
	// return nil
}
