package hash

import (
	"hash/crc32"
)

func HashCrc32(input []byte) uint32 {
	return crc32.ChecksumIEEE(input)
}

func HashCrc32String(input string) uint32 {
	return HashCrc32([]byte(input))
}
