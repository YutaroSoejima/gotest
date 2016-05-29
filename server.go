package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
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
		Count        int64   `json:"count"`
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
		Score           float64 `json:"_score"`
		EntryId         int64   `json:"entryId"`
		AmebaId         string  `json:"amebaId"`
		BlogTitle       string  `json:"blogTitle"`
		EntryTitle      string  `json:"entryTitle"`
		EntryContent    string  `json:"entryContent"`
		NumberOfLetters string  `json:"numberOfLetters"`
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
	fmt.Println(queryTopics)
	// [Not Yet] when requesting to elastic, also use query_topics

	var data []resultItem
	client, _ := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL("http://52.68.230.203:9200/"))
	query := elastic.NewMatchQuery("_all", queryParam)
	searchResult, _ := client.Search().Index("google").Type("ameblo").Query(query).Do()

	// [Not Yet] compairing scores
	fmt.Println(*searchResult.Hits.Hits[0].Score)

	meta := metaInfo{queryParam, 0.34, searchResult.TotalHits(), 10, 1}

	var ttyp esItem
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		if i, ok := item.(esItem); ok {
			title := i.EntryTitle + " | " + i.BlogTitle + "-" + "アメーバブログ"
			uri := "ameblo.jp/" + i.AmebaId + "/" + "entry-" + strconv.FormatInt(i.EntryId, 10) + ".html"
			text := ""
			if text_len := utf8.RuneCountInString(i.EntryContent); text_len > 140 {
				text = string([]rune(i.EntryContent)[:140]) + "..."
			}
			data = append(data, resultItem{title, uri, text})
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
