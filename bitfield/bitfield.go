package bitfield

type Bitfield []byte

func (bf Bitfield) PieceAvailable(index int) bool {
	byteIndex := index / 8
	bitIndex := index % 8
	// edge case
	if byteIndex < 0 || byteIndex >= len(bf) {
		return false
	}
	return bf[byteIndex]&(1<<uint(7-bitIndex)) != 0
}

func (bf Bitfield) SetAvailablePiece(index int) {
	byteIndex := index / 8
	bitIndex := index / 8
	// edge case
	if byteIndex < 0 || byteIndex >= len(bf) {
		return
	}
	bf[byteIndex] |= (1 << uint(7-bitIndex))
}
