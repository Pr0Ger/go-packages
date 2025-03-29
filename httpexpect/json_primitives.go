package httpexpect

import (
	"math"
	"strings"
)

type JSONNumber struct {
	expectation *Expectation

	path  string
	value float64
}

func (n *JSONNumber) Value() float64 {
	return n.value
}

func (n *JSONNumber) Equal(value float64) {
	if n.value != value {
		n.expectation.errorf(`key "%s" is not equal %f. Actual value %f`, n.path, value, n.value)
	}
}

func (n *JSONNumber) EqualDelta(value, delta float64) {
	if math.Abs(n.value-value) >= delta {
		n.expectation.errorf(`key "%s" is not equal %f ± %f. Actual value %f`, n.path, value, delta, n.value)
	}
}

type JSONString struct {
	expectation *Expectation

	path  string
	value string
}

func (s *JSONString) Value() string {
	return s.value
}

func (s *JSONString) Equal(value string) {
	if s.value != value {
		s.expectation.errorf(`key "%s" is not equal "%s". Actual value "%s"`, s.path, value, s.value)
	}
}

func (s *JSONString) Len() *JSONNumber {
	return &JSONNumber{
		expectation: s.expectation,
		path:        s.path + ".$len",
		value:       float64(len(s.value)),
	}
}

func (s *JSONString) HasPrefix(prefix string) {
	if !strings.HasPrefix(s.value, prefix) {
		s.expectation.errorf(`key "%s" doesn't have prefix "%s""`, s.path, prefix)
	}
}

func (s *JSONString) HasSuffix(suffix string) {
	if !strings.HasSuffix(s.value, suffix) {
		s.expectation.errorf(`key "%s" doesn't have suffix "%s""`, s.path, suffix)
	}
}

type JSONBool struct {
	expectation *Expectation

	path  string
	value bool
}

func (b *JSONBool) Value() bool {
	return b.value
}

func (b *JSONBool) Equal(value bool) {
	if b.value != value {
		b.expectation.errorf(`key "%s" is not equal %t`, b.path, value)
	}
}

func (b *JSONBool) True() {
	b.Equal(true)
}

func (b *JSONBool) False() {
	b.Equal(false)
}
