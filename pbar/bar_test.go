package pbar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBarAsString(t *testing.T) {
	tests := []struct {
		name    string
		current uint
		total   uint
		width   uint
		result  string
	}{
		{
			name:    "empty",
			current: 0,
			total:   10,
			width:   30,
			result:  "[                       ] 0%",
		},
		{
			name:    "filled",
			current: 10,
			total:   10,
			width:   30,
			result:  "[=======================] 100%",
		},
		{
			name:    "filled",
			current: 5,
			total:   10,
			width:   30,
			result:  "[===========>           ] 50%",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bar := Bar{
				Current: test.current,
				Total:   test.total,
			}

			assert.Equal(t, test.result, bar.asString(test.width))
		})
	}
}
