package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/jwt-go/initializers"
	"github.com/jwt-go/models"
	"golang.org/x/crypto/bcrypt"
)

type userInfo struct {
	Email    string
	Password string
}

func Signup(c *gin.Context) {
	// Get the email/pass off req Body

	var ui userInfo

	if c.Bind(&ui) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	// Here we are hashing the password, as we would be storing the hash password in the db
	hash, err := bcrypt.GenerateFromPassword([]byte(ui.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password.",
		})
		return
	}

	// Create the user
	user := models.User{Email: ui.Email, Password: string(hash)}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user.",
		})
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{})
}

func Login(c *gin.Context) {
	// Get email & pass off req body
	var ui userInfo

	if c.Bind(&ui) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	// Look up for requested user
	var user models.User

	initializers.DB.First(&user, "email = ?", ui.Email)

	//the default value for a non-existent primary key is zero
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Compare sent in password with saved users password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(ui.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	// Respond
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{})
}

// func Validate(c *gin.Context){
// 	user,_ := c.Get("user")

// 	// user.(models.User).Email    -->   to access specific data

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": user,
// 	})
// }