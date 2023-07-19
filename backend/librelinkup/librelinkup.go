package librelinkup

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	libreLinkUpVersion = "4.7.0"
	libreLinkUpProduct = "llu.ios"
)

var (
	ErrLoginFailed         = errors.New("could not login to LibreLinkUp, make sure the username and password is correct")
	ErrWrongRegionEndpoint = errors.New("user called wrong region")
	linkupHeaders          = map[string][]string{
		"User-Agent":      {"LibreLink"},
		"Content-Type":    {"application/json"},
		"Version":         {libreLinkUpVersion},
		"Product":         {libreLinkUpProduct},
		"Accept-Encoding": {"gzip, deflate, br"},
		"Connection":      {"keep-alive"},
		"Pragma":          {"no-cache"},
		"Cache-Control":   {"no-cache"},
	}
)

type Ticket struct {
	Token    string `json:"token"`
	Expires  int64  `json:"expires"`
	Duration int64  `json:"duration"`
	endpoint Endpoint
}

type LibreLinkUpResponse struct {
	Status int `json:"status"`
	Error  struct {
		Message string `json:"message"`
	} `json:"error"`
}

type LoginResponse struct {
	LibreLinkUpResponse
	Data struct {
		User struct {
			ID string `json:"id"`
		} `json:"user"`
		AuthTicket Ticket `json:"authTicket"`
		Redirect   bool   `json:"redirect"`
		Region     string `json:"region"`
	} `json:"data"`
}

type ConnectionsResponse struct {
	LibreLinkUpResponse
	Data   []Connection `json:"data"`
	Ticket Ticket       `json:"ticket"`
}

type GraphResponse struct {
	LibreLinkUpResponse
	Data struct {
		Connection    Connection           `json:"connection"`
		ActiveSensors []ActiveSensor       `json:"activeSensors"`
		GraphData     []GlucoseMeasurement `json:"graphData"`
	} `json:"data"`
	AuthTicket Ticket `json:"ticket"`
}

type Connection struct {
	ID                 string             `json:"id"`
	PatientID          string             `json:"patientId"`
	FirstName          string             `json:"firstName"`
	LastName           string             `json:"lastName"`
	TargetLow          int                `json:"targetLow"`
	TargetHigh         int                `json:"targetHigh"`
	GlucoseMeasurement GlucoseMeasurement `json:"glucoseMeasurement"`
}

type GlucoseMeasurement struct {
	FactoryTimestamp string  `json:"FactoryTimestamp"`
	Timestamp        string  `json:"Timestamp"`
	Type             int     `json:"Type"`
	ValueInMgPerDl   int     `json:"ValueInMgPerDl"`
	MeasurementColor int     `json:"MeasurementColor"`
	GlucoseUnits     int     `json:"GlucoseUnits"`
	Value            float64 `json:"Value"`
	IsHigh           bool    `json:"isHigh"`
	IsLow            bool    `json:"isLow"`
}

type ActiveSensor struct {
	DeviceId       string    `json:"deviceId"`
	SerialNumber   string    `json:"sn"`
	ActivationTime time.Time `json:"a"`
}

func ToTime(libreTimestamp string) (time.Time, error) {
	return time.Parse("1/02/2006 3:04:05 PM", libreTimestamp)
}

