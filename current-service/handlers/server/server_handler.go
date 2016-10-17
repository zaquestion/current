package handler

// This file contains the Service definition, and a basic service
// implementation. It also includes service middlewares.

import (
	_ "errors"
	_ "time"

	"golang.org/x/net/context"

	_ "github.com/go-kit/kit/log"
	_ "github.com/go-kit/kit/metrics"

	pb "github.com/zaquestion/current/current-service"
	"github.com/zaquestion/current/internal"
)

// NewService returns a naïve, stateless implementation of Service.
func NewService() Service {
	internal.DBInit()
	return currentService{}
}

type currentService struct{}

// PostLocationFromBigBrother implements Service.
func (s currentService) PostLocationFromBigBrother(ctx context.Context, in *pb.PostLocationFromBigBrotherRequest) (*pb.PostLocationFromBigBrotherReply, error) {
	_ = ctx
	_ = in
	err := internal.PostLocationFromBigBrother(in)
	response := pb.PostLocationFromBigBrotherReply{}
	if err != nil {
		response.Err = err.Error()
	}
	return &response, nil
}

// GetLocation implements Service.
func (s currentService) GetLocation(ctx context.Context, in *pb.GetLocationRequest) (*pb.GetLocationReply, error) {
	_ = ctx
	_ = in
	out, err := internal.GetLocation()
	response := pb.GetLocationReply{
		Latitude:         out.Latitude,
		Longitude:        out.Longitude,
		LastUpdated:      out.Time,
		BatteryRemaining: out.Battlevel,
	}
	if err != nil {
		response.Err = err.Error()
	}
	return &response, nil
}

type Service interface {
	PostLocationFromBigBrother(ctx context.Context, in *pb.PostLocationFromBigBrotherRequest) (*pb.PostLocationFromBigBrotherReply, error)
	GetLocation(ctx context.Context, in *pb.GetLocationRequest) (*pb.GetLocationReply, error)
}
