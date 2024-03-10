package main

import (
	"os"
	"fmt"
	"log"
	"time"
	"net/http"
	bolt "go.etcd.io/bbolt"
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
	env "github.com/joho/godotenv"
)

const (
	ADDR = "127.0.0.1:6969"
)

var (
	TrustedProxies = []string{
		"127.0.0.1",
	}
)

func checkErr(err error, exit bool) {
	if err != nil {
		log.Println("ERROR:", err)
		if exit {
			os.Exit(1)
		}
	}
}

// get video's id, title, thumbnail; add created video to db and vm
func addVideo(c *gin.Context, db **bolt.DB, vm *VideoManager, buck string, url string, client *http.Client, apiKey string) (*Video, error) {
	id, err := extractYouTubeID(url); if err != nil {
		return nil, fmt.Errorf("Failed to extract YouTube ID: %v", err)
	}
	log.Println("Extracted id:", id)
	title, err := getYouTubeTitle(client, id, apiKey); if err != nil {
		return nil, fmt.Errorf("failed to get YouTube title: %v", err)
	}
	log.Println("Extracted title:", title)

	thumbnail := getYouTubeThumbnail(id)
	key := uuid.New()

	video := Video{}.New(title, thumbnail, key)
	if err = DBaddVideo(db, []byte(buck), video); err != nil {
		return nil, fmt.Errorf("Failed to put video in database: %v", err)
	}
	vm.AddVideo(buck, key, video)
	return video, nil
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
		return;
	}
	defer db.Close()

	gin.SetMode(gin.DebugMode)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.LoadHTMLGlob("static/*")
	r.Static("/static", "./static")
	r.SetTrustedProxies(TrustedProxies)

	client := http.Client{ Timeout: 5 * time.Second }
	vm := VideoManager{}.New()

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"PostAddr": "/",
		})
	})
	r.POST("/", func(c *gin.Context) {
		url := c.PostForm("link")
		log.Println("Catched url:", url)

		latestvideo, err := addVideo(c, &db, vm, YT_BUCK, url, &client, YT_API_KEY); checkErr(err, false)
		*latestvideo, err = DBgetVideo(&db, []byte(YT_BUCK), (*latestvideo).key); if err == nil {
			log.Println("------DB GET VIDEO------")
			log.Println(latestvideo.String())
		} else {
			log.Println("DB ERROR:", err)
		}

		kvs, err := vm.GetKeyVidsFromBucket(YT_BUCK); checkErr(err, false)
		if len(kvs) == 0 {
			log.Println("No videos found")
			c.HTML(http.StatusOK, "index.html", gin.H{
				"PostAddr": "/",
				"Videos":   nil,
			})
		} else {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"PostAddr": "/",
				"Videos": KeyVids2Videos(kvs),
			})
		}
	})

	done := make(chan error, 1)
	defer close(done)
	go func() {
		if err := r.Run(ADDR); err != nil {
			done <- err
		}
	}()

	checkErr(<-done, true)
}
