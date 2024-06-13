package torrentfile

import (
	"bufio"
	"github.com/adityjoshi/t/entities"
	"os"
)

// splitting the pieces 

func batch(data []byte, batch int) []entities.SHAhash {
  var result []entities.SHAhash

}

func ParseTorrentFile(filePath string) (*entities.Torrent, error) {
	fp, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	return getTorrent(bufio.NewReader(fp))
}
