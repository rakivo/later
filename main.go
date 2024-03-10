package main

import (
	"os"
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

func addSomething(c *gin.Context) {
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

	id, err := extractYouTubeID(TEST_YT_URL); checkErr(err, false)

	client := http.Client{ Timeout: 5 * time.Second }
	title, err := getYouTubeTitle(&client, id, YT_API_KEY); checkErr(err, false)

	thumbnail := getYouTubeThumbnail(id)

	video := Video{}.New(title, thumbnail, uuid.New())

	err = DBputVideo(&db, []byte(YT_BUCK), video)

	*video, err = DBgetVideo(&db, []byte(YT_BUCK), video.key); if err == nil {
		log.Println("------------------------")
		log.Println(video.String())
		log.Println("------------------------")
	}

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

	checkErr(<-done, true)
}
