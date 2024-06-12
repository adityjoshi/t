package torrentfile

import (
	"bufio"
	"github.com/adityjoshi/t/entities"
	"os"
)

func ParseTorrentFile(filePath string) (*entities.Torrent, error) {
	fp, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	return getTorrent(bufio.NewReader(fp))
}
