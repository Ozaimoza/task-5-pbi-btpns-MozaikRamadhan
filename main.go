package main

import (
	"fmt"
	"log"
	"os"
	"task-5-pbi-btpns/database"
	"task-5-pbi-btpns/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// db Connection
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// migrate db
	err = database.AutoMigrate(db)
	if err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	fmt.Println("Database migration successful")

	// Init
	r := gin.Default()

	// Setup router untuk pengguna (users)
	routes.UserRouter(r)
	routes.PhotoRouter(r)

	r.Run(os.Getenv("PORT"))
}
