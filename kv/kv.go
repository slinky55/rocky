package kv

import bolt "go.etcd.io/bbolt"

var Conn *bolt.DB
var ChallengeBucket *bolt.Bucket
var SessionBucket *bolt.Bucket

func Init() (err error) {
	Conn, err = bolt.Open("rocky.kv", 0600, nil)

	err = Conn.Update(func(tx *bolt.Tx) error {
		ChallengeBucket, err = tx.CreateBucketIfNotExists([]byte("challenges"))
		return err
	})
	if err != nil {
		return err
	}

	err = Conn.Update(func(tx *bolt.Tx) error {
		SessionBucket, err = tx.CreateBucketIfNotExists([]byte("sessions"))
		return err
	})
	if err != nil {
		return err
	}

	return nil
}
