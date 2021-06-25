package routers

import (
	"github.com/elastic-jupyter-dashboard-server/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	//router.GET("/", controllers.IndexHome)

	router.GET("/notebooks", controllers.GetNotebooks)
	router.DELETE("/notebooks", controllers.DeleteNotebook)
	router.POST("/notebooks/create", controllers.CreateNotebook)
}
