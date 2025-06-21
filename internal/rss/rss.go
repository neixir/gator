// CH3 L1 https://www.boot.dev/lessons/7347666d-7967-4c77-84c5-a0306bee8d05
package rss

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// It should fetch a feed from the given URL, and, assuming that nothing goes wrong,
// return a filled-out RSSFeed struct.
func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// Overviews
	// https://pkg.go.dev/net/http#pkg-overview
	
	// http.NewRequestWithContext
	// https://pkg.go.dev/net/http#NewRequestWithContext
	// body := io.Reader
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)

	// I set the User-Agent header to gator in the request with request.Header.Set.
	// This is a common practice to identify your program to the server.
	// https://pkg.go.dev/net/http#Header.Set
	req.Header.Set("User-Agent", "gator")

	// http.Client.Do
	// https://pkg.go.dev/net/http#Client.Do
	client := http.Client{}
	res, err := client.Do(req)

	// io.ReadAll
	// https://pkg.go.dev/io#ReadAll
	data, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	// xml.Unmarshal (works the same as json.Unmarshal)
	// https://pkg.go.dev/encoding/xml#Unmarshal
	feed := RSSFeed{}
	err = xml.Unmarshal([]byte(data), &feed)
	if err != nil {
		return &feed, err
	}

	// Use the html.UnescapeString function to decode escaped HTML entities (like &ldquo;).
	// You'll need to run the Title and Description fields
	// (of both the entire channel as well as the items) through this function.
	for _, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}

	return &feed, nil
}

func FetchFeed(feedURL string) (*RSSFeed, error) {
	return fetchFeed(context.Background(), feedURL)
}