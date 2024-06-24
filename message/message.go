package message

import "encoding/binary"

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
