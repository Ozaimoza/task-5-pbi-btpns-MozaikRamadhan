package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"task-5-pbi-btpns/app"
	"task-5-pbi-btpns/database"
	"task-5-pbi-btpns/helpers"
	"task-5-pbi-btpns/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

// Register
func RegisterUser(c *gin.Context) {
	var body app.User
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Initialize validator
	v := validator.New()

	err := v.Struct(body)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	//check username
	var findUser models.UserModel
	database.DB.First(&findUser, "username = ?", body.Username)

	if findUser.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is already taken"})
		return
	}

	//check Email
	database.DB.First(&findUser, "email = ?", body.Email)

	if findUser.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email Is Already In Use"})
		return
	}

	// Hash password
	hashedPassword, err := helpers.HashPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	var newUser models.UserModel
	newUser.Password = hashedPassword
	newUser.Email = body.Email
	newUser.Username = body.Username

	// Save New User To Database
	database.DB.Create(&newUser)

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "data": newUser})
}

// Login
func LoginUser(c *gin.Context) {

	var body app.User
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.UserModel
	database.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Email or Password"})
		return
	}

	//compare password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Password"})
		return
	}

	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"sub": user.ID,
	// 	"exp": time.Now().Add(time.Hour * 24 * 2).Unix(),
	// })

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := helpers.GenerateJWTToken(user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed To Create Token"})
		return
	}

	helpers.SetTokenCookie(c, tokenString, user.ID)

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

// Logout
func LogoutUser(c *gin.Context) {

	// Setting MaxAge to a negative value will make the browser delete the cookie.
	c.SetCookie("token", "", -1, "", "", false, true)
	c.SetCookie("currentUser", "", -1, "", "", false, true)

	c.String(http.StatusOK, "Success Logout")
}

// Update User
func UpdateUser(c *gin.Context) {
	// Get Body
	var body app.User
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Initialize validator
	v := validator.New()

	err := v.Struct(body)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// Get User
	var user models.UserModel
	userId := c.Param("userId")
	database.DB.First(&user, userId)

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User Not Found !!"})
		return
	}

	var findUser models.UserModel
	//check username
	if body.Username != user.Username {
		database.DB.First(&findUser, "username = ?", body.Username)

		if findUser.ID != 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username is already taken"})
			return
		}

		user.Username = body.Username
	}

	//check Email
	if body.Email != user.Email {
		database.DB.First(&findUser, "email = ?", body.Email)

		if findUser.ID != 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email Is Already In Use"})
			return
		}

		user.Email = body.Email
	}

	// Hash password

	er := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if er != nil {
		hashedPassword, er := helpers.HashPassword(body.Password)
		if er != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user.Password = hashedPassword
	}

	//save Updated Data
	database.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "Update Success", "data": user})
}

func DeleteUser(c *gin.Context) {
	userIDInterface, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	// convert type interface{} to int
	userID := int(userIDInterface.(float64))
	param := c.Param("userId")

	userId, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid userId"})
		return
	}

	//get user data
	var user models.UserModel
	database.DB.First(&user, userId)

	// check if not exists
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	// Compare userID (from middleware) with userId (from converted URL parameters)
	if userID != userId {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	// If userID is equal to userId, delete the user from the database (soft delete)
	database.DB.Delete(&user)
	c.JSON(http.StatusOK, gin.H{"message": "Success Delete Data User", "data": user})
}
