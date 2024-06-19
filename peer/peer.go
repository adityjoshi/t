package peer

import (
	"encoding/binary"
	"errors"
	"net"

	"github.com/adityjoshi/t/entities"
)

func Unmarshal(peersBin []byte) ([]entities.Peer, error) {
	// Peers or peersSize is another long binary blob which holds the ip address and port the first 4 bit stores the ip and the last 2 store port
	peerSize := 6
	numPeers := len(peersBin) / peerSize
	if len(peersBin)%peerSize != 0 {
		return nil, errors.New("Invalid Error")
	}
	peers := make([]entities.Peer, numPeers)
	for i := 0; i < numPeers; i++ {
		offset := i * peerSize
		peers[i].IP = net.IP(peersBin[offset : offset+4]).String()
		peers[i].Port = int64(binary.BigEndian.Uint16([]byte(peersBin[offset+4 : offset+6])))
	}
	return peers, nil
}
