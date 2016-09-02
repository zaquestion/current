package main

import (
	"flag"
	"fmt"
	"os"
	//"strings"
	"time"

	//"github.com/lightstep/lightstep-tracer-go"
	//stdopentracing "github.com/opentracing/opentracing-go"
	//zipkin "github.com/openzipkin/zipkin-go-opentracing"
	//appdashot "github.com/sourcegraph/appdash/opentracing"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	//"sourcegraph.com/sourcegraph/appdash"

	//"github.com/go-kit/kit/log"

	// This Service
	grpcclient "github.com/zaquestion/current/service/DONOTEDIT/client/grpc"
	httpclient "github.com/zaquestion/current/service/DONOTEDIT/client/http"
	clientHandler "github.com/zaquestion/current/service/client"
	handler "github.com/zaquestion/current/service/server"
)

func main() {
	// The addcli presumes no service discovery system, and expects users to
	// provide the direct address of an addsvc. This presumption is reflected in
	// the addcli binary and the the client packages: the -transport.addr flags
	// and various client constructors both expect host:port strings. For an
	// example service with a client built on top of a service discovery system,
	// see profilesvc.

	var (
		httpAddr = flag.String("http.addr", "", "HTTP address of addsvc")
		grpcAddr = flag.String("grpc.addr", "", "gRPC (HTTP) address of addsvc")
		//zipkinAddr     = flag.String("zipkin.addr", "", "Enable Zipkin tracing via a Kafka Collector host:port")
		//appdashAddr    = flag.String("appdash.addr", "", "Enable Appdash tracing via an Appdash server host:port")
		//lightstepToken = flag.String("lightstep.token", "", "Enable LightStep tracing via a LightStep access token")
		method = flag.String("method", "postlocationfrombigbrother", "postlocationfrombigbrother,getlocation")
	)

	var (
		LatitudePostLocationFromBigBrother  = flag.Float64("postlocationfrombigbrother.latitude", 0.0, "")
		LongitudePostLocationFromBigBrother = flag.Float64("postlocationfrombigbrother.longitude", 0.0, "")
		AccuracyPostLocationFromBigBrother  = flag.Float64("postlocationfrombigbrother.accuracy", 0.0, "")
		AltitudePostLocationFromBigBrother  = flag.Float64("postlocationfrombigbrother.altitude", 0.0, "")
		BearingPostLocationFromBigBrother   = flag.Float64("postlocationfrombigbrother.bearing", 0.0, "")
		SpeedPostLocationFromBigBrother     = flag.Float64("postlocationfrombigbrother.speed", 0.0, "")
	)
	flag.Parse()

	// This is a demonstration client, which supports multiple tracers.
	// Your clients will probably just use one tracer.
	//var tracer stdopentracing.Tracer
	//{
	//if *zipkinAddr != "" {
	//collector, err := zipkin.NewKafkaCollector(
	//strings.Split(*zipkinAddr, ","),
	//zipkin.KafkaLogger(log.NewNopLogger()),
	//)
	//if err != nil {
	//fmt.Fprintf(os.Stderr, "%v\n", err)
	//os.Exit(1)
	//}
	//tracer, err = zipkin.NewTracer(
	//zipkin.NewRecorder(collector, false, "localhost:8000", "addcli"),
	//)
	//if err != nil {
	//fmt.Fprintf(os.Stderr, "%v\n", err)
	//os.Exit(1)
	//}
	//} else if *appdashAddr != "" {
	//tracer = appdashot.NewTracer(appdash.NewRemoteCollector(*appdashAddr))
	//} else if *lightstepToken != "" {
	//tracer = lightstep.NewTracer(lightstep.Options{
	//AccessToken: *lightstepToken,
	//})
	//defer lightstep.FlushLightStepTracer(tracer)
	//} else {
	//tracer = stdopentracing.GlobalTracer() // no-op
	//}
	//}

	// This is a demonstration client, which supports multiple transports.
	// Your clients will probably just define and stick with 1 transport.

	var (
		service handler.Service
		err     error
	)
	if *httpAddr != "" {
		//service, err = httpclient.New(*httpAddr, tracer, log.NewNopLogger())
		service, err = httpclient.New(*httpAddr)
	} else if *grpcAddr != "" {
		conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v", err)
			os.Exit(1)
		}
		defer conn.Close()
		service = grpcclient.New(conn /*, tracer, log.NewNopLogger()*/)
	} else {
		fmt.Fprintf(os.Stderr, "error: no remote address specified\n")
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	switch *method {

	case "postlocationfrombigbrother":
		request, _ := clientHandler.PostLocationFromBigBrother(*LatitudePostLocationFromBigBrother, *LongitudePostLocationFromBigBrother, *AccuracyPostLocationFromBigBrother, *AltitudePostLocationFromBigBrother, *BearingPostLocationFromBigBrother, *SpeedPostLocationFromBigBrother)
		v, err := service.PostLocationFromBigBrother(context.Background(), request)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Client Requested with:")
		fmt.Println(*LatitudePostLocationFromBigBrother, *LongitudePostLocationFromBigBrother, *AccuracyPostLocationFromBigBrother, *AltitudePostLocationFromBigBrother, *BearingPostLocationFromBigBrother, *SpeedPostLocationFromBigBrother)
		fmt.Println("Server Responded with:")
		fmt.Println(v)

	case "getlocation":
		request, _ := clientHandler.GetLocation()
		v, err := service.GetLocation(context.Background(), request)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Client Requested with:")
		fmt.Println()
		fmt.Println("Server Responded with:")
		fmt.Println(v)

	default:
		fmt.Fprintf(os.Stderr, "error: invalid method %q\n", method)
		os.Exit(1)
	}
}
