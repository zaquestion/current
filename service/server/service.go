package handler

import (
	_ "errors"
	"log"
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

	bs.db, err = as.NewClient("f.iles.io", 3000)
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
	_ = in

	k, err := as.NewKey("ns", "set", "key")
	if err != nil {
		log.Println(err)
	}

	err = s.db.PutObject(nil, k, &in)
	if err != nil {
		log.Println(err)
	}

	response := pb.PostLocationFromBigBrotherReply{}
	return response, nil
}

func (s basicService) GetLocation(ctx context.Context, in pb.GetLocationRequest) (pb.GetLocationReply, error) {
	_ = ctx
	_ = in

	k, err := as.NewKey("testspace", "testset", "testkey")
	var out interface{}

	err = s.db.GetObject(nil, k, out)

	log.Println(out)

	response := pb.GetLocationReply{
		Err: err.Error(),
	}
	return response, nil
}

type Service interface {
	PostLocationFromBigBrother(ctx context.Context, in pb.PostLocationFromBigBrotherRequest) (pb.PostLocationFromBigBrotherReply, error)
	GetLocation(ctx context.Context, in pb.GetLocationRequest) (pb.GetLocationReply, error)
}
