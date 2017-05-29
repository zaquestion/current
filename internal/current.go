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

func PutLocation(loc pb.Location) error {
	k, err := as.NewKey("locations", "zaq", "latest")
	if err != nil {
		return err
	}

	err = db.PutObject(nil, k, loc)
	return err
}

func GetLocation() (*pb.Location, error) {
	k, err := as.NewKey("locations", "zaq", "latest")
	if err != nil {
		return nil, err
	}

	var loc pb.Location
	err = db.GetObject(nil, k, &loc)
	return &loc, err
}
