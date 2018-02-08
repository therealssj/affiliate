package postgresql

import (
	"encoding/binary"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func uint64ToByteArray(num uint64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, num)
	return buf
}
