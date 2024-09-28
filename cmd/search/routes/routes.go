package routes

import (
	search "github.com/adityasharma3/goSearch/cmd/search/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()
	router.GET("/:criteria/:value", search.Search)
	return router
}
