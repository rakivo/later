package main

import (
	"fmt"
	"time"
	"context"
	"net/http"
	re "regexp"
	"encoding/json"
)

const (
	YT_BUCK          = "YT_BUCKET"
	YT_REGEXP        = `^.*(youtu\.be\/|v\/|u\/\w\/|embed\/|watch\?v=|\&v=)([^#\&\?]*).*`
	YT_GET_TITLE     = "https://www.googleapis.com/youtube/v3/videos?part=snippet&id=%s&key=%s"
	YT_GET_THUMBNAIL = "https://img.youtube.com/vi/%s/hqdefault.jpg"
)

func extractYouTubeID(url string) (string, error) {
	match := re.MustCompile(YT_REGEXP).FindStringSubmatch(url)
	if len(match) >= 3 && len(match[2]) == 11 {
		return match[2], nil
	}
	return "", fmt.Errorf("Invalid YouTube URL: %s", url)
}

func getYouTubeTitle(client *http.Client, id string, apiKey string) (string, error) {
	url := fmt.Sprintf(YT_GET_TITLE, id, apiKey)

	req, err := http.NewRequest("GET", url, nil); if err != nil {
		return "", fmt.Errorf("Error creating HTTP request: %v", err)
	}
	ctx, cancel := context.WithTimeout(req.Context(), 5 * time.Second)
	defer cancel()

	resp, err := client.Do(req.WithContext(ctx)); if err != nil {
		return "", fmt.Errorf("Error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}

	var vidResp YouTubeVidResp
	if err = json.NewDecoder(resp.Body).Decode(&vidResp); err != nil {
		return "", fmt.Errorf("Error decoding JSON: %v", err)
	}; if len(vidResp.Items) == 0 {
		return "", fmt.Errorf("Video not found")
	}

	return vidResp.Items[0].Snippet.Title, nil
}

func getYouTubeThumbnail(id string) string {
	return fmt.Sprintf(YT_GET_THUMBNAIL, id)
}
