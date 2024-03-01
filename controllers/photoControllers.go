package controllers

import (
	"net/http"
	"os"
	"task-5-pbi-btpns/app"
	"task-5-pbi-btpns/database"
	"task-5-pbi-btpns/helpers"
	"task-5-pbi-btpns/models"

	"github.com/gin-gonic/gin"
)

func GetAllPhoto(c *gin.Context) {
	var photo []models.PhotoModel

	database.DB.Find(&photo)

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"data":    photo,
	})
}

func PostPhoto(c *gin.Context) {
	var body app.Photo

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	// check file upload
	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, "No file uploaded")
		return
	}
	defer file.Close()

	// check file type
	if header.Header.Get("Content-Type") != "image/jpeg" && header.Header.Get("Content-Type") != "image/png" {
		c.JSON(http.StatusBadRequest, "Invalid file type. Only JPEG and PNG are allowed")
		return
	}

	// check size
	if header.Size > (5 << 20) { // maksimum 5MB
		c.JSON(http.StatusBadRequest, "File size too large")
		return
	}

	// Generate new name
	newName := helpers.GenerateNewFileName(header.Filename)

	// save file with new name to local
	err = c.SaveUploadedFile(header, "assets/"+newName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Unknown Error")
		return
	}

	photoURL := "assets/" + newName

	// get user id from cookie
	userIDInterface, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	// convert to uint
	userID := uint(userIDInterface.(float64))

	// create photo model
	var photo models.PhotoModel
	photo.Title = body.Title
	photo.Caption = body.Caption
	photo.PhotoUrl = photoURL
	photo.UserID = userID

	// Save New photo to database
	database.DB.Create(&photo)

	c.JSON(http.StatusOK, gin.H{"message": "Success", "data": photo})
}

func UpdatePhoto(c *gin.Context) {
	userIDInterface, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	// convert type interface{} to int
	userID := uint(userIDInterface.(float64))
	//get params
	photoID := c.Param("photoId")

	// get photo
	var photo models.PhotoModel
	if err := database.DB.Find(&photo, photoID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Photo not found"})
		return
	}

	// get photo from request
	var newPhoto app.Photo
	if err := c.ShouldBind(&newPhoto); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	// Check if there is a new photo uploaded
	if newPhoto.Photo != nil {
		// remove old photo from system file
		if err := os.Remove(photo.PhotoUrl); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to remove old photo"})
			return
		}

		// Generate new filename
		newName := helpers.GenerateNewFileName(newPhoto.Photo.Filename)

		// save file with new filename
		if err := c.SaveUploadedFile(newPhoto.Photo, "assets/"+newName); err != nil {
			c.JSON(http.StatusInternalServerError, "Unknown Error")
			return
		}

		// Update photo URL in the database
		photo.PhotoUrl = "assets/" + newName
	}

	// Update other photo data
	photo.Title = newPhoto.Title
	photo.Caption = newPhoto.Caption

	// validate User
	if photo.UserID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	//save new data
	if err := database.DB.Save(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update photo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Photo updated successfully", "data": photo})
}

func DeletePhoto(c *gin.Context) {
	userIDInterface, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	// convert type interface{} to int
	userID := uint(userIDInterface.(float64))
	//get params
	photoID := c.Param("photoId")

	// get photo data
	var photo models.PhotoModel
	if err := database.DB.Find(&photo, photoID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Photo not found"})
		return
	}

	// Compare userID (from middleware) with userId (from photo data)
	if userID != photo.UserID {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	// remove photo from system file
	if err := os.Remove(photo.PhotoUrl); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to remove old photo"})
		return
	}

	// If userID is equal to userId, delete the user from the database (soft delete)
	database.DB.Delete(&photo)

	c.JSON(http.StatusOK, gin.H{"message": "Success Delete Photo", "data": photo})
}
