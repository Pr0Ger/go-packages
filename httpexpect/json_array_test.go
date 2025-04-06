package httpexpect

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type JSONArraySuite struct {
	JSONValueSuite
}

func (s *JSONArraySuite) TestOutOfBounds() {
	s.t.EXPECT().Errorf(gomock.Eq(`index %d is out of bounds (len=%d)`), 0, 0).MinTimes(6)
	s.t.EXPECT().FailNow().MinTimes(1)

	obj := JSONArray{
		expectation: s.expectation,
		value:       JSONArrayType{},
	}

	obj.Number(0)
	obj.String(0)
	obj.Bool(0)
	obj.Array(0)
	obj.Object(0)
	obj.Null(0)
}

func (s *JSONArraySuite) TestInvalidType() {
	s.t.EXPECT().FailNow().MinTimes(1)
	s.t.EXPECT().Errorf(gomock.Eq(`element at %d is not a number`), 0)
	s.t.EXPECT().Errorf(gomock.Eq(`element at %d is not a string`), 1)
	s.t.EXPECT().Errorf(gomock.Eq(`element at %d is not a bool`), 2)
	s.t.EXPECT().Errorf(gomock.Eq(`element at %d is not an array`), 3)
	s.t.EXPECT().Errorf(gomock.Eq(`element at %d is not an object`), 4)
	s.t.EXPECT().Errorf(gomock.Eq(`element at %d is not a null`), 5)

	obj := JSONArray{
		expectation: s.expectation,
		value: JSONArrayType{
			"",               // 0
			true,             // 1
			JSONArrayType{},  // 2
			JSONObjectType{}, // 3
			nil,              // 4
			0,                // 5
		},
	}

	obj.Number(0)
	obj.String(1)
	obj.Bool(2)
	obj.Array(3)
	obj.Object(4)
	obj.Null(5)
}

func (s *JSONArraySuite) TestLen() {
	obj := JSONArray{
		expectation: s.expectation,
		value: JSONArrayType{
			1,
		},
	}

	arrLen := obj.Len()
	s.InDelta(1, arrLen.Value(), 0.0)
}

func (s *JSONArraySuite) TestNumber() {
	array := JSONArray{
		expectation: s.expectation,
		value: JSONArrayType{
			float64(1337),
		},
	}

	s.NotNil(array.Number(0))
}

func (s *JSONArraySuite) TestString() {
	array := JSONArray{
		expectation: s.expectation,
		value: JSONArrayType{
			"test",
		},
	}

	s.NotNil(array.String(0))
}

func (s *JSONArraySuite) TestBool() {
	array := JSONArray{
		expectation: s.expectation,
		value: JSONArrayType{
			true,
		},
	}

	s.NotNil(array.Bool(0))
}

func (s *JSONArraySuite) TestArray() {
	array := JSONArray{
		expectation: s.expectation,
		value: JSONArrayType{
			JSONArrayType{},
		},
	}

	s.NotNil(array.Array(0))
}

func (s *JSONArraySuite) TestObject() {
	array := JSONArray{
		expectation: s.expectation,
		value: JSONArrayType{
			JSONObjectType{},
		},
	}

	s.NotNil(array.Object(0))
}

func (s *JSONArraySuite) TestNull() {
	array := JSONArray{
		expectation: s.expectation,
		value: JSONArrayType{
			nil,
		},
	}

	array.Null(0)
}

func TestJSONArray(t *testing.T) {
	suite.Run(t, new(JSONArraySuite))
}
