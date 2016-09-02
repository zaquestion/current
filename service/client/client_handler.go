package clienthandler

import (
	pb "github.com/zaquestion/current/service/DONOTEDIT/pb"
)

func PostLocationFromBigBrother(LatitudePostLocationFromBigBrother float64, LongitudePostLocationFromBigBrother float64, AccuracyPostLocationFromBigBrother float64, AltitudePostLocationFromBigBrother float64, BearingPostLocationFromBigBrother float64, SpeedPostLocationFromBigBrother float64) (pb.PostLocationFromBigBrotherRequest, error) {
	request := pb.PostLocationFromBigBrotherRequest{
		Latitude:  LatitudePostLocationFromBigBrother,
		Longitude: LongitudePostLocationFromBigBrother,
		Accuracy:  AccuracyPostLocationFromBigBrother,
		Altitude:  AltitudePostLocationFromBigBrother,
		Bearing:   BearingPostLocationFromBigBrother,
		Speed:     SpeedPostLocationFromBigBrother}
	return request, nil
}

func GetLocation() (pb.GetLocationRequest, error) {
	request := pb.GetLocationRequest{}
	return request, nil
}
