package entities

type FileInfo struct {
	Length      int64
	PieceLength int64
	Name        string
	Pieces      []SHAhash
}

type Torrent struct {
	Announce string
	Info     FileInfo
	InfoRaw  map[string]interface{}
  InfoHash [20]byte
}
