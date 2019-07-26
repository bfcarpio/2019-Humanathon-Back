package main

import (
	"github.com/labstack/echo"
)

func main() {
	// Create a new instance of Echo
	e := echo.New()

	e.GET("/hello", func(c echo.Context) error { return c.JSON(200, "Hello World") })

	// Start as a web server
	e.Start(":8080")
}