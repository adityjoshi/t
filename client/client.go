package client

import (
	"net"
	"time"
  "github.com/adityjoshi/t/peer"
	"github.com/adityjoshi/t/entities"
)

type Client struct {
	Conn   net.Conn
	Choked bool
	// add bitfield
	Peer     entities.Peer
	infoHash [20]byte
	peerID  
}


func NewConn(peer entities.Peer,peerID,infoHash[20]byte) (*Client, error) {
  conn , err := net.DialTimeout("tcp",wrapperPeer.String(),3*time.Second) 
  if err != nil {
    return nil, err
  }
}
