package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
)

type (
  metaInfo struct {
    Query string `json:"query"`
    SearchTime float64 `json:"searchTime"`
    Count int `json:"count"`
    PageNumber int `json:"pageNumber"`
  }

  resultItem struct {
    Title string `json:"title"`
    URL string `json:"url"`
    Content string `json:"content"`
  }

	result struct {
    Meta metaInfo `json:"meta"`
    Data []resultItem `json:"data"`
	}
)

func search(c echo.Context) error {
  meta := metaInfo{ c.QueryParam("q"), 0.34, 12345, 1 }
  data := []resultItem {
    { "残り３日間（土日含む）頑張ろう！", "http://example.com", "hoge hoge foo bar" },
    { "好きな女優は芦田愛菜", "http://example.com", "ロリコンではなく父性本能" },
  }

  c.Response().Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
  c.Response().Header().Set("Access-Control-Allow-Origin", "*")

	return c.JSON(http.StatusOK, result{ meta, data })
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/search", search)

	e.Run(standard.New(":3000"))
}
