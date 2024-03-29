package bitbuf

import (
	"bytes"
	"encoding/binary"
)

func bytesToUint32(data []byte) (ret uint32) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)

	return ret
}

func bytesToInt64(data []byte) (ret int64) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)

	return ret
}

func bytesToFloat32(data []byte) (ret float32) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)

	return ret
}

func bytesToFloat64(data []byte) (ret float64) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)

	return ret
}

func clamp(e, t, a float32) float32 {
	if e < a {
		if e > t {
			return e
		} else {
			return t
		}
	} else {
		return a
	}
}