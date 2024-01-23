package controllers

import (
	"go-jwt/initializers"
	"go-jwt/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(context *gin.Context) {
	// Get email/password of the request body
	var body struct {
		Email    string
		Password string
	}

	if context.Bind(&body) != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read the body!!"})
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to hash password!!"})
		return
	}

	// Create the user
	newUser := models.User{Email: body.Email, Password: string(hash)}
	result := initializers.DB.Create(&newUser)

	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user!!"})
		return
	}

	// Respond
	context.JSON(http.StatusOK, gin.H{"message": "User created successfully!!"})
}

func Login(context *gin.Context) {
	// Get email/password of the request body
	var body struct {
		Email    string
		Password string
	}

	if context.Bind(&body) != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read the body!!"})
		return
	}

	// Look up requested user
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password!!"})
		return
	}

	// Compare sent in password with saved user hash password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password!!"})
		return
	}

	// Generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to generate token!!"})
		return
	}

	// Setup the cookie
	context.SetSameSite(http.SameSiteLaxMode)
	context.SetCookie("Authorization", string(tokenString), 3600*24*30, "", "", false, true)

	context.JSON(http.StatusOK, gin.H{"message": "Logged in successfully!!"})
}

func Validate(context *gin.Context) {
	user, _ := context.Get("user")
	context.JSON(http.StatusOK, gin.H{"message": user})
}
