package masterpassword

import (
	"encoding/binary"
	"hash/crc32"
	"math/rand"
)

type numberGenerator struct {
	bytes  []byte
	random *rand.Rand
}

func newNumberGenerator(key []byte) *numberGenerator {
	seed := crc32.ChecksumIEEE(key)
	s := rand.NewSource(int64(seed))
	r := rand.New(s)
	return &numberGenerator{
		bytes:  key,
		random: r,
	}
}

func (g *numberGenerator) next() int {
	if len(g.bytes) >= 2 {
		dbyte := []byte{g.bytes[0], g.bytes[1], 0, 0}
		g.bytes = g.bytes[2:]
		data := binary.BigEndian.Uint32(dbyte)
		return int(data)
	}
	return g.random.Int()
}
