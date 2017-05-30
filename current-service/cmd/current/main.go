// Code generated by truss.
// Rerunning truss will overwrite this file.
// DO NOT EDIT!
// Version: a41ee29fb6
// Version Date: Fri May 26 18:19:22 UTC 2017

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/pkg/errors"

	// This Service
	pb "github.com/zaquestion/current/current-service"
	"github.com/zaquestion/current/current-service/svc/client/cli/handlers"
	grpcclient "github.com/zaquestion/current/current-service/svc/client/grpc"
	httpclient "github.com/zaquestion/current/current-service/svc/client/http"
)

var (
	_ = strconv.ParseInt
	_ = strings.Split
	_ = json.Compact
	_ = errors.Wrapf
	_ = pb.RegisterCurrentServer
)

func main() {
	os.Exit(submain())
}

type headerSeries []string

func (h *headerSeries) Set(val string) error {
	const requiredParts int = 2
	parts := strings.SplitN(val, ":", requiredParts)
	if len(parts) != requiredParts {
		return fmt.Errorf("value %q cannot be split in two; must contain at least one ':' character", val)
	}
	parts[1] = strings.TrimSpace(parts[1])
	*h = append(*h, parts...)
	return nil
}

func (h *headerSeries) String() string {
	return fmt.Sprintf("%v", []string(*h))
}

