package cryptosource

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"math/rand"
)

type source struct{}

// Uint64 returns a random 64-bit value as a uint64.
func (s source) Uint64() uint64 {
	var buf [8]byte
	if _, err := cryptorand.Read(buf[:]); err != nil {
		panic("cryptosource error: " + err.Error())
	}

	return binary.LittleEndian.Uint64(buf[:])
}

// Int63 returns a non-negative random 63-bit integer as an int64.
func (s source) Int63() int64 {
	var buf [8]byte
	if _, err := cryptorand.Read(buf[:]); err != nil {
		panic("cryptosource error: " + err.Error())
	}

	return int64(binary.LittleEndian.Uint64(buf[:]) >> 1)
}

// Seed should not be called for cryptosource.
func (s source) Seed(int64) {
	panic("this source should not be seeded")
}

// NewSource returns a new random Source which uses randomness from crypto/math.
func NewSource() rand.Source {
	return source{}
}
