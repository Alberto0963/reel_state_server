package controllers

import (
	// "io"
	// "mime/multipart"

	// "os"

	// "encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"reelState/auth"
	"reelState/models"
	"reelState/utils/token"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	// "golang.org/x/crypto/nacl/auth"
)

func LoginHandler(c *gin.Context) {
	// secretKey := []byte("your-secret-key") // Replace with your own secret key
	auth.Login(c, models.Pool) // Pass the DB connection and secret key to the Login function
}

func LoginWithGoogleHandler(c *gin.Context) {
	// secretKey := []byte("your-secret-key") // Replace with your own secret key
	auth.HandleGoogleLogin(c, models.Pool) // Pass the DB connection and secret key to the Login function
}

func RegisterHandler(c *gin.Context) {
	auth.Register(c)
}

func RegisterHandlerWithGogle(c *gin.Context) {
	auth.HandleGoogleRegister(c)
}

func RegisterHandlerWithFacebook(c *gin.Context) {
	auth.HandleFacebookRegister(c)
}

func LoginHandlerWithFacebook(c *gin.Context) {
	auth.HandleFacebookLogin(c)
}

func UpdatePasswordHandler(c *gin.Context) {
	auth.ResetPassword(c)
}

func UpdateUsernameHandler(c *gin.Context) {
	auth.UpdateUserName(c)
}

func UpdatePhoneNumberHandler(c *gin.Context) {
	auth.UpdatePhoneNumber(c)
}

func CurrentUserHandler(c *gin.Context) {
	auth.CurrentUser(c)
}

func UserByIdHandler(c *gin.Context) {

	actualUserID, _ := token.ExtractTokenID(c)

	id := c.Query("id")

	profileId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	result, err := models.GetUserByID(uint(profileId))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	followers, err := models.Getfollowers(int(profileId))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	imFollower, err := models.Imfollower(int(profileId), int(actualUserID))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "success",
		"data":       result,
		"followers":  followers,
		"actualUser": actualUserID,
		"imfollower": imFollower,
	})
}

func DeleteUser(c *gin.Context) {
	// Get the user ID from the URL parameter
	// userID := c.Param("id")
	userID, _ := token.ExtractTokenID(c)

	// var user models.UserUpdate

	err := models.DeleteUserByID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func SendVerificationCode(c *gin.Context) {
	auth.SendVerificationCode(c)
}

func ValidatePhone(c *gin.Context) {
	auth.ValidatePhone(c)
}

func ValidateVerificationCode(c *gin.Context) {
	auth.ValidateVerificationCode(c)
}

func SearchProfile(c *gin.Context) {

	// userID, _ := token.ExtractTokenID(c)
	username := c.Query("username")
	// username, err := strconv.ParseUint(p, 10, 64)

	// page, _ := strconv.Atoi(c.DefaultPostForm("page", "1"))

	cat, err := models.SearchProfile(username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": cat})
}

func GetMyVideos(c *gin.Context) {

	userID, _ := token.ExtractTokenID(c)
	p := c.Query("page")
	page, err := strconv.ParseUint(p, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page"})
		return
	}

	typeV := c.Query("type")
	typeVideo, err := strconv.ParseUint(typeV, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Type"})
		return
	}
	// page, _ := strconv.Atoi(c.DefaultPostForm("page", "1"))

	cat, err := models.GetMyVideos(int(userID), int(page), int(typeVideo))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": cat})
}

func GetUserVideos(c *gin.Context) {

	id := c.Query("id")

	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	typeV := c.Query("type")
	typeVideo, err := strconv.ParseUint(typeV, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page"})
		return
	}
	// userID, _ := token.ExtractTokenID(c)

	p := c.Query("page")
	page, err := strconv.ParseUint(p, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page"})
		return
	}

	cat, err := models.GetMyVideos(int(userID), int(page), int(typeVideo))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": cat})
}

func GetUserFavoritesVideos(c *gin.Context) {

	userID, _ := token.ExtractTokenID(c)

	p := c.Query("page")
	page, err := strconv.ParseUint(p, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page"})
		return
	}

	cat, err := models.GetMyFavoritesVideos(int(userID), int(page))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": cat})
}

func ValidateUserName(c *gin.Context) {

	username := c.DefaultPostForm("UserName", "")

	if models.UsernameExists(username) {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Username is available"})
}

func UpdateProfileImageUserName(c *gin.Context) {

	profileImage, err := c.FormFile("profile_image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from the request"})
		return
	}
	// Generate a random file name for the profile image
	imageFileName := models.GenerateRandomName() + filepath.Ext(profileImage.Filename)
	url := os.Getenv("MY_URL")
	profileImagePath := filepath.Join("public/profile_images", imageFileName)

	userID, _ := token.ExtractTokenID(c)

	u := models.UserUpdate{}

	u.ProfileImage = profileImagePath
	u.ID = userID

	_, err = u.UpdateProfileImageUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.SaveUploadedFile(profileImage, url+profileImagePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save profile image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "image is updated"})
}

