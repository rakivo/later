package main

import (
	"os"
	"fmt"
	"log"
	"time"
	"bufio"
	"strings"
	"net/http"
	"html/template"
	bolt "go.etcd.io/bbolt"
	"github.com/google/uuid"
)

const (
	ADDR    = "127.0.0.1:6969"
	DB_FILE = "my.db"
)

// get video's id, title, thumbnail; add created video to db and vm
func addVideo (
	db **bolt.DB,
	vm *VideoManager,
	buck string,
	url string,
	client *http.Client,
	apiKey string,
	dbChan chan DBreq) (*Video, error) {

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

   3. gif preview
*/

func main() {
	YT_API_KEY := os.Getenv("LATER_YOUTUBE_API_KEY")
	LATER_DIR  := os.Getenv("LATER_PROJECT_DIR")

	db, err := bolt.Open(LATER_DIR + DB_FILE, 0600, nil); if err != nil {
		log.Fatal(err)
		return
	}; defer db.Close()

	vm := VideoManager{}.New()
	y, n := checkArgs_(os.Args, len(os.Args))
	if y {
		err := DBrecover(&db, &vm); checkErr_(err, true)
	} else if !y && !n {
		if ask(bufio.NewReader(os.Stdin)) {
			DBrecover(&db, &vm); checkErr_(err, true)
		}
	}

	fmt.Println("Starting server on: http://" + ADDR)

	client := http.Client{ Timeout: 5 * time.Second }
	dbChan := make(chan DBreq) // channel to send requests to the DB

	go func() {
		for req := range dbChan {
			if err := DBaddVideo(&db, []byte(req.Bucket), req.Video)
			err != nil {
				log.Printf("Failed to put video in database: %v", err)
			}
		}
	}()

	tmpl, err := template.ParseFiles(LATER_DIR + "static/index.html"); checkErr_(err, true)
	http.Handle(LATER_DIR + "/static/", http.StripPrefix(LATER_DIR + "/static/", http.FileServer(http.Dir(LATER_DIR + "/static"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			url := r.FormValue("link")
			log.Println("Catched url:", url)

			_, err := addVideo(
				&db,
				&vm, YT_BUCK,
				url, &client,
				YT_API_KEY, dbChan); checkErr_(err, false)
			checkGetAndRender_(&tmpl, &w, &vm, YT_BUCK); checkErr_(err, false)
		} else {
			checkGetAndRender_(&tmpl, &w, &vm, YT_BUCK)
		}
	})
	checkErr_(http.ListenAndServe(ADDR, nil), true)
}

func checkArgs_(args []string, argsLen int) (bool, bool) {
	ifY := false
	ifN := false
	if argsLen == 2 {
		trimmed := strings.TrimSpace(args[1])
		ifY = strings.Compare(trimmed, "y") == 0
		ifN = strings.Compare(trimmed, "n") == 0
		return ifY, ifN
	}; if argsLen >= 2 {
		fmt.Println("Usage: ./later <y/n>")
		fmt.Println("y means recover videos from the previous sessions")
		fmt.Println("n means do not")
		fmt.Println("By the way, these flags are optional")
	}
	return false, false
}

func checkErr_(err error, exit bool) {
	if err != nil {
		log.Println("ERROR:", err)
		if exit {
			os.Exit(1)
		}
	}
}

func checkGetAndRender_(tmpl **template.Template, w *http.ResponseWriter, vm *VideoManager, buck string) {
	videos, err := vm.GetVideosFromBucket(buck);
	if err == nil {
		if err = (*tmpl).Execute(*w, struct {
			PostAddr string
			Videos []Video
		} {
			PostAddr: "/",
			Videos: videos,
		}); err != nil {
			http.Error(*w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		if err = (*tmpl).Execute(*w, struct {
			PostAddr string
			Videos []Video
		} {
			PostAddr: "/",
			Videos: []Video{},
		}); err != nil {
			http.Error(*w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func ask(ior *bufio.Reader) bool {
	for {
		fmt.Printf("Recover videos from previous sessions? [y/n] ")
		ioans, err := ior.ReadString('\n'); if err != nil {
			fmt.Println("ERROR reading input:", err)
			continue
		}
		ioans = strings.TrimSpace(ioans)
		if len(ioans) == 0 {
			fmt.Println("Enter y or n")
			continue
		}
		fields := strings.Fields(ioans)
		if len(fields) > 1 {
			fmt.Println("Enter y or n")
			continue
		}

		ifY := strings.Compare(fields[0], "y")
		ifN := strings.Compare(fields[0], "n")
		if ifY != 0 && ifN != 0 {
			fmt.Println("Enter y or n")
			continue
		} else if ifY == 0 {
			return true
		} else if ifN == 0 {
			return false
		}
	}
}
