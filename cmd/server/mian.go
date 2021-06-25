package main

import (
	_ "github.com/elastic-jupyter-dashboard-server/pkg/driver"
	"github.com/elastic-jupyter-dashboard-server/pkg/server"
)

func main() {
	server.Run()
}
