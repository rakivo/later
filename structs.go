package main

import (
	"fmt"
	"log"
	"github.com/google/uuid"
)

type Video struct {
	Title     string
	Thumbnail string
	Url       string
	key       uuid.UUID
}

func (_ Video) New(Title string, Thumbnail string, Url string, key uuid.UUID) *Video {
	return &Video {
		Title,
		Thumbnail,
		Url,
		key,
	}
}

func (v Video) String() string {
	return fmt.Sprintf("Video{\n    Title: \t%s,\n    Thumbnail: \t%s,\n    Url: \t%s, \n}", v.Title, v.Thumbnail, v.Url)
}

type KeyVid struct {
	key uuid.UUID
	vid Video
}

func (kv KeyVid) String() string {
	return fmt.Sprintf("KeyVid{\n    key: \t-,\n    vid: \t%s    }", kv.vid.String())
}

func (_ KeyVid) New(key uuid.UUID, vid *Video) KeyVid {
	return KeyVid{
		key,
		*vid,
	}
}

func KeyVids2Videos(kvs []KeyVid) []Video {
	videos := make([]Video, len(kvs))
	for i, v := range kvs {
		videos[i] = v.vid
	}
	return videos
}

type VideoManager struct {
	order map[string][]KeyVid
	sizes map[string]uint32
}

func (_ VideoManager) New() VideoManager {
	return VideoManager {
		order: make(map[string][]KeyVid),
		sizes: make(map[string]uint32),
	}
}

func (vm *VideoManager) String() string {
	var result string
	result += "VideoManager {\n"
	for bucket, keyVids := range vm.order {
		result += fmt.Sprintf("  Bucket: %s\n", bucket)
		for i, keyVid := range keyVids {
			result += fmt.Sprintf("    Index: %d, KeyVid: %v\n", i, keyVid)
		}
	}
	result += "}\n"
	return result
}

func (vm VideoManager) GetKeyVidsFromBucket(buck string) ([]KeyVid, error) {
	if _, ok := vm.order[buck]; !ok {
		return nil, fmt.Errorf("In GetKeyVidsFromBucket: No such bucket: %s", buck)
	}
	return vm.order[buck], nil
}

func (vm VideoManager) GetVideosFromBucket(buck string) ([]Video, error) {
	log.Println("Getting video from bucket, vm.order: ", vm.order)
	if _, ok := vm.order[buck]; !ok {
		return nil, fmt.Errorf("In GetVideosFromBucket: No such bucket: %s", buck)
	}
	videos := make([]Video, vm.sizes[buck])
	for i := range vm.order[buck] {
		videos[i] = vm.order[buck][i].vid
	}
	return videos, nil
}

func (vm VideoManager) AddVideo(buck string, key uuid.UUID, video *Video) {
	if _, ok := vm.order[buck]; !ok {
		log.Println("Creating new bucket:", buck)
		vm.order[buck] = make([]KeyVid, 0)
	}
	keyvid := KeyVid{}.New(key, video)
	vm.order[buck] = append(vm.order[buck], keyvid)
	vm.sizes[buck]++
	log.Println("Added video: vm.order:", vm.order, "vm.sizes:", vm.sizes)
}

func (vm VideoManager) GetVideo(buck string, index uint32) (*Video, error) {
	log.Println("Getting video:", vm.order)
	if _, ok := vm.order[buck]; !ok {
		return nil, fmt.Errorf("In GetVideo: No such bucket: %s", buck)
	}
	size := vm.sizes[buck]
	if index > size {
		return nil, fmt.Errorf("No such video with index: %d, size of slice with bucket %s, is %d", index, buck, size)
	}
	return &vm.order[buck][index].vid, nil
}

type DBreq struct {
	Bucket []byte
	Video  *Video
}

func (_ DBreq) New(Bucket []byte, Video *Video) DBreq {
	return DBreq{
		Bucket,
		Video,
	}
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
