package routes

import (
	searchController "github.com/adityasharma3/goSearch/cmd/search/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()
	router.GET("/:criteria/:value", searchController.Search)
	return router
}
