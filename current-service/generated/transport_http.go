package svc

// This file provides server-side bindings for the HTTP transport.
// It utilizes the transport/http.Server.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	//stdopentracing "github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"

	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	//"github.com/go-kit/kit/endpoint"
	//"github.com/go-kit/kit/tracing/opentracing"
	httptransport "github.com/go-kit/kit/transport/http"

	// This service
	pb "github.com/zaquestion/current/current-service"
)

var (
	_ = fmt.Sprint
	_ = bytes.Compare
	_ = strconv.Atoi
	_ = httptransport.NewServer
	_ = ioutil.NopCloser
	_ = pb.RegisterCurrentServiceServer
	_ = io.Copy
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
		rv, err := next(ctx, r)
		if err != nil {
			logger.Log("method", r.Method, "url", r.URL.String(), "Error", err)
		}
		return rv, err
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
	// err = io.EOF if r.Body was empty
	if err != nil && err != io.EOF {
		return nil, errors.Wrap(err, "decoding body of http request")
	}

	pathParams, err := PathParams(r.URL.Path, "/location/bigbrother")
	_ = pathParams
	if err != nil {
		fmt.Printf("Error while reading path params: %v\n", err)
		return nil, errors.Wrap(err, "couldn't unmarshal path parameters")
	}
	queryParams, err := QueryParams(r.URL.Query())
	_ = queryParams
	if err != nil {
		fmt.Printf("Error while reading query params: %v\n", err)
		return nil, errors.Wrapf(err, "Error while reading query params: %v", r.URL.Query())
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

	BattlevelPostLocationFromBigBrotherStr := queryParams["battlevel"]
	BattlevelPostLocationFromBigBrother, err := strconv.ParseInt(BattlevelPostLocationFromBigBrotherStr, 10, 32)
	// TODO: Better error handling
	if err != nil {
		fmt.Printf("Error while extracting BattlevelPostLocationFromBigBrother from query: %v\n", err)
		fmt.Printf("queryParams: %v\n", queryParams)
		return nil, err
	}
	req.Battlevel = int32(BattlevelPostLocationFromBigBrother)

	TimePostLocationFromBigBrotherStr := queryParams["time"]
	TimePostLocationFromBigBrother := TimePostLocationFromBigBrotherStr
	// TODO: Better error handling
	if err != nil {
		fmt.Printf("Error while extracting TimePostLocationFromBigBrother from query: %v\n", err)
		fmt.Printf("queryParams: %v\n", queryParams)
		return nil, err
	}
	req.Time = TimePostLocationFromBigBrother

	return &req, err
}

// DecodeHTTPGetLocationZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded getlocation request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPGetLocationZeroRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req pb.GetLocationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	// err = io.EOF if r.Body was empty
	if err != nil && err != io.EOF {
		return nil, errors.Wrap(err, "decoding body of http request")
	}

	pathParams, err := PathParams(r.URL.Path, "/location")
	_ = pathParams
	if err != nil {
		fmt.Printf("Error while reading path params: %v\n", err)
		return nil, errors.Wrap(err, "couldn't unmarshal path parameters")
	}
	queryParams, err := QueryParams(r.URL.Query())
	_ = queryParams
	if err != nil {
		fmt.Printf("Error while reading query params: %v\n", err)
		return nil, errors.Wrapf(err, "Error while reading query params: %v", r.URL.Query())
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
	req := request.(*pb.PostLocationFromBigBrotherRequest)
	_ = req

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"location",
		"bigbrother",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	values.Add("latitude", fmt.Sprint(req.Latitude))

	values.Add("longitude", fmt.Sprint(req.Longitude))

	values.Add("accuracy", fmt.Sprint(req.Accuracy))

	values.Add("altitude", fmt.Sprint(req.Altitude))

	values.Add("bearing", fmt.Sprint(req.Bearing))

	values.Add("speed", fmt.Sprint(req.Speed))

	values.Add("battlevel", fmt.Sprint(req.Battlevel))

	values.Add("time", fmt.Sprint(req.Time))

	r.URL.RawQuery = values.Encode()

	// Set the body parameters
	var buf bytes.Buffer
	toRet := map[string]interface{}{}
	if err := json.NewEncoder(&buf).Encode(toRet); err != nil {
		return errors.Wrapf(err, "couldn't encode body as json %v", toRet)
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
	req := request.(*pb.GetLocationRequest)
	_ = req

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"location",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	r.URL.RawQuery = values.Encode()

	// Set the body parameters
	var buf bytes.Buffer
	toRet := map[string]interface{}{}
	if err := json.NewEncoder(&buf).Encode(toRet); err != nil {
		return errors.Wrapf(err, "couldn't encode body as json %v", toRet)
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

// PathParams takes a url and a gRPC-annotation style url template, and
// returns a map of the named parameters in the template and their values in
// the given url.
//
// PathParams does not support the entirety of the URL template syntax defined
// in third_party/googleapis/google/api/httprule.proto. Only a small subset of
// the functionality defined there is implemented here.
func PathParams(url string, urlTmpl string) (map[string]string, error) {
	rv := map[string]string{}
	pmp := BuildParamMap(urlTmpl)

	parts := strings.Split(url, "/")
	for k, v := range pmp {
		rv[k] = parts[v]
	}

	return rv, nil
}

// BuildParamMap takes a string representing a url template and returns a map
// indicating the location of each parameter within that url, where the
// location is the index as if in a slash-separated sequence of path
// components. For example, given the url template:
//
//     "/v1/{a}/{b}"
//
// The returned param map would look like:
//
//     map[string]int {
//         "a": 2,
//         "b": 3,
//     }
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

// RemoveBraces replace all curly braces in the provided string, opening and
// closing, with empty strings.
func RemoveBraces(val string) string {
	val = strings.Replace(val, "{", "", -1)
	val = strings.Replace(val, "}", "", -1)
	return val
}

// QueryParams takes query parameters in the form of url.Values, and returns a
// bare map of the string representation of each key to the string
// representation for each value. The representations of repeated query
// parameters is undefined.
func QueryParams(vals url.Values) (map[string]string, error) {

	rv := map[string]string{}
	for k, v := range vals {
		rv[k] = v[0]
	}
	return rv, nil
}
