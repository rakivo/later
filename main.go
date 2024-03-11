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
	ADDR    = "127.0.0.1:6969"
	DB_FILE = "my.db"
)

var (
	TrustedProxies = []string{
		"127.0.0.1", "localhost:6969",
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
	log.Println("Extracting id from url:", url)
	id, err := extractYouTubeID(url); if err != nil {
		return nil, fmt.Errorf("Failed to extract YouTube ID: %v", err)
	}
	log.Println("Extracting title from id:", id)
	title, err := getYouTubeTitle(client, id, apiKey); if err != nil {
		return nil, fmt.Errorf("failed to get YouTube title: %v", err)
	}
	log.Println("Extracted title:", title)

	thumbnail := getYouTubeThumbnail(id)
	key := uuid.New()

	video := Video{}.New(title, thumbnail, url, key)
	log.Println("Adding video to the vm, video:", video)
	(*vm).AddVideo(buck, key, video)
	log.Println("Adding video to the db, video:", video)
	if err = DBaddVideo(db, []byte(buck), video); err != nil {
		return nil, fmt.Errorf("Failed to put video in database: %v", err)
	}
	return video, nil
}

/* TODO:
   1. proper frontend
	  proper readme

   2. get rid of using youtube api
	  get rid of using gin

   3. gif preview
*/

func main() {
	err := env.Load(); checkErr(err, true)
	YT_API_KEY := os.Getenv("YOUTUBE_API_KEY")

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)

	db, err := bolt.Open(DB_FILE, 0600, nil); if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	vm := VideoManager{}.New()
	if err = DBrecover(&db, &vm); err != nil {
		log.Fatal(err)
		return
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.LoadHTMLGlob("static/*")
	r.Static("/static", "./static")
	r.SetTrustedProxies(TrustedProxies)

	log.Println("Starting server on: http://127.0.0.1:6969/")

	r.GET("/", func(c *gin.Context) {
		videos, err := vm.GetVideosFromBucket(YT_BUCK);
		if err == nil {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"PostAddr": "/",
				"Videos": videos,
			})
		} else {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"PostAddr": "/",
			})
		}
	})

	client := http.Client{ Timeout: 5 * time.Second }
	r.POST("/", func(c *gin.Context) {
		url := c.PostForm("link")
		log.Println("Catched url:", url)

		_, err := addVideo(c, &db, &vm, YT_BUCK, url, &client, YT_API_KEY); checkErr(err, false)
		videos, err := vm.GetVideosFromBucket(YT_BUCK); checkErr(err, false)
		if err == nil {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"PostAddr": "/",
				"Videos": videos,
			})
		} else {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"PostAddr": "/",
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
