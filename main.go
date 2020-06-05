package main

import (
	"log"
	"mime"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	domain string
)

func main() {
	port, _ := os.LookupEnv("PORT")
	domain, _ = os.LookupEnv("DOMAIN")

	if port == "" {
		port = "80"
	}

	if domain == "" {
		log.Fatal("DOMAIN env variable need to be set")
	}

	_ = mime.AddExtensionType(".js", "application/javascript")
	_ = mime.AddExtensionType(".css", "text/css")

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	e.GET("/", home)

	e.GET("/player/*", player)
	e.GET("/v1/jwplayer6/ping.gif", stats)
	e.GET("/v1/error/ping.gif", statsError)
	e.GET("/canPlayAds.json", canPlayAds)

	e.Logger.Fatal(e.Start(":" + port))
}

func home(c echo.Context) error {
	return c.String(http.StatusOK, "JWPlayer")
}
