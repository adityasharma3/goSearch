package main

import (
	"log"
	"os"

	"github.com/adityasharma3/goSearch/cmd/search/routes"
	elasticsearch "github.com/adityasharma3/goSearch/cmd/search/searchclient"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	elasticsearch.InitializeElasticSearch()
	app := routes.SetupRoutes()
	enviornment := os.Getenv("ENVIORNMENT")
	app.Run(enviornment)
}
