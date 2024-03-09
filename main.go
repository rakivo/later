package main

import (
	"os"
	"log"
	"fmt"
	"time"
	"context"
	"net/http"
	re "regexp"
	"encoding/json"
	bolt "go.etcd.io/bbolt"
	"github.com/gin-gonic/gin"
	env "github.com/joho/godotenv"
)

const (
	ADDR          = "127.0.0.1:6969"
	YT_REGEXP     = `(?:youtube\.com\/(?:[^\/\n\s]+\/\S+\/|(?:v|e(?:mbed)?)\/|\S*?[?&]v=)|youtu\.be\/)([a-zA-Z0-9_-]{11})`
	TEST_YT_URL   = "https://www.youtube.com/watch?v=86j5alcXLJE&list=PLDZ_9qD1hkzMdre6oedUdyDTgoJYq-_AY&index=4"
	GOOGLE_YT_API = "https://www.googleapis.com/youtube/v3/videos?part=snippet&id=%s&key=%s"
)

var (
	TrustedProxies = []string{
		"127.0.0.1",
	}
)

func checkErr(err error, exit bool) {
	if err != nil {
		log.Println(err)
		if exit {
			os.Exit(1)
		}
	}
}

func addSomething(c *gin.Context) {
}

func extractYouTubeID(url string) (string, error) {
	matches := re.MustCompile(YT_REGEXP).FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", fmt.Errorf("Video ID not found")
	}
	return matches[1], nil
}

func getYouTubeTitle(client *http.Client, id string, apiKey string) (string, error) {
	url := fmt.Sprintf(GOOGLE_YT_API, id, apiKey)

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

/* TODO:
   "/" POST function
   maybe feature to extract thumbnails
   proper frontend
   proper README
   LICENSE
   proper db integration
   etc
*/

func main() {
	err := env.Load(); checkErr(err, true)
	YT_API_KEY := os.Getenv("YOUTUBE_API_KEY")

	db, err := bolt.Open("my.db", 0600, nil); if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer db.Close()

	gin.SetMode(gin.DebugMode)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.SetTrustedProxies(TrustedProxies)
	r.LoadHTMLGlob("static/*")

	id, err := extractYouTubeID(TEST_YT_URL); checkErr(err, false)
	if err == nil { log.Println("ID:", id) }

	client := http.Client{Timeout: 5 * time.Second}
	title, err := getYouTubeTitle(&client, id, YT_API_KEY); checkErr(err, false)
	if err == nil { log.Println("TITLE:", title) }

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"PostAddr": "/",
		})
	})
	r.POST("/", addSomething)

	done := make(chan error, 1)
	defer close(done)
	go func() {
		err := r.Run(ADDR)
		if err != nil {
			done <- err
		}
	}()

	err = <-done; checkErr(err, true)
}
