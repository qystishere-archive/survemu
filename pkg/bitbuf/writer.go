package bitbuf

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"unsafe"
)

type Writer struct {
	internalBuffer []byte
	totalBits uint
	currentBit uint
	bitsWritten uint
}

func (writer *Writer) Data() []byte {
	if writer.BytesWritten() == 0 {
		return make([]byte, 0)
	}
	return writer.internalBuffer[:writer.BytesWritten()]
}

func (writer *Writer) BitsWritten() uint {
	return writer.bitsWritten
}

func (writer *Writer) BytesWritten() int {
	return int(math.Ceil(float64(writer.bitsWritten) / 8))
}

func (writer *Writer) Seek(position uint) {
	if writer.totalBits > position {
		writer.currentBit = writer.totalBits
	}
	writer.currentBit = position
}

func (writer *Writer) WriteByte(val byte) {
	writer.WriteUnsignedBitInt32(uint32(val), uint(unsafe.Sizeof(val)) << 3)
}

func (writer *Writer) WriteBytes(val []byte) {
	for _,b := range val {
		writer.WriteByte(b)
	}
}

func (writer *Writer) WriteBool(is bool) {
	if is {
		writer.WriteUnsignedBitInt32(1, 1)
	} else {
		writer.WriteUnsignedBitInt32(0, 1)
	}
}

func (writer *Writer) WriteInt8(val int8) {
	writer.WriteSignedBitInt32(int32(val), uint(unsafe.Sizeof(val)) << 3)
}

func (writer *Writer) WriteUint8(val uint8) {
	writer.WriteUnsignedBitInt32(uint32(val), uint(unsafe.Sizeof(val)) << 3)
}

func (writer *Writer) WriteFloat(val, min, max float32, numBits uint) {
	r := 1 << numBits -1
	o := clamp(val, min, max)
	n := (o - min) / (max - min)
	s := uint32(float32(n * float32(r)) + .5)
	writer.WriteUnsignedBitInt32(s, numBits)
}

func (writer *Writer) WriteInt16(val int16) {
	writer.WriteSignedBitInt32(int32(val), uint(unsafe.Sizeof(val)) << 3)
}

func (writer *Writer) WriteUint16(val uint16) {
	writer.WriteUnsignedBitInt32(uint32(val), uint(unsafe.Sizeof(val)) << 3)
}

func (writer *Writer) WriteInt32(val int32) {
	writer.WriteSignedBitInt32(int32(val), uint(unsafe.Sizeof(val)) << 3)
}

func (writer *Writer) WriteUint32(val uint32) {
	writer.WriteUnsignedBitInt32(uint32(val), uint(unsafe.Sizeof(val)) << 3)
}

func (writer *Writer) WriteInt64(val int64) {
	writer.WriteUint64(uint64(val))
}

func (writer *Writer) WriteUint64(val uint64) {
	raw := make([]byte, 8)
	binary.LittleEndian.PutUint64(raw, uint64(val))
	writer.WriteUnsignedBitInt32(uint32(binary.LittleEndian.Uint32(raw[:4])), uint(unsafe.Sizeof(val)) << 3)
	writer.WriteUnsignedBitInt32(uint32(binary.LittleEndian.Uint32(raw[4:8])), uint(unsafe.Sizeof(val)) << 3)
}

func (writer *Writer) WriteString(val string) {
	for _,b := range []byte(val) {
		writer.WriteByte(b)
	}
	writer.WriteByte(0)
}

func (writer *Writer) WriteUnsignedBitInt32(data uint32, numBits uint) {
	// Force the sign-extension bit to be correct even in the case of overflow.
	//nValue := uint(data)
	//nPreserveBits := (0x7FFFFFFF >> (32 - numBits))
	//nSignExtension := (nValue >> 31) & ^nPreserveBits
	//nValue &= nPreserveBits
	//nValue |= nSignExtension

	writer.writeInternal(uint32(data), numBits, false)
}

func (writer *Writer) WriteSignedBitInt32(data int32, numBits uint) {
	// Force the sign-extension bit to be correct even in the case of overflow.
	nValue := int(data)
	nPreserveBits := (0x7FFFFFFF >> (32 - numBits))
	nSignExtension := (nValue >> 31) & ^nPreserveBits
	nValue &= nPreserveBits
	nValue |= nSignExtension

	writer.writeInternal(uint32(nValue), numBits, false)
}

func (writer *Writer) writeInternal(curData uint32, numBits uint, checkRange bool) error {
	if err := writer.ensureInBounds(numBits); err != nil {
		writer.currentBit = writer.totalBits
		return err
	}

	iCurBitMasked := writer.currentBit & 31
	iDWord := uint32(writer.currentBit >> 5)
	if writer.currentBit == writer.bitsWritten {
		writer.bitsWritten += numBits
	}
	writer.currentBit += numBits

	// Mask in a dword.
	//Assert((iDWord * 4 + sizeof(long)) <= (unsigned int)m_nDataBytes)
	pOut := []uint32{
		bytesToUint32(writer.internalBuffer[(iDWord*4):(iDWord*4)+4]),
		bytesToUint32(writer.internalBuffer[(iDWord*4)+4:(iDWord*4)+8]),
	}

	// Rotate data into dword alignment
	curData = (curData << iCurBitMasked) | (curData >> (32 - iCurBitMasked))

	// Calculate bitmasks for first and second word
	temp := uint(1 << (numBits - 1))
	mask1 := uint32((temp * 2 - 1) << iCurBitMasked)
	mask2 := uint32((temp - 1) >> (31 - iCurBitMasked))

	// Only look beyond current word if necessary (avoid access violation)
	i := mask2 & 1
	dword1 := pOut[0]
	dword2 := pOut[i]

	// Drop bits into place
	dword1 ^= (mask1 & (curData ^ dword1))
	dword2 ^= (mask2 & (curData ^ dword2))

	// Note reversed order of writes so that dword1 wins if mask2 == 0 && i == 0
	binary.LittleEndian.PutUint32(writer.internalBuffer[(iDWord*4) + (i*4):(iDWord*4) + (i*4) + 4], dword2)
	binary.LittleEndian.PutUint32(writer.internalBuffer[(iDWord*4):(iDWord*4) + 4], dword1)

	return nil
}

func (writer *Writer) ensureInBounds(numBits uint) error {
	if writer.currentBit + numBits > writer.totalBits {
		return errors.New(fmt.Sprintf("bitbuf attempt oob write by %d bits", (writer.currentBit + numBits) - writer.totalBits))
	}
	return nil
}


func NewWriter(length int) *Writer {
	return & Writer{
		internalBuffer: make([]byte, length + 4),
		totalBits:	    uint(length * 8) + 32,
		currentBit:     0,
	}
}
