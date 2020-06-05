package main

import (
	"net/http"

	"github.com/labstack/echo"
)

func canPlayAds(c echo.Context) error {
	return c.JSON(http.StatusOK, struct {
		CanPlayAds bool `json:"canPlayAds"`
	}{
		CanPlayAds: true,
	})
}
