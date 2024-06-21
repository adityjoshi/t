package entities
import (
"github.com/jackpal/bencode-go"
"crypto/sha1"
"bytes"
)

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


func (t *Torrent) CalculateInfoHash() error {
	var infoBuffer bytes.Buffer
  err := bencode.Marshal(&infoBuffer, t.InfoRaw)
	if err != nil {
		return err
	}
   t.InfoHash = sha1.Sum(infoBuffer.Bytes())
  return nil
}
