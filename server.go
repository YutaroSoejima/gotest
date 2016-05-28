package main

import (
	"net/http"
	"reflect"
	"unicode/utf8"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"

	"gopkg.in/olivere/elastic.v3"
)

type (
	// For responce
	metaInfo struct {
		Query      string  `json:"query"`
		SearchTime float64 `json:"searchTime"`
		Count      int     `json:"count"`
		PageNumber int     `json:"pageNumber"`
	}

	resultItem struct {
		Title   string `json:"title"`
		URL     string `json:"url"`
		Content string `json:"content"`
	}

	result struct {
		Meta metaInfo     `json:"meta"`
		Data []resultItem `json:"data"`
	}

	// For Serialization
	esItem struct {
		URL       string `json:"url"`
		WholeText string
		PageRank  float64
	}
)

func search(c echo.Context) error {
	meta := metaInfo{c.QueryParam("query"), 0.34, 12345, 1}
	var data []resultItem
	client, _ := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL("http://52.68.230.203:9200/"))
	query := elastic.NewMatchQuery("wholeText", c.QueryParam("query"))
	searchResult, _ := client.Search().Index("google").Query(query).Do()

	var ttyp esItem
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		if i, ok := item.(esItem); ok {
			text := ""
			if text_len := utf8.RuneCountInString(i.WholeText); text_len > 140 {
				text = string([]rune(i.WholeText)[:140]) + "..."
			}
			data = append(data, resultItem{"titleはまだない", i.URL, text})
		}
	}

	c.Response().Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")

	return c.JSON(http.StatusOK, result{meta, data})
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/search", search)

	e.Run(standard.New(":3000"))
}
