package handlers

import (
	"errors"
	"time"

	"github.com/bradfitz/latlong"
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
func (s currentService) PostLocationBigBrother(ctx context.Context, in *pb.PostLocationBigBrotherRequest) (*pb.Empty, error) {
	return nil, errors.New("Currently unsupported")
	response := pb.Empty{}
	loc := pb.Location{
		Latitude:    in.Latitude,
		Longitude:   in.Longitude,
		Speed:       in.Speed,
		LastUpdated: in.Time,
		Battery:     in.Battlevel,
	}
	err := internal.PutLocation("", loc)
	return &response, err
}

// PostLocationTasker implements Service.
func (s currentService) PostLocationTasker(ctx context.Context, in *pb.PostLocationTaskerRequest) (*pb.Empty, error) {
	response := pb.Empty{}
	if len(in.Location) < 2 {
		err := errors.New("No location provided")
		return &response, err
	}
	lat := in.Location[0]
	long := in.Location[1]
	zone, err := time.LoadLocation(latlong.LookupZoneName(lat, long))
	if err != nil {
		return &response, err
	}

	datetime, err := time.ParseInLocation("1-2-06 15.04", in.DateTime, zone)
	if err != nil {
		return &response, err
	}
	loc := pb.Location{
		Latitude:    lat,
		Longitude:   long,
		Charging:    in.Charging,
		Speed:       in.Speed,
		LastUpdated: datetime.UTC().Format("2006-01-02T15:04:05.00Z"),
		Battery:     in.Battery,
	}
	err = internal.PutLocation(in.Secret, loc)
	return &response, err
}

// GetLocation implements Service.
func (s currentService) GetLocation(ctx context.Context, in *pb.GetLocationRequest) (*pb.Location, error) {
	return internal.GetLocation(in.Secret)
}
