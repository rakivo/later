package main

import (
	"fmt"
	"log"
	"encoding/json"
	bolt "go.etcd.io/bbolt"
	"github.com/google/uuid"
)

func DBaddVideo(db **bolt.DB, bucket []byte, video *Video) error {
	return (*db).Update(func(tx *bolt.Tx) error {
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

func DBrecover(db **bolt.DB, vm *VideoManager) error {
	return (*db).View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(bucketName []byte, bucket *bolt.Bucket) error {
			bucketNameStr := string(bucketName)
			log.Printf("Recovering bucket: %s", bucketNameStr)

			var keyVids []KeyVid
			var size uint32 = 0
			bucket.ForEach(func(k, v []byte) error {
				var video Video
				if err := json.Unmarshal(v, &video); err != nil {
					return fmt.Errorf("Error unmarshalling video: %v", err)
				}
				log.Println("Recovered video:", video)
				uuidKey, err := uuid.FromBytes(k); if err != nil {
					return fmt.Errorf("Error converting key to UUID: %v", err)
				}
				keyVid := KeyVid{}.New(uuidKey, &video)
				keyVids = append(keyVids, keyVid)
				size++
				return nil
			})
			(*vm).order[bucketNameStr] = keyVids
			(*vm).sizes[bucketNameStr] = size
			return nil
		})
	})
}
