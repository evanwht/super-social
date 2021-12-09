package main

import (
	"fmt"
	"log"
	"os"

	"superhuman-social/internal/clearbit"
	"superhuman-social/internal/database"
	"superhuman-social/internal/endpoints"

	"github.com/gin-gonic/gin"
)

var (
	clearbitAPIKey string
	projectID      string
)

func initialize() error {
	clearbitAPIKey = os.Getenv("CLEARBIT_API_KEY")
	if len(clearbitAPIKey) == 0 {
		return fmt.Errorf("clear bit api key not found")
	}

	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	if len(projectID) == 0 {
		return fmt.Errorf("project not found")
	}
	return nil
}

func main() {
	if err := initialize(); err != nil {
		log.Fatalf("error intitializing server: %v", err)
	}

	clearbitAPI := clearbit.NewClearbitAPI(clearbitAPIKey)
	db, err := database.NewDB(projectID)
	if err != nil {
		log.Fatalf("error initializing database client: %v", err)
	}

	r := gin.Default()
	r.GET("/lookup", endpoints.LookupPerson(clearbitAPI, db))
	r.GET("/popularity", endpoints.ListPopular(db))

	// TODO graceful shutdown
	log.Fatal(r.Run(":8080"))
}
