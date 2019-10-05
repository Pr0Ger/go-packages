package cryptosource

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleRand() {
	source := rand.New(NewSource())

	fmt.Println(source.Int63() != 0)
	// Output: true
}

func TestNew(t *testing.T) {
	source := NewSource()

	assert.Implements(t, (*rand.Source)(nil), source)
	assert.Implements(t, (*rand.Source64)(nil), source)
}

func TestSource_Int63(t *testing.T) {
	source := NewSource()

	for i := 0; i < 1000; i++ {
		value := source.Int63()

		assert.GreaterOrEqual(t, value, int64(0), ">= 0")
		assert.LessOrEqual(t, value, int64(math.MaxInt64), "< 2 ** 63")
	}
}

func TestSource_Uint64(t *testing.T) {
	source := NewSource().(rand.Source64)

	seen := make(map[uint64]struct{})
	for i := 0; i < 1000; i++ {
		value := source.Uint64()

		if _, ok := seen[value]; ok {
			t.Fail()
		}
		seen[value] = struct{}{}

		assert.GreaterOrEqual(t, value, uint64(0), ">= 0")
		assert.LessOrEqual(t, value, uint64(math.MaxUint64), "< 2 ** 64")
	}
}

func TestSource_Seed(t *testing.T) {
	source := NewSource()

	assert.Panics(t, func() {
		source.Seed(0)
	})

	rand.New(NewSource()).Uint32()
}

func BenchmarkSource_Uint64(b *testing.B) {
	source := rand.New(NewSource())

	for n := 0; n < b.N; n++ {
		source.Int63()
	}
}
