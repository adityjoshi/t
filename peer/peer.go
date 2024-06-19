package peer

import "github.com/adityjoshi/t/entities"


func Unmarshal(peersBin []byte) (entities.Peer, error) {
  // Peers or peersSize is another long binary blob which holds the ip address and port the first 4 bit stores the ip and the last 2 store port
  peerSize := 6 

}
