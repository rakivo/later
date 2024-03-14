package main

import (
	"fmt"
	"log"
	"time"
	"context"
	"net/http"
	"encoding/json"
)

const (
	YT_BUCK     = "YT_BUCKET"
	YT_GET_DATA = "https://noembed.com/embed?url=%s"
)

func getYouTubeTitleAndThumbnail(client **http.Client, url string) (string, string, error) {
	url = fmt.Sprintf(YT_GET_DATA, url)
	log.Println("URL:", url)

	req, err := http.NewRequest("GET", url, nil); if err != nil {
		return "", "", fmt.Errorf("Error creating HTTP request: %v", err)
	}
	ctx, cancel := context.WithTimeout(req.Context(), 5 * time.Second)
	defer cancel()

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3");
	resp, err := (*client).Do(req.WithContext(ctx)); if err != nil {
		return "", "", fmt.Errorf("Error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}

	var ytResp YTresp
	if err = json.NewDecoder(resp.Body).Decode(&ytResp); err != nil {
		return "", "", fmt.Errorf("Error decoding JSON: %v", err)
	};

	return ytResp.Title, ytResp.ThumbUrl, nil
}
