package main

import (
	"github.com/labstack/echo/v4"
	"github.com/salamanderman234/pos-backend/config"
	"github.com/salamanderman234/pos-backend/jobs"
	"github.com/salamanderman234/pos-backend/routes"
)

func init() {
	// init setup
	config.StartSetup()
}

func main() {
	// create new instance
	server := echo.New()
	// setup routes
	routes.RouteSetup(server)
	// start scheduled work
	jobs.StartJob()
	// start
	port := config.ApplicationPort()
	server.Logger.Fatal(
		server.Start(port),
	)
}
