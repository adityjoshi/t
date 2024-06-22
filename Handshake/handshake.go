package handshake


type Handshake struct {
  Pstr string
  InfoHash [20]byte
  PeerID [20]byte
}

func New(infoHash, peerID [20]byte) *Handshake {
  return &Handshake{
    Pstr:     "BitTorrent protocol",
		InfoHash: infoHash,
		PeerID:   peerID,
  }
}


