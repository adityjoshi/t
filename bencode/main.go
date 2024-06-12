package bencode

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
)

func Decode(read *bufio.Reader) (interface{}, error) {
	ch, err := read.ReadByte()
	if err != nil {
		return nil, err
	}

	switch ch {
	case 'i':
		var buffer []byte
		for {
			ch, err := read.ReadByte()
			if err != nil {
				return nil, err
			}
			// if i stumble upon 'e', dict is complete i will return
			if ch == 'e' {
				value, err := strconv.ParseInt(string(buffer), 10, 64)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("invalid integer %s", string(buffer)))
				}
				return value, nil
			}
			buffer = append(buffer, ch)
		}
	case 'l':
		var listValueHolder []interface{}
		for {
			ch, err := read.ReadByte()
			if err != nil {
				return nil, err
			}
			// if i stumble upon 'e', list complete, and we return
			if ch == 'e' {
				return listValueHolder, err
			}
			// read the key
			read.UnreadByte()
			data, err := Decode(read)
			if err != nil {
				return nil, err
			}
			listValueHolder = append(listValueHolder, data)
		}
	case 'd':
		dictHolder := map[string]interface{}{}
		for {
			ch, err := read.ReadByte()
			if err != nil {
				return nil, err
			}
			// if we stumble upon 'e', dictionary complete, and we return
			if ch == 'e' {
				return dictHolder, nil
			}
			read.UnreadByte()
			data, err := Decode(read)
			if err != nil {
				return nil, err
			}

			// key has to be a string, if not then throw err
			key, ok := data.(string)
			if !ok {
				return nil, errors.New("key of the dictionary is not a string")
			}
			// read the value
			value, err := Decode(read)
			if err != nil {
				return nil, err
			}
			// put key and value in dictionary
			dictHolder[key] = value
		}
		// string
	default:
		read.UnreadByte()

		var lengthBuf []byte
		for {
			ch, err := read.ReadByte()
			if err != nil {
				return nil, err
			}
			if ch == ':' {
				break
			}
			lengthBuf = append(lengthBuf, ch)
		}
		length, err := strconv.Atoi(string(lengthBuf))
		if err != nil {
			return nil, errors.New(fmt.Sprintf("invaid integer %s", string(lengthBuf)))
		}
		var strBuf []byte
		for i := 0; i < length; i++ {
			ch, err := read.ReadByte()
			if err != nil {
				return nil, err
			}
			strBuf = append(strBuf, ch)
		}
		return string(strBuf), err
	}
}
