package entities

type Peer struct {
  IP string
  Port int64
  PeerID string
}

type TrackerResponse struct {
  Complete int64
  Incomplete int64
  Interval int64
  Peer []Peer
  FailureReason string 
}