func callLogin(email, password string, endpoint Endpoint) (LoginResponse, error) {
	type Credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	byteCredentials, err := json.Marshal(Credentials{Email: email, Password: password})
	if err != nil {
		return LoginResponse{}, fmt.Errorf("error marshaling credentials to JSON: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint.LoginURL(), bytes.NewBuffer(byteCredentials))
	if err != nil {
		return LoginResponse{}, fmt.Errorf("error creating request to '%s': %w", endpoint.LoginURL(), err)
	}

	var lr LoginResponse
	if err := doRequest(req, &Ticket{}, &lr); err != nil {
		return LoginResponse{}, err
	}
	if lr.Status != 0 {
		if lr.Status == 2 {
			return LoginResponse{}, ErrLoginFailed
		}
		return LoginResponse{}, fmt.Errorf("error during login, status code %v, message: %s", lr.Status, lr.Error.Message)
	}
	return lr, nil
}

func Login(email, password string, endpoint Endpoint) (*Ticket, error) {
	resp, err := callLogin(email, password, endpoint)
	if err != nil {
		return &Ticket{}, err
	}

	if resp.Data.Redirect {
		return &Ticket{}, ErrWrongRegionEndpoint
	}

	t := resp.Data.AuthTicket
	t.endpoint = endpoint
	return &t, nil
}

func FindEndpoint(email, password string) (Endpoint, error) {
	resp, err := callLogin(email, password, EndpointDefault)
	if err != nil {
		return EndpointDefault, err
	}

	if resp.Data.Redirect {
		e, found := EndpointByRegion(resp.Data.Region)
		if !found {
			return Endpoint{}, fmt.Errorf("endpoint with region '%s' could not be found", resp.Data.Region)
		}
		return e, nil
	}

	return EndpointDefault, nil
}

func (ticket *Ticket) Connections() ([]Connection, error) {
	req, err := http.NewRequest(http.MethodGet, ticket.endpoint.ConnectionsURL(), nil)
	if err != nil {
		return []Connection{}, err
	}

	var cr ConnectionsResponse
	if err := doRequest(req, ticket, &cr); err != nil {
		return []Connection{}, err
	}
	if cr.Status != 0 {
		return []Connection{}, fmt.Errorf("error occured, status code %v, message: %s", cr.Status, cr.Error.Message)
	}

	ticket.Token = cr.Ticket.Token
	ticket.Expires = cr.Ticket.Expires
	ticket.Duration = cr.Ticket.Duration

	return cr.Data, nil
}

func (ticket *Ticket) Graph(patientID string) (Connection, []GlucoseMeasurement, error) {
	req, err := http.NewRequest(http.MethodGet, ticket.endpoint.GraphURL(patientID), nil)
	if err != nil {
		return Connection{}, []GlucoseMeasurement{}, err
	}

	var gr GraphResponse
	if err := doRequest(req, ticket, &gr); err != nil {
		return Connection{}, []GlucoseMeasurement{}, err
	}
	if gr.Status != 0 {
		return Connection{}, []GlucoseMeasurement{}, fmt.Errorf("error occured, status code %v, message: %s", gr.Status, gr.Error.Message)
	}

	ticket.Token = gr.AuthTicket.Token
	ticket.Expires = gr.AuthTicket.Expires
	ticket.Duration = gr.AuthTicket.Duration

	return gr.Data.Connection, gr.Data.GraphData, nil
}

func doRequest[T any](req *http.Request, ticket *Ticket, response T) error {
	h := linkupHeaders
	if ticket.Token != "" {
		http.Header(h).Set("Authorization", "Bearer "+ticket.Token)
	}
	req.Header = h

	if isDebug() {
		printRequest(req)
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error executing request to '%s': %w", req.URL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server responded with status %v: %s", resp.StatusCode, resp.Status)
	}
	result, err := bodyToString(resp)
	if err != nil {
		return fmt.Errorf("error while reading response: %w", err)
	}

	if err := json.Unmarshal([]byte(result), response); err != nil {
		return fmt.Errorf("error while unmarshaling response from JSON: %w", err)
	}

	if isDebug() {
		printResponse(resp, result)
	}

	return nil
}

func bodyToString(resp *http.Response) (string, error) {
	var result []byte
	var err error
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return "", err
		}
		defer reader.Close()
		result, err = io.ReadAll(reader)
		if err != nil {
			return "", err
		}
	} else {
		result, err = io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
	}
	return string(result), nil
}

func isDebug() bool {
	_, found := os.LookupEnv("GOLIBRE_DEBUG")
	return found
}

func printRequest(req *http.Request) {
	fmt.Println("> ", req.Method, req.URL, req.Proto)
	for key, value := range req.Header {
		fmt.Println("> ", key, strings.Join(value, " "))
	}
	fmt.Println("> ")
}

func printResponse(resp *http.Response, body string) {
	fmt.Println("< ", resp.Proto, resp.Status)
	for key, value := range resp.Header {
		fmt.Println("< ", key+":", strings.Join(value, " "))
	}
	fmt.Printf("\n%s\n", body)
}
