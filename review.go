package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type JsonResponse struct {
	Feed Feed `json:"feed"`
}

type Feed struct {
	Entry []Entry `json:"entry"`
}

type Entry struct {
	Author  Author  `json:"author"`
	Updated Updated `json:"updated"`
	Rating  Label   `json:"im:rating"`
	Title   Label   `json:"title"`
	Content Label   `json:"content"`
}

type Author struct {
	Uri  Uri   `json:"author"`
	Name Label `json:"name"`
}

type Updated struct {
	Label time.Time `json:"label"`
}

type Uri struct {
	Uri string `json:"uri"`
}

type Label struct {
	Label string `json:"label"`
}

type SlackBlock struct {
	Blocks []SlackEntity `json:"blocks"`
}

type SlackEntity struct {
	Type string    `json:"type"`
	Text SlackText `json:"text"`
}

type SlackText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func main() {
	flag.Parse()
	id := flag.Arg(0)
	webhookUrl := flag.Arg(1)

	if len(id) == 0 {
		log.Fatal("id is Empty")
	}

	if len(webhookUrl) == 0 {
		log.Fatal("WebhookUrl is Empty")
	}

	entries, err := fetchEntries(id)
	if err != nil {
		log.Fatal("Fail to Fetch,", err)
	}

	yesterday := time.Now().AddDate(0, -1, -1)
	fmt.Println(yesterday)

	for i := 0; i < len(entries); i++ {
		entry := entries[i]
		if !yesterday.After(entry.Updated.Label) {
			postSlack(webhookUrl, entry)
		}
	}
}

func fetchEntries(id string) (entries []Entry, err error) {
	// API呼び出し
	res, err := http.Get("https://itunes.apple.com/jp/rss/customerreviews/id=" + id + "/sortBy=mostRecent/json")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// jsonを読み込む
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// JSONデコード
	var response JsonResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return response.Feed.Entry, nil
}

func postSlack(url string, entry Entry) (err error) {
	title := entry.Title.Label
	authorName := entry.Author.Name.Label
	content := entry.Content.Label
	rating, err := strconv.Atoi(entry.Rating.Label)
	if err != nil {
		return err
	}
	rateString := ""
	for i := 0; i < rating; i++ {
		rateString = rateString + ":star:"
	}

	text := "*" + title + "*\n" + rateString + " by " + authorName + "\n" + content
	slackText := SlackText{Type: "mrkdwn", Text: text}
	slackEntity := SlackEntity{Type: "section", Text: slackText}
	slackEntities := []SlackEntity{}
	slackEntities = append(slackEntities, slackEntity)
	slackBlock := SlackBlock{Blocks: slackEntities}
	jsonValue, _ := json.Marshal(slackBlock)

	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	fmt.Println(response)
	if err != nil {
		return err
	}
	return nil
}
