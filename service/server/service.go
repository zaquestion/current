package handler

import (
	_ "errors"
	"log"
	"os"
	_ "time"

	"golang.org/x/net/context"

	_ "github.com/go-kit/kit/log"
	_ "github.com/go-kit/kit/metrics"

	as "github.com/aerospike/aerospike-client-go"
	pb "github.com/zaquestion/current/service/DONOTEDIT/pb"
)

func NewBasicService() Service {
	var err error
	var bs basicService

	AEROSPIKE_HOST := os.Getenv("AEROSPIKE_HOST")

	bs.db, err = as.NewClient(AEROSPIKE_HOST, 3000)
	if err != nil {
		log.Fatal(err)
	}
	return bs
}

type basicService struct {
	db *as.Client
}

func (s basicService) PostLocationFromBigBrother(ctx context.Context, in pb.PostLocationFromBigBrotherRequest) (pb.PostLocationFromBigBrotherReply, error) {
	_ = ctx

	k, err := as.NewKey("locations", "zaq", "latest")
	if err != nil {
		log.Println(err)
		response := pb.PostLocationFromBigBrotherReply{
			Err: err.Error(),
		}
		return response, err
	}

	err = s.db.PutObject(nil, k, &in)
	if err != nil {
		log.Println(err)
		response := pb.PostLocationFromBigBrotherReply{
			Err: err.Error(),
		}
		return response, err
	}

	response := pb.PostLocationFromBigBrotherReply{}
	return response, nil
}

func (s basicService) GetLocation(ctx context.Context, in pb.GetLocationRequest) (pb.GetLocationReply, error) {
	_ = ctx
	_ = in

	k, err := as.NewKey("locations", "zaq", "latest")
	var response pb.GetLocationReply
	var out pb.PostLocationFromBigBrotherRequest

	err = s.db.GetObject(nil, k, &out)
	if err != nil {
		log.Println(err)
		response.Err = err.Error()
		return response, err
	}

	response.Latitude = out.Latitude
	response.Longitude = out.Longitude

	return response, nil
}

type Service interface {
	PostLocationFromBigBrother(ctx context.Context, in pb.PostLocationFromBigBrotherRequest) (pb.PostLocationFromBigBrotherReply, error)
	GetLocation(ctx context.Context, in pb.GetLocationRequest) (pb.GetLocationReply, error)
}