func UpdateLinkUserName(c *gin.Context) {

	link := c.Query("link")
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from the request"})
	// 	return
	// }
	// Generate a random file name for the profile image
	// imageFileName := models.GenerateRandomName() + filepath.Ext(profileImage.Filename)
	// url := os.Getenv("MY_URL")
	// profileImagePath := filepath.Join("public/profile_images", imageFileName)

	userID, _ := token.ExtractTokenID(c)

	u, err := models.GetUserByIdToUpdate(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get user"})
		return
	}
	u.Link = link
	// u.ID = userID

	_, err = u.UpdateUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "link is updated"})
}

func UpdateCoverImageUserName(c *gin.Context) {

	coverImage, err := c.FormFile("cover_image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from the request"})
		return
	}
	// Generate a random file name for the profile image
	imageFileName := models.GenerateRandomName() + filepath.Ext(coverImage.Filename)
	url := os.Getenv("MY_URL")
	coverImagePath := filepath.Join("public/profile_images", imageFileName)

	userID, _ := token.ExtractTokenID(c)

	u := models.User{}

	u.Cover_image = coverImagePath
	u.ID = userID

	_, err = u.UpdateCoverImageUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.SaveUploadedFile(coverImage, url+coverImagePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save profile image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "image is updated"})
}

func GetPublicMemberShips(c *gin.Context) {

	// userID, _ := token.ExtractTokenID(c)
	p := c.Query("page")
	page, err := strconv.ParseUint(p, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page"})
		return
	}

	currencyCode := c.Query("currencyCode")
	// currencyCode, err := strconv.ParseUint(p, 10, 64)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page"})
	// 	return
	// }

	cat, err := models.GetMemberShips(currencyCode, int(page))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": cat})
}

