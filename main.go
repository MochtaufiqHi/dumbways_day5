package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, WOrld!!")
	})

	e.GET("/about", func(c echo.Context) error {
		return c.String(200, "Ini Halaman tentang saya")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