// submain exists to act as the functional main, but will return exit codes to
// the actual main instead of calling os.Exit directly. This is done to allow
// the defered functions to be called, since if os.Exit where called directly
// from this function, none of the defered functions set up by this function
// would be called.
func submain() int {
	var headers headerSeries
	flag.Var(&headers, "header", "Header(s) to be sent in the transport (follows cURL style)")
	var (
		httpAddr = flag.String("http.addr", "", "HTTP address of addsvc")
		grpcAddr = flag.String("grpc.addr", ":5040", "gRPC (HTTP) address of addsvc")
	)

	// The addcli presumes no service discovery system, and expects users to
	// provide the direct address of an service. This presumption is reflected in
	// the cli binary and the the client packages: the -transport.addr flags
	// and various client constructors both expect host:port strings.

	fsGetLocation := flag.NewFlagSet("getlocation", flag.ExitOnError)

	fsPostLocationBigBrother := flag.NewFlagSet("postlocationbigbrother", flag.ExitOnError)

	fsPostLocationTasker := flag.NewFlagSet("postlocationtasker", flag.ExitOnError)

	var (
		flagLatitudePostLocationBigBrother  = fsPostLocationBigBrother.Float64("latitude", 0.0, "")
		flagLongitudePostLocationBigBrother = fsPostLocationBigBrother.Float64("longitude", 0.0, "")
		flagAccuracyPostLocationBigBrother  = fsPostLocationBigBrother.Float64("accuracy", 0.0, "")
		flagAltitudePostLocationBigBrother  = fsPostLocationBigBrother.Float64("altitude", 0.0, "")
		flagBearingPostLocationBigBrother   = fsPostLocationBigBrother.Float64("bearing", 0.0, "")
		flagSpeedPostLocationBigBrother     = fsPostLocationBigBrother.Float64("speed", 0.0, "")
		flagBattlevelPostLocationBigBrother = fsPostLocationBigBrother.Int("battlevel", 0, "")
		flagTimePostLocationBigBrother      = fsPostLocationBigBrother.String("time", "", "")
		flagLocationPostLocationTasker      = fsPostLocationTasker.String("location", "", "")
		flagSpeedPostLocationTasker         = fsPostLocationTasker.Float64("speed", 0.0, "")
		flagBatteryPostLocationTasker       = fsPostLocationTasker.Int("battery", 0, "")
		flagChargingPostLocationTasker      = fsPostLocationTasker.Bool("charging", false, "")
		flagDateTimePostLocationTasker      = fsPostLocationTasker.String("datetime", "", "")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Subcommands:\n")
		fmt.Fprintf(os.Stderr, "  %s\n", "getlocation")
		fmt.Fprintf(os.Stderr, "  %s\n", "postlocationbigbrother")
		fmt.Fprintf(os.Stderr, "  %s\n", "postlocationtasker")
	}
	if len(os.Args) < 2 {
		flag.Usage()
		return 1
	}

	flag.Parse()

	var (
		service pb.CurrentServer
		err     error
	)

	if *httpAddr != "" {
		service, err = httpclient.New(*httpAddr, httpclient.CtxValuesToSend(headers...))
	} else if *grpcAddr != "" {
		conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while dialing grpc connection: %v", err)
			return 1
		}
		defer conn.Close()
		service, err = grpcclient.New(conn, grpcclient.CtxValuesToSend(headers...))
	} else {
		fmt.Fprintf(os.Stderr, "error: no remote address specified\n")
		return 1
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	if len(flag.Args()) < 1 {
		fmt.Printf("No 'method' subcommand provided; exiting.")
		flag.Usage()
		return 1
	}

	ctx := context.Background()
	for i := 0; i < len(headers); i += 2 {
		ctx = context.WithValue(ctx, headers[i], headers[i+1])
	}

	switch flag.Args()[0] {

	case "getlocation":
		fsGetLocation.Parse(flag.Args()[1:])

		request, err := handlers.GetLocation()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error calling handlers.GetLocation: %v\n", err)
			return 1
		}

		v, err := service.GetLocation(ctx, request)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error calling service.GetLocation: %v\n", err)
			return 1
		}
		fmt.Println("Client Requested with:")
		fmt.Println()
		fmt.Println("Server Responded with:")
		fmt.Println(v)

	case "postlocationbigbrother":
		fsPostLocationBigBrother.Parse(flag.Args()[1:])

		LatitudePostLocationBigBrother := *flagLatitudePostLocationBigBrother
		LongitudePostLocationBigBrother := *flagLongitudePostLocationBigBrother
		AccuracyPostLocationBigBrother := *flagAccuracyPostLocationBigBrother
		AltitudePostLocationBigBrother := *flagAltitudePostLocationBigBrother
		BearingPostLocationBigBrother := *flagBearingPostLocationBigBrother
		SpeedPostLocationBigBrother := *flagSpeedPostLocationBigBrother
		BattlevelPostLocationBigBrother := int32(*flagBattlevelPostLocationBigBrother)
		TimePostLocationBigBrother := *flagTimePostLocationBigBrother

		request, err := handlers.PostLocationBigBrother(LatitudePostLocationBigBrother, LongitudePostLocationBigBrother, AccuracyPostLocationBigBrother, AltitudePostLocationBigBrother, BearingPostLocationBigBrother, SpeedPostLocationBigBrother, BattlevelPostLocationBigBrother, TimePostLocationBigBrother)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error calling handlers.PostLocationBigBrother: %v\n", err)
			return 1
		}

		v, err := service.PostLocationBigBrother(ctx, request)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error calling service.PostLocationBigBrother: %v\n", err)
			return 1
		}
		fmt.Println("Client Requested with:")
		fmt.Println(LatitudePostLocationBigBrother, LongitudePostLocationBigBrother, AccuracyPostLocationBigBrother, AltitudePostLocationBigBrother, BearingPostLocationBigBrother, SpeedPostLocationBigBrother, BattlevelPostLocationBigBrother, TimePostLocationBigBrother)
		fmt.Println("Server Responded with:")
		fmt.Println(v)

	case "postlocationtasker":
		fsPostLocationTasker.Parse(flag.Args()[1:])

		var LocationPostLocationTasker []float64
		if flagLocationPostLocationTasker != nil && len(*flagLocationPostLocationTasker) > 0 {
			err = json.Unmarshal([]byte(*flagLocationPostLocationTasker), &LocationPostLocationTasker)
			if err != nil {
				panic(errors.Wrapf(err, "unmarshalling LocationPostLocationTasker from %v:", flagLocationPostLocationTasker))
			}
		}

		SpeedPostLocationTasker := *flagSpeedPostLocationTasker
		BatteryPostLocationTasker := int32(*flagBatteryPostLocationTasker)
		ChargingPostLocationTasker := *flagChargingPostLocationTasker
		DateTimePostLocationTasker := *flagDateTimePostLocationTasker

		request, err := handlers.PostLocationTasker(LocationPostLocationTasker, SpeedPostLocationTasker, BatteryPostLocationTasker, ChargingPostLocationTasker, DateTimePostLocationTasker)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error calling handlers.PostLocationTasker: %v\n", err)
			return 1
		}

		v, err := service.PostLocationTasker(ctx, request)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error calling service.PostLocationTasker: %v\n", err)
			return 1
		}
		fmt.Println("Client Requested with:")
		fmt.Println(LocationPostLocationTasker, SpeedPostLocationTasker, BatteryPostLocationTasker, ChargingPostLocationTasker, DateTimePostLocationTasker)
		fmt.Println("Server Responded with:")
		fmt.Println(v)

	default:
		flag.Usage()
		return 1
	}

	return 0
}
