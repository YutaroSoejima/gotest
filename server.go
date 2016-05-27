package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
)

type (
	result struct {
		Query string
	}
)

func search(c echo.Context) error {
	q := c.QueryParam("q")
	return c.JSON(http.StatusOK, result{ Query: q })
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/search", search)

	e.Run(standard.New(":3000"))
}
