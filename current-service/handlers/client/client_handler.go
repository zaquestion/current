package clienthandler

import (
	pb "github.com/zaquestion/current/current-service"
)

// PostLocationFromBigBrother implements Service.
func PostLocationFromBigBrother(LatitudePostLocationFromBigBrother float64, LongitudePostLocationFromBigBrother float64, AccuracyPostLocationFromBigBrother float64, AltitudePostLocationFromBigBrother float64, BearingPostLocationFromBigBrother float64, SpeedPostLocationFromBigBrother float64, BattlevelPostLocationFromBigBrother int32, TimePostLocationFromBigBrother string) (*pb.PostLocationFromBigBrotherRequest, error) {

	request := pb.PostLocationFromBigBrotherRequest{
		Latitude:  LatitudePostLocationFromBigBrother,
		Longitude: LongitudePostLocationFromBigBrother,
		Accuracy:  AccuracyPostLocationFromBigBrother,
		Altitude:  AltitudePostLocationFromBigBrother,
		Bearing:   BearingPostLocationFromBigBrother,
		Speed:     SpeedPostLocationFromBigBrother,
		Battlevel: BattlevelPostLocationFromBigBrother,
		Time:      TimePostLocationFromBigBrother,
	}
	return &request, nil
}

// GetLocation implements Service.
func GetLocation() (*pb.GetLocationRequest, error) {

	request := pb.GetLocationRequest{}
	return &request, nil
}
