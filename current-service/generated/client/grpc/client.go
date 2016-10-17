// Package grpc provides a gRPC client for the add service.
package grpc

import (
	//"time"

	//jujuratelimit "github.com/juju/ratelimit"
	//stdopentracing "github.com/opentracing/opentracing-go"
	//"github.com/sony/gobreaker"
	"google.golang.org/grpc"

	//"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	//"github.com/go-kit/kit/log"
	//"github.com/go-kit/kit/ratelimit"
	//"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"

	// This Service
	pb "github.com/zaquestion/current/current-service"
	svc "github.com/zaquestion/current/current-service/generated"
	handler "github.com/zaquestion/current/current-service/handlers/server"
)

// New returns an AddService backed by a gRPC client connection. It is the
// responsibility of the caller to dial, and later close, the connection.
func New(conn *grpc.ClientConn /*, tracer stdopentracing.Tracer, logger log.Logger*/) handler.Service {
	// We construct a single ratelimiter middleware, to limit the total outgoing
	// QPS from this client to all methods on the remote instance. We also
	// construct per-endpoint circuitbreaker middlewares to demonstrate how
	// that's done, although they could easily be combined into a single breaker
	// for the entire remote instance, too.

	//limiter := ratelimit.NewTokenBucketLimiter(jujuratelimit.NewBucketWithRate(100, 100))

	var postlocationfrombigbrotherEndpoint endpoint.Endpoint
	{
		postlocationfrombigbrotherEndpoint = grpctransport.NewClient(
			conn,
			"current.CurrentService",
			"PostLocationFromBigBrother",
			svc.EncodeGRPCPostLocationFromBigBrotherRequest,
			svc.DecodeGRPCPostLocationFromBigBrotherResponse,
			pb.PostLocationFromBigBrotherReply{},
			//grpctransport.ClientBefore(opentracing.FromGRPCRequest(tracer, "PostLocationFromBigBrother", logger)),
		).Endpoint()
		//postlocationfrombigbrotherEndpoint = opentracing.TraceClient(tracer, "PostLocationFromBigBrother")(postlocationfrombigbrotherEndpoint)
		//postlocationfrombigbrotherEndpoint = limiter(postlocationfrombigbrotherEndpoint)
		//postlocationfrombigbrotherEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		//Name:    "PostLocationFromBigBrother",
		//Timeout: 30 * time.Second,
		//}))(postlocationfrombigbrotherEndpoint)
	}

	var getlocationEndpoint endpoint.Endpoint
	{
		getlocationEndpoint = grpctransport.NewClient(
			conn,
			"current.CurrentService",
			"GetLocation",
			svc.EncodeGRPCGetLocationRequest,
			svc.DecodeGRPCGetLocationResponse,
			pb.GetLocationReply{},
			//grpctransport.ClientBefore(opentracing.FromGRPCRequest(tracer, "GetLocation", logger)),
		).Endpoint()
		//getlocationEndpoint = opentracing.TraceClient(tracer, "GetLocation")(getlocationEndpoint)
		//getlocationEndpoint = limiter(getlocationEndpoint)
		//getlocationEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		//Name:    "GetLocation",
		//Timeout: 30 * time.Second,
		//}))(getlocationEndpoint)
	}

	return svc.Endpoints{

		PostLocationFromBigBrotherEndpoint: postlocationfrombigbrotherEndpoint,
		GetLocationEndpoint:                getlocationEndpoint,
	}
}
