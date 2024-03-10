package main

import (
	"fmt"
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
	return fmt.Sprintf("Video{\n    Title: \t%s,\n    Thumbnail: \t%s,\n    Url: \t%s, \n    Key: \t%s\n}", v.Title, v.Thumbnail, v.Url, v.key)
}

type KeyVid struct {
	key uuid.UUID
	vid Video
}

func (kv KeyVid) String() string {
	return fmt.Sprintf("KeyVid{\n    key: \t%s,\n    vid: \t%s    }", kv.key.String(), kv.vid.String())
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
		return nil, fmt.Errorf("ERROR: No such bucket: %s", buck)
	}
	return vm.order[buck], nil
}

func (vm VideoManager) AddVideo(buck string, key uuid.UUID, video *Video) {
	if _, ok := vm.order[buck]; !ok {
		vm.order[buck] = make([]KeyVid, 0)
	}
	keyvid := KeyVid{}.New(key, video)
	vm.order[buck] = append(vm.order[buck], keyvid)
	vm.sizes[buck]++
}

func (vm VideoManager) GetVideo(buck string, index uint32) (*Video, error) {
	if _, ok := vm.order[buck]; !ok {
		return nil, fmt.Errorf("No such bucket: %s", buck)
	}
	size := vm.sizes[buck]
	if index > size {
		return nil, fmt.Errorf("No such video with index: %d, size of slice with bucket %s, is %d", index, buck, size)
	}
	return &vm.order[buck][index].vid, nil
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
