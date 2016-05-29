package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/franela/goreq"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"

	"gopkg.in/olivere/elastic.v3"
)

type (
	// For responce
	metaInfo struct {
		Query        string  `json:"query"`
		SearchTime   float64 `json:"searchTime"`
		Count        int     `json:"count"`
		CountPerPage int     `json:"countPerPage"`
		PageNumber   int     `json:"pageNumber"`
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

	// .../prod/classify?word={QueryParameter.Word}
	QueryParameter struct {
		Word string
	}
)

// return topics for query(::query_topics)
func getTopics(query string) map[string]string {
	word := QueryParameter{Word: query}
	response, _ := goreq.Request{
		Uri:         "https://jpdtd1hnzf.execute-api.ap-northeast-1.amazonaws.com/prod/classify",
		QueryString: word,
	}.Do()

	if len(response.Uri) > 0 {
		body, _ := response.Body.ToString()
		body = strings.Trim(body, "{}")
		topics := strings.Split(body, ",")
		query_topics := make(map[string]string)
		for _, topic := range topics {
			elm := strings.Split(topic, ":")
			elm[0] = strings.Trim(elm[0], " ")
			query_topics[elm[0]] = elm[1]
		}
		// {topic: probability, topic: probability, ...}
		return query_topics
	}

	return nil
}

func search(c echo.Context) error {
	queryParam := c.QueryParam("query")
	queryTopics := getTopics(queryParam)
	fmt.Println(reflect.TypeOf(queryTopics))
	// [REMIND] when requesting to elastic, also use query_topics

	meta := metaInfo{queryParam, 0.34, 12345, 10, 1}

	var data []resultItem
	client, _ := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL("http://52.68.230.203:9200/"))
	query := elastic.NewMatchQuery("_all", queryParam)
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
