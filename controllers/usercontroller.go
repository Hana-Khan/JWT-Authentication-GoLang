package controllers
import (
	"jwt-authentication-golang/auth"
	"jwt-authentication-golang/database"
	"jwt-authentication-golang/models"
	"net/http"
	"github.com/gin-gonic/gin"
)
func RegisterUser(context *gin.Context) {
	//Here we declare a local variable of type models.User.
	var user models.User
	// Whatever is sent by the client as a JSON body will be mapped into the user variable.
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	// Here, we hash the password using the bcrypt helpers that we added earlier to the models/user.go file.
	if err := user.HashPassword(user.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	// Once hashed, we store the user data into the database using the GORM global instance that we initialized earlier in the main file.
	record := database.Instance.Create(&user)
	// If there is an error while saving the data, the application would throw an HTTP Internal Server Error Code 500 and abort the request.
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}
	// Finally, if everything goes well, we send back the user id, name, and email to the client along with a 200 SUCCESS status code.
	context.JSON(http.StatusCreated, gin.H{"userId": user.ID, "email": user.Email, "username": user.Username})
}
func Login(context *gin.Context) {
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
	context.JSON(http.StatusOK, gin.H{"user":gin.H{"userId": user.ID, "email": user.Email, "username": user.Username},"token": tokenString})
}