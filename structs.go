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

func (vid Video) New(Title string, Thumbnail string, key uuid.UUID) *Video {
	return &Video {
		Title,
		Thumbnail,
		key,
	}
}

func (v Video) String() string {
	return fmt.Sprintf("Video{Title: %s, Thumbnail: %s, Key: %s}", v.Title, v.Thumbnail, v.key)
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
