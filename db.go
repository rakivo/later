package main

import (
	"fmt"
	"encoding/json"
	bolt "go.etcd.io/bbolt"
	"github.com/google/uuid"
)

func DBputVideo(db **bolt.DB, bucket []byte, video *Video) error {
	err := (*db).Update(func(tx *bolt.Tx) error {
		buck, err := tx.CreateBucketIfNotExists(bucket); if err != nil {
			return fmt.Errorf("Error creating bucket: %v", err)
		}

		keyBytes := video.key[:]
		data, err := json.Marshal(video); if err != nil {
			return fmt.Errorf("Error marshalling video data: %v", err)
		}

		if err = buck.Put(keyBytes, data); err != nil {
			return fmt.Errorf("Error putting data: %v", err)
		}
		return nil
	})
	return err
}

func DBgetVideo(db **bolt.DB, bucket []byte, key uuid.UUID) (Video, error) {
	var video Video
	err := (*db).View(func(tx *bolt.Tx) error {
		buck := tx.Bucket(bucket); if buck == nil {
			return fmt.Errorf("Bucket %s not found", bucket)
		}

		keyBytes := key[:]
		data := buck.Get(keyBytes); if data == nil {
			return fmt.Errorf("Key not found: %s", key)
		}

		if err := json.Unmarshal(data, &video); err != nil {
			return fmt.Errorf("Error unmarshalling video data: %v", err)
		}
		return nil
	});
	return video, err
}
