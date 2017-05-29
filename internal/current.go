package internal

import (
	"errors"
	"log"
	"os"

	as "github.com/aerospike/aerospike-client-go"
	pb "github.com/zaquestion/current/current-service"
)

var (
	key *as.Key
	db  *as.Client
)

func DBInit() {
	var err error
	AEROSPIKE_HOST := os.Getenv("AEROSPIKE_HOST")

	db, err = as.NewClient(AEROSPIKE_HOST, 3000)
	if err != nil {
		log.Fatal(err)
	}

	k, err := as.NewKey("current", "", "loc")
	if err != nil {
		log.Fatal(err)
	}
	key = k
}

func PutLocation(loc pb.Location) error {
	if key == nil {
		return errors.New("DBInit not called")
	}
	return db.PutObject(nil, key, &loc)
}

func GetLocation() (*pb.Location, error) {
	if key == nil {
		return nil, errors.New("DBInit not called")
	}
	var loc pb.Location
	err := db.GetObject(nil, key, &loc)
	return &loc, err
}
