package addsvc

// This file provides server-side bindings for the HTTP transport.
// It utilizes the transport/http.Server.

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	//stdopentracing "github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"

	"github.com/go-kit/kit/log"
	//"github.com/go-kit/kit/endpoint"
	//"github.com/go-kit/kit/tracing/opentracing"
	httptransport "github.com/go-kit/kit/transport/http"

	// This service
	pb "github.com/zaquestion/current/service/DONOTEDIT/pb"
)

var (
	_ = fmt.Sprint
	_ = bytes.Compare
	_ = strconv.Atoi
	_ = httptransport.NewServer
	_ = ioutil.NopCloser
)

// MakeHTTPHandler returns a handler that makes a set of endpoints available
// on predefined paths.
func MakeHTTPHandler(ctx context.Context, endpoints Endpoints, logger log.Logger) http.Handler {
	//func MakeHTTPHandler(ctx context.Context, endpoints Endpoints, /*tracer stdopentracing.Tracer,*/ logger log.Logger) http.Handler {
	/*options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorLogger(logger),
	}*/
	m := http.NewServeMux()

	m.Handle("/location/bigbrother", httptransport.NewServer(
		ctx,
		endpoints.PostLocationFromBigBrotherEndpoint,
		HttpDecodeLogger(DecodeHTTPPostLocationFromBigBrotherZeroRequest, logger),
		EncodeHTTPGenericResponse,
	))

	m.Handle("/location", httptransport.NewServer(
		ctx,
		endpoints.GetLocationEndpoint,
		HttpDecodeLogger(DecodeHTTPGetLocationZeroRequest, logger),
		EncodeHTTPGenericResponse,
	))
	return m
}

func HttpDecodeLogger(next httptransport.DecodeRequestFunc, logger log.Logger) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		logger.Log("method", r.Method, "url", r.URL.String())
		return next(ctx, r)
	}
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	code := http.StatusInternalServerError
	msg := err.Error()

	/*if e, ok := err.(httptransport.Error); ok {
		msg = e.Err.Error()
		switch e.Domain {
		case httptransport.DomainDecode:
			code = http.StatusBadRequest

		case httptransport.DomainDo:
			switch e.Err {
			case ErrTwoZeroes, ErrMaxSizeExceeded, ErrIntOverflow:
				code = http.StatusBadRequest
			}
		}
	}*/

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errorWrapper{Error: msg})
}

func errorDecoder(r *http.Response) error {
	var w errorWrapper
	if err := json.NewDecoder(r.Body).Decode(&w); err != nil {
		return err
	}
	return errors.New(w.Error)
}

type errorWrapper struct {
	Error string `json:"error"`
}

// Server Decode

// DecodeHTTPPostLocationFromBigBrotherZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded postlocationfrombigbrother request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPPostLocationFromBigBrotherZeroRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req pb.PostLocationFromBigBrotherRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	pathParams, err := PathParams(r.URL.Path, "/location/bigbrother")
	_ = pathParams
	// TODO: Better error handling
	if err != nil {
		fmt.Printf("Error while reading path params: %v\n", err)
		return nil, err
	}
	queryParams, err := QueryParams(r.URL.Query())
	_ = queryParams
	// TODO: Better error handling
	if err != nil {
		fmt.Printf("Error while reading query params: %v\n", err)
		return nil, err
	}

	LatitudePostLocationFromBigBrotherStr := queryParams["latitude"]
	LatitudePostLocationFromBigBrother, err := strconv.ParseFloat(LatitudePostLocationFromBigBrotherStr, 64)
	// TODO: Better error handling
	if err != nil {
		fmt.Printf("Error while extracting LatitudePostLocationFromBigBrother from query: %v\n", err)
		fmt.Printf("queryParams: %v\n", queryParams)
		return nil, err
	}
	req.Latitude = LatitudePostLocationFromBigBrother

	LongitudePostLocationFromBigBrotherStr := queryParams["longitude"]
	LongitudePostLocationFromBigBrother, err := strconv.ParseFloat(LongitudePostLocationFromBigBrotherStr, 64)
	// TODO: Better error handling
	if err != nil {
		fmt.Printf("Error while extracting LongitudePostLocationFromBigBrother from query: %v\n", err)
		fmt.Printf("queryParams: %v\n", queryParams)
		return nil, err
	}
	req.Longitude = LongitudePostLocationFromBigBrother

	AccuracyPostLocationFromBigBrotherStr := queryParams["accuracy"]
	AccuracyPostLocationFromBigBrother, err := strconv.ParseFloat(AccuracyPostLocationFromBigBrotherStr, 64)
	// TODO: Better error handling
	if err != nil {
		fmt.Printf("Error while extracting AccuracyPostLocationFromBigBrother from query: %v\n", err)
		fmt.Printf("queryParams: %v\n", queryParams)
		return nil, err
	}
	req.Accuracy = AccuracyPostLocationFromBigBrother

	AltitudePostLocationFromBigBrotherStr := queryParams["altitude"]
	AltitudePostLocationFromBigBrother, err := strconv.ParseFloat(AltitudePostLocationFromBigBrotherStr, 64)
	// TODO: Better error handling
	if err != nil {
		fmt.Printf("Error while extracting AltitudePostLocationFromBigBrother from query: %v\n", err)
		fmt.Printf("queryParams: %v\n", queryParams)
		return nil, err
	}
	req.Altitude = AltitudePostLocationFromBigBrother

	BearingPostLocationFromBigBrotherStr := queryParams["bearing"]
	BearingPostLocationFromBigBrother, err := strconv.ParseFloat(BearingPostLocationFromBigBrotherStr, 64)
	// TODO: Better error handling
	if err != nil {
		fmt.Printf("Error while extracting BearingPostLocationFromBigBrother from query: %v\n", err)
		fmt.Printf("queryParams: %v\n", queryParams)
		return nil, err
	}
	req.Bearing = BearingPostLocationFromBigBrother

	SpeedPostLocationFromBigBrotherStr := queryParams["speed"]
	SpeedPostLocationFromBigBrother, err := strconv.ParseFloat(SpeedPostLocationFromBigBrotherStr, 64)
	// TODO: Better error handling
	if err != nil {
		fmt.Printf("Error while extracting SpeedPostLocationFromBigBrother from query: %v\n", err)
		fmt.Printf("queryParams: %v\n", queryParams)
		return nil, err
	}
	req.Speed = SpeedPostLocationFromBigBrother

	return &req, err
}

// DecodeHTTPGetLocationZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded getlocation request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPGetLocationZeroRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req pb.GetLocationRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	pathParams, err := PathParams(r.URL.Path, "/location")
	_ = pathParams
	// TODO: Better error handling
	if err != nil {
		fmt.Printf("Error while reading path params: %v\n", err)
		return nil, err
	}
	queryParams, err := QueryParams(r.URL.Query())
	_ = queryParams
	// TODO: Better error handling
	if err != nil {
		fmt.Printf("Error while reading query params: %v\n", err)
		return nil, err
	}

	return &req, err
}

// Client Decode

// DecodeHTTPPostLocationFromBigBrother is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded PostLocationFromBigBrotherReply response from the HTTP response body.
// If the response has a non-200 status code, we will interpret that as an
// error and attempt to decode the specific error message from the response
// body. Primarily useful in a client.
func DecodeHTTPPostLocationFromBigBrotherResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errorDecoder(r)
	}
	var resp pb.PostLocationFromBigBrotherReply
	err := json.NewDecoder(r.Body).Decode(&resp)
	return &resp, err
}

// DecodeHTTPGetLocation is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded GetLocationReply response from the HTTP response body.
// If the response has a non-200 status code, we will interpret that as an
// error and attempt to decode the specific error message from the response
// body. Primarily useful in a client.
func DecodeHTTPGetLocationResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errorDecoder(r)
	}
	var resp pb.GetLocationReply
	err := json.NewDecoder(r.Body).Decode(&resp)
	return &resp, err
}

// Client Encode

// EncodeHTTPPostLocationFromBigBrotherZeroRequest is a transport/http.EncodeRequestFunc
// that encodes a postlocationfrombigbrother request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPPostLocationFromBigBrotherZeroRequest(_ context.Context, r *http.Request, request interface{}) error {
	fmt.Printf("Encoding request %v\n", request)
	req := request.(pb.PostLocationFromBigBrotherRequest)
	_ = req

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"location",
		"bigbrother",
	}, "/")
	//r.URL.Scheme,
	//r.URL.Host,
	u, err := url.Parse(path)
	if err != nil {
		return err
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	values.Add("latitude", fmt.Sprint(req.Latitude))
	values.Add("longitude", fmt.Sprint(req.Longitude))
	values.Add("accuracy", fmt.Sprint(req.Accuracy))
	values.Add("altitude", fmt.Sprint(req.Altitude))
	values.Add("bearing", fmt.Sprint(req.Bearing))
	values.Add("speed", fmt.Sprint(req.Speed))

	r.URL.RawQuery = values.Encode()

	// Set the body parameters
	var buf bytes.Buffer
	toRet := map[string]interface{}{}
	if err := json.NewEncoder(&buf).Encode(toRet); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	fmt.Printf("URL: %v\n", r.URL)
	return nil
}

// EncodeHTTPGetLocationZeroRequest is a transport/http.EncodeRequestFunc
// that encodes a getlocation request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPGetLocationZeroRequest(_ context.Context, r *http.Request, request interface{}) error {
	fmt.Printf("Encoding request %v\n", request)
	req := request.(pb.GetLocationRequest)
	_ = req

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"location",
	}, "/")
	//r.URL.Scheme,
	//r.URL.Host,
	u, err := url.Parse(path)
	if err != nil {
		return err
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()

	r.URL.RawQuery = values.Encode()

	// Set the body parameters
	var buf bytes.Buffer
	toRet := map[string]interface{}{}
	if err := json.NewEncoder(&buf).Encode(toRet); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	fmt.Printf("URL: %v\n", r.URL)
	return nil
}

// EncodeHTTPGenericResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeHTTPGenericResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func PathParams(url string, urlTmpl string) (map[string]string, error) {
	rv := map[string]string{}
	pmp := BuildParamMap(urlTmpl)

	parts := strings.Split(url, "/")
	for k, v := range pmp {
		rv[k] = parts[v]
	}

	return rv, nil
}
func BuildParamMap(urlTmpl string) map[string]int {
	rv := map[string]int{}

	parts := strings.Split(urlTmpl, "/")
	for idx, part := range parts {
		if strings.ContainsAny(part, "{}") {
			param := RemoveBraces(part)
			rv[param] = idx
		}
	}
	return rv
}
func RemoveBraces(val string) string {
	val = strings.Replace(val, "{", "", -1)
	val = strings.Replace(val, "}", "", -1)
	return val
}
func QueryParams(vals url.Values) (map[string]string, error) {

	rv := map[string]string{}
	for k, v := range vals {
		rv[k] = v[0]
	}
	return rv, nil
}
