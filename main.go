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

// get video's id, title, thumbnail; add created video to db and vm
func addVideo(c *gin.Context, db **bolt.DB, vm *VideoManager, buck string, url string, client *http.Client, apiKey string, dbChan chan DBreq) (*Video, error) {
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
	dbChan <- DBreq{}.New([]byte(buck), video)
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
	err := env.Load(); checkErr_(err, true)
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
	r.SetTrustedProxies("127.0.0.1")

	log.Println("Starting server on: http://" + ADDR)

	client := http.Client{ Timeout: 5 * time.Second }
	dbChan := make(chan DBreq) // channel to send requests to the DB

	go func() {
		for req := range dbChan {
			err := DBaddVideo(&db, []byte(req.Bucket), req.Video)
			if err != nil {
				log.Printf("Failed to put video in database: %v", err)
			}
		}
	}()

	r.GET("/", func(c *gin.Context) {
		checkGetAndRender_(c, &vm, YT_BUCK)
	})
	r.POST("/", func(c *gin.Context) {
		url := c.PostForm("link")
		log.Println("Catched url:", url)

		_, err := addVideo(c, &db, &vm, YT_BUCK, url, &client, YT_API_KEY, dbChan); checkErr_(err, false)
		checkGetAndRender_(c, &vm, YT_BUCK); checkErr_(err, false)
	})

	checkErr_(r.Run(ADDR), true)
}

func checkErr_(err error, exit bool) {
	if err != nil {
		log.Println("ERROR:", err)
		if exit {
			os.Exit(1)
		}
	}
}

func checkGetAndRender_(c *gin.Context, vm *VideoManager, buck string) {
	videos, err := vm.GetVideosFromBucket(buck);
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
}
