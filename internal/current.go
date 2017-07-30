package internal

import (
	"log"
	"os"

	as "github.com/aerospike/aerospike-client-go"
	pb "github.com/zaquestion/current/current-service"
)

var (
	db *as.Client
)

func DBInit() {
	var err error
	AEROSPIKE_HOST := os.Getenv("AEROSPIKE_HOST")

	db, err = as.NewClient(AEROSPIKE_HOST, 3000)
	if err != nil {
		log.Fatal(err)
	}
}

func PutLocation(secret string, loc pb.Location) error {
	key, err := as.NewKey("current", "location", secret)
	if err != nil {
		return err
	}
	return db.PutObject(nil, key, &loc)
}

func GetLocation(secret string) (*pb.Location, error) {
	key, err := as.NewKey("current", "location", secret)
	if err != nil {
		return nil, err
	}
	var loc pb.Location
	err = db.GetObject(nil, key, &loc)
	return &loc, err
}
