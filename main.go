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
		log.Println(err)
		if exit {
			os.Exit(1)
		}
	}
}

func addVideo(c *gin.Context, db **bolt.DB, vm *VideoManager, buck string, url string, client *http.Client, apiKey string) (*Video, error) {
	id, err := extractYouTubeID(url); if err != nil {
		return nil, fmt.Errorf("failed to extract YouTube ID: %v", err)
	}
	title, err := getYouTubeTitle(client, id, apiKey); if err != nil {
		return nil, fmt.Errorf("failed to get YouTube title: %v", err)
	}

	thumbnail := getYouTubeThumbnail(id)
	key := uuid.New()
	video := Video{}.New(title, thumbnail, key)
	vm.AddVideo(buck, key, video)
	if err = DBputVideo(db, []byte(buck), video); err != nil {
		return nil, fmt.Errorf("failed to put video in database: %v", err)
	}
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
		os.Exit(1)
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

	latestvideo := Video{}.New("", "", uuid.New())

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"PostAddr": "/",
		})
	})
	r.POST("/", func(c *gin.Context) {
		latestvideo, err = addVideo(c, &db, vm, YT_BUCK, TEST_YT_URL, &client, YT_API_KEY)
		// JUST FOR TEST:
		*latestvideo, err = DBgetVideo(&db, []byte(YT_BUCK), (*latestvideo).key); if err == nil {
			log.Println("------DB GET VIDEO------")
			log.Println(latestvideo.String())
		} else {
			log.Println("ERROR:", err)
		}
		ordervid, err := vm.GetVideo(YT_BUCK, 0); if err == nil {
			log.Println("-----ORDER GET VIDEO----")
			log.Println(ordervid.String())
		} else {
			log.Println("ERROR:", err)
		}
	})

	latestvideo.Thumbnail = ";";

	done := make(chan error, 1)
	defer close(done)
	go func() {
		if err := r.Run(ADDR); err != nil {
			done <- err
		}
	}()

	checkErr(<-done, true)
}
