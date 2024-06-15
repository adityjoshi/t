package tracker

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	bencode "github.com/jackpal/bencode-go"

	"github.com/adityjoshi/t/entities"
)

type TrackerClient struct {
	t *entities.Torrent
}

func NewTrackerClient(t *entities.Torrent) *TrackerClient {
	return &TrackerClient{t}
}

// these safe functions are created to safely convert the value to int64 and string

func safeInt64(v interface{}) int64 {
	x, ok := v.(int64)
	if !ok {
		return 0
	}
	return x
}

func safeString(v interface{}) string {
	x, ok := v.(string)
	if !ok {
		return ""
	}
	return x
}

func (tr *TrackerClient) FetchPeers() ([]entities.Peer, error) {
	base, err := url.Parse(tr.t.Announce)
	if err != nil {
		return nil, nil
	}

	/*
	   we are creatina a infoBuffer which a slice of bytes
	   and we are using jackpal bencode to encode the tr.t.InfoRaw to the infoBuffer slice
	   and checking the error
	   then we are calculating the sha1 hash value of the infoBuffer slice and storing the value in the infoHash which is a 20 byte array
	*/

	var infoBuffer bytes.Buffer
	err = bencode.Marshal(&infoBuffer, tr.t.InfoRaw)
	if err != nil {
		return nil, err
	}
	infoHash := sha1.Sum(infoBuffer.Bytes())

	// request params
	params := url.Values{
		"info_hash":  []string{string(infoHash[:])},
		"peer_id":    []string{string("~T Torrent~")},
		"port":       []string{strconv.Itoa(6881)},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"left":       []string{strconv.Itoa(int(tr.t.Info.Length))},
	}
	base.RawQuery = params.Encode()
	url := base.String()

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	tResponse, err := prepareTrackerResponse(resp.Body)
	if err != nil {
		return nil, err
	}

	if tResponse.FailureReason != "" {
		return nil, errors.New(tResponse.FailureReason)
	}

	return tResponse.Peers, nil
}
