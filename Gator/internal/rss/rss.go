package rss

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

func FetchFeed(ctx context.Context, feedURL string) (feed *RSSFeed, err error) {
	// Get new client
	client := &http.Client{}

	// Build request
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	// Let the remote server know we're a bot
	req.Header.Add("User-Agent", "gator")

	// Do callout
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var parsedResponse RSSFeed
	err = xml.Unmarshal(responseBody, &parsedResponse)

	// Unescape the title and description fields
	parsedResponse.Channel.Title = html.UnescapeString(parsedResponse.Channel.Title)
	parsedResponse.Channel.Description = html.UnescapeString(parsedResponse.Channel.Description)
	for i, item := range parsedResponse.Channel.Item {
		parsedResponse.Channel.Item[i].Title = html.UnescapeString(item.Title)
		parsedResponse.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}

	return &parsedResponse, nil
}
