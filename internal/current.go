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

func PostLocationFromBigBrother(in *pb.PostLocationFromBigBrotherRequest) error {
	k, err := as.NewKey("locations", "zaq", "latest")
	if err != nil {
		return err
	}

	err = db.PutObject(nil, k, in)
	return err
}

func GetLocation() (*pb.PostLocationFromBigBrotherRequest, error) {
	k, err := as.NewKey("locations", "zaq", "latest")
	if err != nil {
		return nil, err
	}

	var out pb.PostLocationFromBigBrotherRequest
	err = db.GetObject(nil, k, &out)
	return &out, err
}