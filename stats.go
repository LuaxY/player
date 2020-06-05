package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

func stats(c echo.Context) error {
	queryParams := c.QueryParams()
	queryParams.Add("ip", c.RealIP())

	data := make(map[string]string)

	for k, v := range queryParams {
		data[k] = v[0]
	}

	log.Info("stats ", data)

	return c.NoContent(http.StatusNoContent)
}

func statsError(c echo.Context) error {
	queryParams := c.QueryParams()
	queryParams.Add("ip", c.RealIP())

	data := make(map[string]string)

	for k, v := range queryParams {
		data[k] = v[0]
	}

	log.Info("stats_error ", data)

	return c.NoContent(http.StatusNoContent)
}
