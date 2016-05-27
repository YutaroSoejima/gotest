package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
)

type (
  metaInfo struct {
    Query string
    SearchTime float64
    TotalResults int
    PageNumber int
  }

  resultItem struct {
    Title string
    Url string
    Content string
  }

	result struct {
    Meta metaInfo
    Data []resultItem
	}
)

func search(c echo.Context) error {
  meta := metaInfo{ c.QueryParam("q"), 0.34, 12345, 1 }
  data := []resultItem {
    { "残り３日間（土日含む）頑張ろう！", "http://example.com", "hoge hoge foo bar" },
    { "好きな女優は芦田愛菜", "http://example.com", "(>_< *)" },
  }

	return c.JSON(http.StatusOK, result{ meta, data })
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/search", search)

	e.Run(standard.New(":3000"))
}