func GetMemberShips(c *gin.Context) {

	// userID, _ := token.ExtractTokenID(c)
	user_id, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p := c.Query("page")
	page, err := strconv.ParseUint(p, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page"})
		return
	}

	currencyCode := c.Query("currencyCode")
	// currencyCode, err := strconv.ParseUint(p, 10, 64)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page"})
	// 	return
	// }

	cat, err := models.GetMemberShips(currencyCode, int(page))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	u, err := models.GetUserByID(user_id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": cat, "user": u})
}

func DeleteUserVideo(c *gin.Context) {

	// userID, _ := token.ExtractTokenID(c)
	user_id, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p := c.Query("idVideo")
	idVideo, err := strconv.ParseUint(p, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	err = models.DeleteUserVideo(int(idVideo), int(user_id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func Setlike(c *gin.Context) {

	actualUserID, _ := token.ExtractTokenID(c)
	// var u *models.Likes
	id := c.Query("id")
	profileId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	u, err := models.GetProfilefollowers(int(profileId), int(actualUserID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting user"})
		return
	}

	if u.Id == 0 {
		u = models.Likes{}
		u.Id_profile = uint(profileId)
		u.Id_user = actualUserID
		models.SaveLikeProfile(&u)

	} else {
		models.UpdateLikeProfile(&u)
	}

	// if(u == nil){
	// 	u = *models.Likes{}

	// 	u.Id_profile = uint(profileId)
	// 	u.Id_user = actualUserID
	// 	models.SaveLikeProfile(u)
	// }

	// _, err = u.UpdateUser()
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// token, err := token.GenerateToken(u.ID)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{"message": "update success", "d": u})
	// c.JSON(http.StatusOK, gin.H{"message": "validated!"})

}

func CreateSubscription(c *gin.Context) {

	var sub models.Createsubscription
	actualUserID, _ := token.ExtractTokenID(c)

	// Bind JSON to struct
	if err := c.ShouldBind(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub.IdUser = int(actualUserID)
	layout := "2006-01-02 15:04:05.999999999 -0700 MST"
	now := time.Now()
	formattedDate := now.Format("2006-01-02 15:04:05.999999999 -0700 MST")

	// Parse the date string
	parsedTime, err := time.Parse(layout, formattedDate)
	if err != nil {
		// fmt.Println("Error parsing date:", err)
		c.JSON(http.StatusBadRequest, gin.H{"Error Parse Time": err.Error()})

		return
	}

	models.CancelSubscriptionIfActive(strconv.FormatUint(uint64(actualUserID), 10), "user Create new membership", parsedTime, "paypal")

	_, err = sub.CreateSubscription()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	// Respond with success
	c.JSON(http.StatusCreated, gin.H{"message": "Subscription created successfully"})
}

func CancelSubscription(c *gin.Context) {
	var sub models.CancelSubscription

	// Bind JSON to struct
	if err := c.ShouldBind(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	layout := "2006-01-02 15:04:05.999999999 -0700 MST"
	now := time.Now()
	formattedDate := now.Format("2006-01-02 15:04:05.999999999 -0700 MST")

	// Parse the date string
	parsedTime, err := time.Parse(layout, formattedDate)
	if err != nil {
		// fmt.Println("Error parsing date:", err)
		c.JSON(http.StatusBadRequest, gin.H{"Error Parse Time": err.Error()})

		return
	}

	subid, err := strconv.ParseUint(sub.Id, 10, 32)
	if err != nil {
		// fmt.Println("Error al convertir string a uint:", err)
		c.JSON(http.StatusBadRequest, gin.H{"Error al convertir string a uint:": err.Error()})

		return
	}

	subscriber, err := models.GetSubscriptionByID(uint(subid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error Paypal Sub": err.Error()})
		return
	}

	if subscriber.CustomerId == "" {
		err = models.CancelPaypalSubscription(subscriber.PaypalSubscriptionId, sub.Reason)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error Paypal Sub": err.Error()})
			return
		}
	} else {
		err = models.CancelOpenPaySubscription(subscriber.CustomerId, subscriber.PaypalSubscriptionId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error Paypal Sub": err.Error()})
			return
		}
	}

	err = models.CancelSubscriptionFunction(subscriber.ID, parsedTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error Calcel sub": err.Error()})
		return
	}

	// Respond with success
	c.JSON(http.StatusCreated, gin.H{"message": "Subscription cancelled successfully"})
}

func GetUserSubscription(c *gin.Context) {

	// var sub models.Createsubscription
	actualUserID, _ := token.ExtractTokenID(c)

	// Bind JSON to struct
	// if err := c.ShouldBind(&sub); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// sub.IdUser = int(actualUserID)

	subs, err := models.GetSubscription(int(actualUserID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	// Respond with success
	c.JSON(http.StatusCreated, gin.H{"data": subs})
}

func MakePayOpenPay(c *gin.Context) {
	var card models.ChargeRequest
	actualUserID, _ := token.ExtractTokenID(c)

	// Decodificar los datos de la tarjeta
	if err := c.ShouldBindJSON(&card); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Datos de la tarjeta inválidos", "error": err.Error()})
		return
	}

	user, _ := models.GetUserByID(actualUserID)
	// Obtener el token de la tarjeta
	token := card.SourceID
	amount := card.Amount // Ejemplo de monto
	// description := card.Description       // Ejemplo de descripción
	// deviceSessionID := card.DeviceSession // Se debe obtener el ID de sesión del dispositivo
	planid := card.PlanID

	nombre, apellidos := models.SepararNombreCompleto(card.Name)

	customer := models.Customer{Name: nombre, LastName: apellidos, Email: card.Email, PhoneNumber: user.Phone}
	customerID, err := models.CreateCustomer(customer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error al procesar el pago", "error": err.Error()})
		return
	}
	// Realizar el pago
	var open_sub models.SubscriptionResponse

	if open_sub, err = models.CreateSubscriptionWithToken(customerID, planid, token, amount, 0); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error al procesar el pago", "error": err.Error()})
		return
	}

	var sub models.Createsubscription

	sub.IdUser = int(actualUserID)
	layout := "2006-01-02 15:04:05.999999999 -0700 MST"
	now := time.Now()
	formattedDate := now.Format("2006-01-02 15:04:05.999999999 -0700 MST")

	sub.PaypalSubscriptionId = open_sub.ID
	membershipIDInt64, err := strconv.ParseInt(card.MembershipID, 10, 64)
	sub.Renewal = true
	sub.CustomerId = customerID

	// if err != nil {
	// 	fmt.Println("Error al convertir MembershipID a int64:", err)
	// } else {
	// 	fmt.Println("MembershipID como entero int64:", membershipIDInt64)
	// }
	sub.MembershipId = int(membershipIDInt64)
	// Parse the date string
	parsedTime, err := time.Parse(layout, formattedDate)
	if err != nil {
		// fmt.Println("Error parsing date:", err)
		c.JSON(http.StatusBadRequest, gin.H{"Error Parse Time": err.Error()})

		return
	}

	models.CancelSubscriptionIfActive(strconv.FormatUint(uint64(actualUserID), 10), "user Create new membership", parsedTime, "openpay")

	_, err = sub.CreateSubscription()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusCreated, gin.H{"status": "Pago procesado con éxito", "opensub": open_sub, "reelstateSub": sub})
}

// CreateOrUpdateReview maneja la creación o actualización de reseñas
func CreateOrUpdateReview(c *gin.Context) {
	// Inicializamos una variable del tipo Review
	var review models.Review

	// Vinculamos la solicitud JSON al struct Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validación: id_user no debe ser igual a id_profile
	if review.IDUser == review.IDProfile {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User cannot review their own profile"})
		return
	}

	// Intentamos crear o actualizar la reseña
	if err := review.CreateOrUpdate(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create or update review"})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{"message": "Review created or updated successfully"})
}


// Struct para recibir el JSON
type ProfileRequest struct {
    IdProfile int `json:"id_profile"` // O usa el tipo adecuado según tu modelo
}
// GetReviewsByProfile maneja la obtención de reseñas por id_profile
func GetReviewsByProfile(c *gin.Context) {
	var req ProfileRequest

    // Bind JSON a la estructura
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Formato JSON inválido"})
        return
    }

	var reviews []models.Review
	// var review models.Review

	reviews, err := models.GetReviewsProfile(req.IdProfile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reviews})

	// c.JSON(http.StatusOK, reviews)
}
