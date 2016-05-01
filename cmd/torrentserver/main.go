package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/russellchadwick/rpc"
	pb "github.com/russellchadwick/torrentservice/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type server struct{}

func (t *server) AddTorrent(ctx context.Context, in *pb.AddTorrentRequest) (*pb.AddTorrentResponse, error) {
	log.WithField("url", in.Url).Debug("-> server.AddTorrent")
	start := time.Now()

	discovery, err := rpc.NewDiscovery()
	if err != nil {
		log.Fatalf("failed to create discovery service: %v", err)
		return nil, err
	}

	name := "transmission"
	address, err := discovery.GetRandomServiceAddress(name)
	if err != nil {
		log.WithField("name", name).WithField("error", err).Error("failed to get service from discovery")
		return nil, err
	}
	url := "http://" + *address + "/transmission/rpc"

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	sessionID, err := getSessionID(httpClient, url)
	if err != nil {
		return nil, err
	}

	transmissionRequest, err := getTransmissionRequest(in.Url)
	if err != nil {
		return nil, err
	}

	responseBody, err := postTransmissionRequest(httpClient, url, transmissionRequest, *sessionID)
	if err != nil {
		return nil, err
	}

	torrentAddResponseArguments, err := decodeResponse(responseBody)
	if err != nil {
		return nil, err
	}

	elapsed := float64(time.Since(start)) / float64(time.Microsecond)
	log.WithField("elapsed", elapsed).Debug("<- server.AddTorrent")

	return &pb.AddTorrentResponse{
		Id:   torrentAddResponseArguments.ID,
		Name: torrentAddResponseArguments.Name,
		Hash: torrentAddResponseArguments.HashString,
	}, nil
}

type transmissionRequest struct {
	Method    string      `json:"method"`
	Arguments interface{} `json:"arguments"`
}

type transmissionResponse struct {
	Result    string                     `json:"result"`
	Arguments map[string]json.RawMessage `json:"arguments"`
}

type torrentAddRequestArguments struct {
	Filename string `json:"filename"`
}

type torrentAddResponseArguments struct {
	HashString string `json:"hashString"`
	ID         uint32 `json:"id"`
	Name       string `json:"name"`
}

func getSessionID(httpClient *http.Client, url string) (*string, error) {
	getResponse, err := httpClient.Get(url)
	if err != nil {
		log.WithField("error", err).Error("failed to get http get response")
		return nil, err
	}
	defer func() {
		closeErr := getResponse.Body.Close()
		if closeErr != nil {
			log.WithField("error", closeErr).Error("failed to close response body")
		}
	}()

	responseBody, err := ioutil.ReadAll(getResponse.Body)
	if err != nil {
		log.WithField("error", err).Error("failed to read http response")
		return nil, err
	}

	log.Debug("raw session id response: %s", string(responseBody))

	if getResponse.StatusCode != 409 {
		return nil, fmt.Errorf("expected a %v but got %v", 409, getResponse.StatusCode)
	}

	sessionID := getResponse.Header.Get("X-Transmission-Session-Id")
	if sessionID == "" {
		return nil, errors.New("unable to locate session header")
	}

	return &sessionID, nil
}

func postTransmissionRequest(httpClient *http.Client, url string, request io.Reader, sessionID string) (*[]byte, error) {
	httpRequest, err := http.NewRequest("POST", url, request)
	if err != nil {
		log.WithField("error", err).Error("failed to get http post request")
		return nil, err
	}
	httpRequest.Header.Set("X-Transmission-Session-Id", sessionID)
	httpRequest.Header.Set("Content-Type", "application/json")
	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		log.WithField("error", err).Error("failed to post http request")
		return nil, err
	}

	defer func() {
		closeErr := httpResponse.Body.Close()
		if closeErr != nil {
			log.WithField("error", closeErr).Error("failed to close response body")
		}
	}()

	responseBody, err := ioutil.ReadAll(httpResponse.Body)
	log.Debugf("raw reponse: %s", string(responseBody))
	if err != nil {
		log.WithField("error", err).Error("failed to read response body")
		return nil, err
	}

	return &responseBody, nil
}

func getTransmissionRequest(filename string) (*bytes.Buffer, error) {
	request := transmissionRequest{
		Method: "torrent-add",
		Arguments: torrentAddRequestArguments{
			Filename: filename,
		},
	}

	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(request)
	log.Debugf("raw request: %s", buffer.String())
	if err != nil {
		log.WithField("error", err).Error("failed to encode")
		return nil, err
	}

	return &buffer, nil
}

func decodeResponse(responseBody *[]byte) (*torrentAddResponseArguments, error) {
	var resp transmissionResponse
	err := json.Unmarshal(*responseBody, &resp)
	if err != nil {
		log.WithField("error", err).Error("failed to decode response")
		return nil, err
	}

	if resp.Result != "success" {
		return nil, fmt.Errorf("expected result of sucess but was: %v", resp.Result)
	}

	_, exists := resp.Arguments["torrent-added"]
	if !exists {
		return nil, errors.New("expected torrent-added dictionary is response but not found")
	}

	var respArgs torrentAddResponseArguments
	err = json.Unmarshal(resp.Arguments["torrent-added"], &respArgs)
	if err != nil {
		log.WithField("error", err).Error("failed to decode torrent added response")
		return nil, err
	}

	return &respArgs, nil
}

func main() {
	rpcServer := rpc.Server{}
	go serve(&rpcServer)
	defer func() {
		err := rpcServer.Stop()
		if err != nil {
			log.WithField("error", err).Error("error during stop")
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	signal.Notify(signalChan, syscall.SIGTERM)
	<-signalChan
	log.Info("Received shutdown signal")
}

func serve(rpcServer *rpc.Server) {
	err := rpcServer.Serve("torrent", func(grpcServer *grpc.Server) {
		pb.RegisterTorrentServer(grpcServer, &server{})
	})
	if err != nil {
		log.WithField("error", err).Error("error from rpc serve")
	}
}
