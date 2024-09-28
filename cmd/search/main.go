package main

import (
	"github.com/adityasharma3/goSearch/cmd/search/routes"
)

func main() {
	app := routes.SetupRoutes()
	app.Run("localhost:8080")
}
