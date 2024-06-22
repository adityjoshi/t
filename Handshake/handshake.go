package handshake

type Handshake struct {
	Pstr     string
	InfoHash [20]byte
	PeerID   [20]byte
}

func New(infoHash, peerID [20]byte) *Handshake {
	return &Handshake{
		Pstr:     "BitTorrent protocol",
		InfoHash: infoHash,
		PeerID:   peerID,
	}
}

/*
When we serialize this is how a handshake string looks like
\x13BitTorrent protocol\x00\x00\x00\x00\x00\x00\x00\x00\x86\xd4\xc8\x00\x24\xa4\x69\xbe\x4c\x50\xbc\x5a\x10\x2c\xf7\x17\x80\x31\x00\x74-TR2940-k8hj0wgej6ch
*/
func (h *Handshake) Serialize() []byte {
  buf := make([]byte, len(h.Pstr)+49)
  buf[0] = byte(len(h.Pstr))
  curr := 1
  curr += copy(buf[curr:],h.Pstr)
  curr += copy(buf[curr:],make([]byte,8)) // 8 reserved bytes
  curr += copy(buf[curr:],h.InfoHash[:])
  curr += copy(buf[curr:],h.PeerID[:])
  return buf

}
