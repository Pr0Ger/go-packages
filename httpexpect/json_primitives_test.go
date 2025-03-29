package httpexpect

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type JSONNumberSuite struct {
	JSONValueSuite
}

func (s *JSONNumberSuite) TestValue() {
	s.t.EXPECT().Errorf(gomock.Any(), gomock.Any()).MaxTimes(0)

	val := JSONNumber{
		expectation: s.expectation,
		path:        ".key",
		value:       float64(1337),
	}

	s.InEpsilon(float64(1337), val.Value(), 0.0001)
}

func (s *JSONNumberSuite) TestEqual() {
	s.t.EXPECT().Errorf(
		gomock.Eq(`key "%s" is not equal %f. Actual value %f`),
		gomock.Eq(".key"),
		gomock.Eq(float64(1000)),
		gomock.Eq(float64(1337)),
	)

	val := JSONNumber{
		expectation: s.expectation,
		path:        ".key",
		value:       1337,
	}

	val.Equal(1337)
	val.Equal(1000)
}

func (s *JSONNumberSuite) TestEqualDelta() {
	s.t.EXPECT().Errorf(
		gomock.Eq(`key "%s" is not equal %f ± %f. Actual value %f`),
		gomock.Eq(".key"),
		gomock.Eq(float64(1000)),
		gomock.Eq(float64(10)),
		gomock.Eq(float64(1337)),
	)

	val := JSONNumber{
		expectation: s.expectation,
		path:        ".key",
		value:       1337,
	}

	val.EqualDelta(1336, 10)
	val.EqualDelta(1338, 10)
	val.EqualDelta(1000, 10)
}

func TestJSONNumber(t *testing.T) {
	suite.Run(t, new(JSONNumberSuite))
}

type JSONStringSuite struct {
	JSONValueSuite
}

func (s *JSONStringSuite) TestValue() {
	s.t.EXPECT().Errorf(gomock.Any(), gomock.Any()).MaxTimes(0)

	val := JSONString{
		expectation: s.expectation,
		path:        ".key",
		value:       "test",
	}

	s.Equal("test", val.Value())
}

func (s *JSONStringSuite) TestEqual() {
	s.t.EXPECT().Errorf(
		gomock.Eq(`key "%s" is not equal "%s". Actual value "%s"`),
		gomock.Eq(".key"),
		gomock.Eq("not-test"),
		gomock.Eq("test"),
	)

	val := JSONString{
		expectation: s.expectation,
		path:        ".key",
		value:       "test",
	}

	val.Equal("test")
	val.Equal("not-test")
}

func (s *JSONStringSuite) TestLen() {
	s.t.EXPECT().Errorf(gomock.Any(), gomock.Any()).MaxTimes(0)

	val := JSONString{
		expectation: s.expectation,
		path:        ".key",
		value:       "test",
	}

	s.InDelta(4, val.Len().Value(), 0.0)
}

func (s *JSONStringSuite) TestHasPrefix() {
	s.t.EXPECT().Errorf(
		gomock.Eq(`key "%s" doesn't have prefix "%s""`),
		gomock.Eq(".key"),
		gomock.Eq("bar"),
	)

	val := JSONString{
		expectation: s.expectation,
		path:        ".key",
		value:       "foo-bar",
	}

	val.HasPrefix("foo")
	val.HasPrefix("bar")
}

func (s *JSONStringSuite) TestHasSuffix() {
	s.t.EXPECT().Errorf(
		gomock.Eq(`key "%s" doesn't have suffix "%s""`),
		gomock.Eq(".key"),
		gomock.Eq("foo"),
	)

	val := JSONString{
		expectation: s.expectation,
		path:        ".key",
		value:       "foo-bar",
	}

	val.HasSuffix("foo")
	val.HasSuffix("bar")
}

func TestJSONString(t *testing.T) {
	suite.Run(t, new(JSONStringSuite))
}

type JSONBoolSuite struct {
	JSONValueSuite
}

func (s *JSONBoolSuite) TestValue() {
	s.t.EXPECT().Errorf(gomock.Any(), gomock.Any()).MaxTimes(0)

	val := JSONBool{
		expectation: s.expectation,
		path:        ".key",
		value:       true,
	}

	s.True(val.Value())
}

func (s *JSONBoolSuite) TestEqual() {
	s.t.EXPECT().Errorf(gomock.Any(), gomock.Any()).MaxTimes(0)

	val := JSONBool{
		expectation: s.expectation,
		path:        ".key",
		value:       false,
	}

	val.Equal(false)
}

func (s *JSONBoolSuite) TestHelpers() {
	s.t.EXPECT().Errorf(`key "%s" is not equal %t`, ".key", false)

	val := JSONBool{
		expectation: s.expectation,
		path:        ".key",
		value:       true,
	}
	val.True()
	val.False()
}

func TestJSONBool(t *testing.T) {
	suite.Run(t, new(JSONBoolSuite))
}
