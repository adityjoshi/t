package message

import (
	"encoding/binary"
	"fmt"
	"io"
)

type messageID uint8

const (
	MessageChoke         messageID = 0
	MessageUnchoke       messageID = 1
	MessageInterested    messageID = 2
	MessageNotInterested messageID = 3
	MessageHave          messageID = 4
	MessageBitfield      messageID = 5
	MessageRequest       messageID = 6
	MessagePiece         messageID = 7
	MessageCancel        messageID = 8
)

type Message struct {
	ID      messageID
	Payload []byte
}

/*
We will create a serialize method to serialize the a message into the buffer
the patter it will follow will be = >    <length><message id><optional payload>
*/
func FormatRequest(index, begin, length int) *Message {
	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload[0:4], uint32(index))
	binary.BigEndian.PutUint32(payload[4:8], uint32(begin))
	binary.BigEndian.PutUint32(payload[8:12], uint32(length))
	return &Message{ID: MessageRequest, Payload: payload}
}

// FormatHave creates a HAVE message
func FormatHave(index int) *Message {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, uint32(index))
	return &Message{ID: MessageHave, Payload: payload}
}

// ParsePiece parses a PIECE message and copies its payload into a buffer
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
		return 0, fmt.Errorf("Data too long [%d] for offset %d with length %d", len(data), begin, len(buf))
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
	lengthBuffer := make([]byte, 4)
	_, err := io.ReadFull(r, lengthBuffer)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lengthBuffer)

	if length == 0 {
		return nil, nil
	}
	messageBuffer := make([]byte, length)
	_, err = io.ReadFull(r, messageBuffer)
	if err != nil {
		return nil, err
	}
	m := Message{
		ID:      messageID(messageBuffer[0]),
		Payload: messageBuffer[1:],
	}
	return &m, err
}

func (m *Message) name() string {
	if m == nil {
		return "Keep the network alive!"
	}
	switch m.ID {
	case MessageChoke:
		return "Choked"
	case MessageUnchoke:
		return "Unchoked"
	case MessageInterested:
		return "Interested"
	case MessageNotInterested:
		return "Uninterested"
	case MessageHave:
		return "Message Have"
	case MessageBitfield:
		return "Message Bitfield"
	case MessageRequest:
		return "Request"
	case MessagePiece:
		return "piece"
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
