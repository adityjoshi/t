package message

import (
	"encoding/binary"
	"fmt"
	"io"
)

type messageID uint8

const (
	// MessageChoke chokes the receiver
	MessageChoke messageID = 0
	// MessageUnchoke unchokes the receiver
	MessageUnchoke messageID = 1
	// MessageInterested interested in receiving data
	MessageInterested messageID = 2
	// MessageUninterested not interesred in data
	MessageUninterested messageID = 3
	// MessageHave alerts the receiver that sender has downloaded the piece
	MessageHave messageID = 4
	// MessageBitfield encodes which piece sender has downloaded
	MessageBitfield messageID = 5
	// MessageRequest request a block of data from receiver
	MessageRequest messageID = 6
	// MessagePiece delivers a block of data to fulfill request
	MessagePiece messageID = 7
	// MessageCancel cancles the request
	MessageCancel messageID = 8
)

type Message struct {
	ID      messageID
	Payload []byte
}

// FormatRequest creates a Request message
func FormatRequest(index, begin, length int) *Message {
	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload[0:4], uint32(index))
	binary.BigEndian.PutUint32(payload[4:8], uint32(begin))
	binary.BigEndian.PutUint32(payload[8:12], uint32(length))
	return &Message{ID: MessageRequest, Payload: payload}
}

// FormatHave creates a Have message
func FormatHave(index int) *Message {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, uint32(index))
	return &Message{ID: MessageHave, Payload: payload}
}

func ParsePiece(index int, buf []byte, msg *Message) (int, error) {
	if msg.ID != MessagePiece {
		return 0, fmt.Errorf("Expected PIECE (ID %d), got ID %d", MessagePiece, msg.ID)
	}
	if len(msg.Payload) < 8 {
		return 0, fmt.Errorf("Payload too short. %d < 8", len(msg.Payload))
	}
	parsedIndex := int(binary.BigEndian.Uint32(msg.Payload[0:4]))
	if parsedIndex != index {
		return 0, fmt.Errorf("Expected index %d, got %d", index, parsedIndex)
	}
	begin := int(binary.BigEndian.Uint32(msg.Payload[4:8]))
	if begin >= len(buf) {
		return 0, fmt.Errorf("Begin offset too high. %d >= %d", begin, len(buf))
	}
	data := msg.Payload[8:]
	if begin+len(data) > len(buf) {
		return 0, fmt.Errorf(
			"Data too long [%d] for offset %d with length %d",
			len(data),
			begin,
			len(buf),
		)
	}
	copy(buf[begin:], data)
	return len(data), nil
}

func ParseHave(msg *Message) (int, error) {
	if msg.ID != MessageHave {
		return 0, fmt.Errorf("Expected HAVE (ID %d), got ID %d", MessageHave, msg.ID)
	}
	if len(msg.Payload) != 4 {
		return 0, fmt.Errorf("Expected payload length 4, got length %d", len(msg.Payload))
	}
	index := int(binary.BigEndian.Uint32(msg.Payload))
	return index, nil
}

/*
a message has three fields <length><id><optional Payload>
length is of size 4 , id is of size 1 and the remaining is paylod
*/
func (m *Message) Serialize() []byte {
	if m == nil {
		return make([]byte, 4)
	}
	length := uint32(len(m.Payload) + 1)
	buf := make([]byte, 4+length)
	binary.BigEndian.PutUint32(buf[0:4], length)
	buf[4] = byte(m.ID)
	copy(buf[5:], m.Payload)
	return buf
}

func Read(r io.Reader) (*Message, error) {
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lengthBuf)

	// keep message alive
	if length == 0 {
		return nil, nil
	}
	messageBuf := make([]byte, length)
	_, err = io.ReadFull(r, messageBuf)
	if err != nil {
		return nil, err
	}
	m := Message{
		ID:      messageID(messageBuf[0]),
		Payload: messageBuf[1:],
	}
	return &m, nil
}

func (m *Message) name() string {
	if m == nil {
		return "KeepAlive"
	}
	switch m.ID {
	case MessageChoke:
		return "Choke"
	case MessageUnchoke:
		return "Unchoke"
	case MessageInterested:
		return "Interested"
	case MessageUninterested:
		return "NotInterested"
	case MessageHave:
		return "Have"
	case MessageBitfield:
		return "Bitfield"
	case MessageRequest:
		return "Request"
	case MessagePiece:
		return "Piece"
	case MessageCancel:
		return "Cancel"
	default:
		return fmt.Sprintf("Unknown#%d", m.ID)
	}
}

func (m *Message) String() string {
	if m == nil {
		return m.name()
	}
	return fmt.Sprintf("%s [%d]", m.name(), len(m.Payload))
}
