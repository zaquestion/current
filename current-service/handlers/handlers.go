package handlers

import (
	"errors"

	"golang.org/x/net/context"

	pb "github.com/zaquestion/current/current-service"
	"github.com/zaquestion/current/internal"
)

// NewService returns a naïve, stateless implementation of Service.
func NewService() pb.CurrentServer {
	internal.DBInit()
	return currentService{}
}

type currentService struct{}

// PostLocationBigBrother implements Service.
func (s currentService) PostLocationBigBrother(ctx context.Context, in *pb.PostLocationBigBrotherRequest) (*pb.Error, error) {
	loc := pb.Location{
		Latitude:         in.Latitude,
		Longitude:        in.Longitude,
		LastUpdated:      in.Time,
		BatteryRemaining: in.Battlevel,
	}
	err := internal.PutLocation(loc)
	response := pb.Error{}
	if err != nil {
		response.Err = err.Error()
	}
	return &response, nil
}

// PostLocationTasker implements Service.
func (s currentService) PostLocationTasker(ctx context.Context, in *pb.PostLocationTaskerRequest) (*pb.Error, error) {
	response := pb.Error{}
	if len(in.Location) < 2 {
		err := errors.New("No location provided")
		response.Err = err.Error()
		return &response, err
	}
	loc := pb.Location{
		Latitude:         in.Location[0],
		Longitude:        in.Location[1],
		LastUpdated:      in.Time,
		BatteryRemaining: in.Battery,
	}
	err := internal.PutLocation(loc)
	if err != nil {
		response.Err = err.Error()
	}
	return &response, err
}

// GetLocation implements Service.
func (s currentService) GetLocation(ctx context.Context, in *pb.GetLocationRequest) (*pb.Location, error) {
	loc, err := internal.GetLocation()
	if err != nil {
		loc.Err = err.Error()
	}
	return loc, err
}
