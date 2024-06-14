package torrentfile

import (
	"bufio"
	"errors"
	"os"

	"github.com/adityjoshi/t/bencode"
	"github.com/adityjoshi/t/entities"
)

// splitting the pieces

func batch(data []byte, batch int) []entities.SHAhash {
	var result []entities.SHAhash
	for i := 0; i < len(data); i += batch {
		hash := entities.SHAhash{}
		end := 1 + batch
		if end > len(data) {
			end = len(data)
		}
		copy(hash[:], data[i:end])
		result = append(result, hash)
	}
	return result
}

func getTorrent(bReader *bufio.Reader) (*entities.Torrent, error) {
	data, err := bencode.Decode(bReader)
	if err != nil {
		return nil, err
	}
	torrentData, ok := data.(map[string]interface{})
	if !ok {
		return nil, errors.New("Invalid torrent file")
	}
	torrentInfoData, ok := torrentData["info"].(map[string]interface{})
	if !ok {
		return nil, errors.New("Invalid torrent file")
	}

	torrent := &entities.Torrent{}
	torrent.Announce = torrentData["announce"].(string)
	torrent.Info = entities.FileInfo{
		PieceLength: torrentInfoData["PieceLength"].(int64),
		Length:      torrentInfoData["length"].(int64),
		Name:        torrentInfoData["name"].(string),
		Pieces:      batch([]byte(torrentInfoData["pieces"].(string)), 20),
	}
	torrent.InfoRaw = torrentData["info"].(map[string]interface{})
	return torrent, nil
}

func ParseTorrentFile(filePath string) (*entities.Torrent, error) {
	fp, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	return getTorrent(bufio.NewReader(fp))
}
