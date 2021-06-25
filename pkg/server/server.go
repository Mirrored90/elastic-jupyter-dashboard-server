package server

import (
	"github.com/elastic-jupyter-dashboard-server/pkg/configs"
	"github.com/elastic-jupyter-dashboard-server/pkg/middlewares"
	"github.com/elastic-jupyter-dashboard-server/pkg/routers"

	"github.com/gin-gonic/gin"
)

func Run() {
	serverConfig := configs.GetServerConfig()
	// gin.SetMode(serverConfig["ENV"])

	httpServer := gin.Default()
	httpServer.Use(middlewares.Cors())
	routers.RegisterRoutes(httpServer)

	serverAddr := serverConfig["HOST"] + ":" + serverConfig["PORT"]

	if err := httpServer.Run(serverAddr); err != nil {
		panic("server run error: " + err.Error())
	}
}
