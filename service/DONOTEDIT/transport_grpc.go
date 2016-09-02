package addsvc

// This file provides server-side bindings for the gRPC transport.
// It utilizes the transport/grpc.Server.

import (
	//stdopentracing "github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"

	//"github.com/go-kit/kit/log"
	//"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"

	// This Service
	pb "github.com/zaquestion/current/service/DONOTEDIT/pb"
)

// MakeGRPCServer makes a set of endpoints available as a gRPC AddServer.
func MakeGRPCServer(ctx context.Context, endpoints Endpoints /*, tracer stdopentracing.Tracer, logger log.Logger*/) pb.CurrentServiceServer {
	//options := []grpctransport.ServerOption{
	//grpctransport.ServerErrorLogger(logger),
	//}
	return &grpcServer{
		// currentservice

		postlocationfrombigbrother: grpctransport.NewServer(
			ctx,
			endpoints.PostLocationFromBigBrotherEndpoint,
			DecodeGRPCPostLocationFromBigBrotherRequest,
			EncodeGRPCPostLocationFromBigBrotherResponse,
			//append(options,grpctransport.ServerBefore(opentracing.FromGRPCRequest(tracer, "PostLocationFromBigBrother", logger)))...,
		),
		getlocation: grpctransport.NewServer(
			ctx,
			endpoints.GetLocationEndpoint,
			DecodeGRPCGetLocationRequest,
			EncodeGRPCGetLocationResponse,
			//append(options,grpctransport.ServerBefore(opentracing.FromGRPCRequest(tracer, "GetLocation", logger)))...,
		),
	}
}

type grpcServer struct {
	postlocationfrombigbrother grpctransport.Handler
	getlocation                grpctransport.Handler
}

// Methods

func (s *grpcServer) PostLocationFromBigBrother(ctx context.Context, req *pb.PostLocationFromBigBrotherRequest) (*pb.PostLocationFromBigBrotherReply, error) {
	_, rep, err := s.postlocationfrombigbrother.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.PostLocationFromBigBrotherReply), nil
}

func (s *grpcServer) GetLocation(ctx context.Context, req *pb.GetLocationRequest) (*pb.GetLocationReply, error) {
	_, rep, err := s.getlocation.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GetLocationReply), nil
}

// Server Decode

// DecodeGRPCPostLocationFromBigBrotherRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC postlocationfrombigbrother request to a user-domain postlocationfrombigbrother request. Primarily useful in a server.
func DecodeGRPCPostLocationFromBigBrotherRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.PostLocationFromBigBrotherRequest)
	return req, nil
}

// DecodeGRPCGetLocationRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC getlocation request to a user-domain getlocation request. Primarily useful in a server.
func DecodeGRPCGetLocationRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetLocationRequest)
	return req, nil
}

// Client Decode

// DecodeGRPCPostLocationFromBigBrotherResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC postlocationfrombigbrother reply to a user-domain postlocationfrombigbrother response. Primarily useful in a client.
func DecodeGRPCPostLocationFromBigBrotherResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.PostLocationFromBigBrotherReply)
	return reply, nil
}

// DecodeGRPCGetLocationResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC getlocation reply to a user-domain getlocation response. Primarily useful in a client.
func DecodeGRPCGetLocationResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.GetLocationReply)
	return reply, nil
}

// Server Encode

// EncodeGRPCPostLocationFromBigBrotherResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain postlocationfrombigbrother response to a gRPC postlocationfrombigbrother reply. Primarily useful in a server.
func EncodeGRPCPostLocationFromBigBrotherResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.PostLocationFromBigBrotherReply)
	return resp, nil
}

// EncodeGRPCGetLocationResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain getlocation response to a gRPC getlocation reply. Primarily useful in a server.
func EncodeGRPCGetLocationResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.GetLocationReply)
	return resp, nil
}

// Client Encode

// EncodeGRPCPostLocationFromBigBrotherRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain postlocationfrombigbrother request to a gRPC postlocationfrombigbrother request. Primarily useful in a client.
func EncodeGRPCPostLocationFromBigBrotherRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(pb.PostLocationFromBigBrotherRequest)
	return &req, nil
}

// EncodeGRPCGetLocationRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain getlocation request to a gRPC getlocation request. Primarily useful in a client.
func EncodeGRPCGetLocationRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(pb.GetLocationRequest)
	return &req, nil
}
