package middlewares

import (
	pb "github.com/zaquestion/current/current-service"
)

func WrapService(in pb.CurrentServer) pb.CurrentServer {
	return in
}
