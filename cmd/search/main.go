package main

import (
	"os"

	"github.com/adityasharma3/goSearch/cmd/search/routes"
	elasticsearch "github.com/adityasharma3/goSearch/cmd/search/searchclient"
)

func main() {
	elasticsearch.InitializeElasticSearch()
	app := routes.SetupRoutes()
	enviornment := os.Getenv("ENVIORNMENT")
	app.Run(enviornment)
}
