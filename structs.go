package main

import (
	"fmt"
	"github.com/google/uuid"
)

type Video struct {
	Title     string
	Thumbnail string
	key       uuid.UUID
}

func (_ Video) New(Title string, Thumbnail string, key uuid.UUID) *Video {
	return &Video {
		Title,
		Thumbnail,
		key,
	}
}

func (v Video) String() string {
    return fmt.Sprintf("Video{\n    Title: \t%s, \n    Thumbnail: \t%s, \n    Key: \t%s\n}", v.Title, v.Thumbnail, v.key)
}

type YouTubeSnippet struct {
	Title string `json:"title"`
}

type YouTubeItem struct {
	Snippet YouTubeSnippet `json:"snippet"`
}

type YouTubeVidResp struct {
	Items []YouTubeItem `json:"items"`
}
