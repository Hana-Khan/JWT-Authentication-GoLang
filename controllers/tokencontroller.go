package controllers
import (
	"jwt-authentication-golang/auth"
	"jwt-authentication-golang/database"
	"jwt-authentication-golang/models"
	"net/http"
	"github.com/gin-gonic/gin"
)
// Here we define a simple struct that will essentially be what the endpoint would expect as the request body. This would contain the user’s email id and password.
type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
// In the GenerateToken() function, we Bind the incoming request to the TokenRequest struct. 
// At Line 27, we communicate with the database via GORM to check if the email id passed by the request actually exists in the database. 
// If so, it will fetch the first record that matches. Else, an appropriate error message will be thrown out by the code.
//  Next, we check if the entered password matches the one in the database. For this, we will be using the CheckPassword() method that we created in the user.go file
func GenerateToken(context *gin.Context) {
	var request TokenRequest
	var user models.User
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	// check if email exists and password is correct
	record := database.Instance.Where("email = ?", request.Email).First(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}
	credentialError := user.CheckPassword(request.Password)
	if credentialError != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		context.Abort()
		return
	}

	// If everything goes well, and the password is matched, 
	//  we Generate the JWT using the GenerateJWT() function. 
	// This would return a signed token string with an expiry of 1 hour, which in turn will be sent back to the client as a response with a 200 Status Code.
	tokenString, err:= auth.GenerateJWT(user.Email, user.Username)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	context.JSON(http.StatusOK, gin.H{"token": tokenString})
}