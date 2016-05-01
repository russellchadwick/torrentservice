package torrentservice

import (
	log "github.com/Sirupsen/logrus"
	"github.com/russellchadwick/rpc"
	pb "github.com/russellchadwick/torrentservice/proto"
	"golang.org/x/net/context"
	"time"
)

// Client is used to speak with the torrent service via rpc
type Client struct{}

// AddTorrent adds a torrent to transmission via RPC using a url to a file or a magnet link
func (t *Client) AddTorrent(url string) (*pb.AddTorrentResponse, error) {
	log.Debug("-> client.AddTorrent")
	start := time.Now()

	client := rpc.Client{}
	clientConn, err := client.Dial("torrent")
	if err != nil {
		log.WithField("error", err).Error("error during dial")
		return nil, err
	}
	defer func() {
		closeErr := clientConn.Close()
		if closeErr != nil {
			log.WithField("error", closeErr).Error("error during close")
		}
	}()

	torrentClient := pb.NewTorrentClient(clientConn)

	addTorrentResponse, err := torrentClient.AddTorrent(context.Background(), &pb.AddTorrentRequest{
		Url: url,
	})

	if err != nil {
		log.WithField("error", err).Error("error from rpc")
		return nil, err
	}

	elapsed := float64(time.Since(start)) / float64(time.Microsecond)
	log.WithField("elapsed", elapsed).Debug("<- client.AddTorrent")

	return addTorrentResponse, nil
}
