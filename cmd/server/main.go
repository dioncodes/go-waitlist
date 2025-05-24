package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dioncodes/go-waitlist/internal/db"
	"github.com/dioncodes/go-waitlist/internal/model"
	"github.com/dioncodes/go-waitlist/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("BASE_DIR") == "" {
		os.Setenv("BASE_DIR", ".")
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	fmt.Println("Starting go-waitlist server... Waiting 5 sec for DB to be up and running.")

	time.Sleep(5 * time.Second)

	db.Conn()
	db.Conn().AutoMigrate(&model.Registration{})

	if os.Getenv("ENV") != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}

	// setup routing
	r := gin.Default()
	router.Setup(r)

	r.Run(":80")
}
