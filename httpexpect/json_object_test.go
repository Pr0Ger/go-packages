package httpexpect

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type JSONObjectSuite struct {
	JSONValueSuite
}

func (s JSONObjectSuite) TestNotExists() {
	s.t.EXPECT().Errorf(gomock.Eq(`key "%s" do not exists`), "key").MinTimes(6)
	s.t.EXPECT().FailNow().MinTimes(1)

	obj := JSONObject{
		expectation: s.expectation,
		value:       JSONObjectType{},
	}

	obj.Number("key")
	obj.String("key")
	obj.Bool("key")
	obj.Array("key")
	obj.Object("key")
	obj.Null("key")
}

func (s JSONObjectSuite) TestInvalidType() {
	s.t.EXPECT().FailNow().MinTimes(1)
	s.t.EXPECT().Errorf(gomock.Eq(`key "%s" is not an number`), "not_number")
	s.t.EXPECT().Errorf(gomock.Eq(`key "%s" is not a string`), "not_string")
	s.t.EXPECT().Errorf(gomock.Eq(`key "%s" is not a bool`), "not_bool")
	s.t.EXPECT().Errorf(gomock.Eq(`key "%s" is not an array`), "not_array")
	s.t.EXPECT().Errorf(gomock.Eq(`key "%s" is not an object`), "not_obj")
	s.t.EXPECT().Errorf(gomock.Eq(`key "%s" is not null`), "not_null")

	obj := JSONObject{
		expectation: s.expectation,
		value: JSONObjectType{
			"not_number": "",
			"not_string": true,
			"not_bool":   JSONArrayType{},
			"not_array":  JSONObjectType{},
			"not_obj":    nil,
			"not_null":   0,
		},
	}

	obj.Number("not_number")
	obj.String("not_string")
	obj.Bool("not_bool")
	obj.Array("not_array")
	obj.Object("not_obj")
	obj.Null("not_null")
}

func (s JSONObjectSuite) TestHasKey() {
	s.t.EXPECT().Errorf(gomock.Eq(`key "%s" do not exists`), "not_key").MaxTimes(1)

	obj := JSONObject{
		expectation: s.expectation,
		value: JSONObjectType{
			"key": nil,
		},
	}

	obj.HasKey("key")
	obj.HasKey("not_key")
}

func (s JSONObjectSuite) TestNumber() {
	obj := JSONObject{
		expectation: s.expectation,
		value: JSONObjectType{
			"number": float64(1337),
		},
	}

	s.NotNil(obj.Number("number"))
}

func (s JSONObjectSuite) TestString() {
	obj := JSONObject{
		expectation: s.expectation,
		value: JSONObjectType{
			"string": "test",
		},
	}

	s.NotNil(obj.String("string"))
}

func (s JSONObjectSuite) TestBool() {
	obj := JSONObject{
		expectation: s.expectation,
		value: JSONObjectType{
			"bool": true,
		},
	}

	s.NotNil(obj.Bool("bool"))
}

func (s JSONObjectSuite) TestArray() {
	obj := JSONObject{
		expectation: s.expectation,
		value: JSONObjectType{
			"array": JSONArrayType{},
		},
	}

	s.NotNil(obj.Array("array"))
}

func (s JSONObjectSuite) TestObject() {
	obj := JSONObject{
		expectation: s.expectation,
		value: JSONObjectType{
			"nested_obj": JSONObjectType{},
		},
	}

	s.NotNil(obj.Object("nested_obj"))
}

func (s JSONObjectSuite) TestNull() {
	obj := JSONObject{
		expectation: s.expectation,
		value: JSONObjectType{
			"key": nil,
		},
	}

	obj.Null("key")
}

func TestJSONObject(t *testing.T) {
	suite.Run(t, new(JSONObjectSuite))
}
