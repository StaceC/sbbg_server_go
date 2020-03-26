package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetSequenceNumberGenerator(t *testing.T) {
	assert := assert.New(t)

	seq := []int{9, 1, 4, 10, 7, 5, 3}
	gen := NewSSNG(seq)

	for i := 0; i < len(seq); i++ {
		got := gen.GetInt()
		want := seq[i]
		assert.Equal(got, want)
	}

}

func TestRandomNumberGenerator(t *testing.T) {
	assert := assert.New(t)
	gen := NewRNG(MaxNum)

	for i := 0; i < 10000; i++ {
		got := gen.GetInt()
		assert.True(got <= MaxNum && got >= MinNum)
	}
}
