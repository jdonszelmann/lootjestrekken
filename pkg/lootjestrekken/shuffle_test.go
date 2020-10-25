package lootjestrekken

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShuffle(t *testing.T) {
	for i := 0; i < 100000; i++ {
		a := []string{"0","1","2","3","4","5","6","7","8","9"}
		b := Derange(a)

		for index, i := range b {
			assert.NotEqual(t, i, a[index])
		}
	}
}
