package httpexpect_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"go.pr0ger.dev/x/httpexpect"
	"go.pr0ger.dev/x/httpexpect/mocks"
)

//go:generate mockgen -package mocks -destination mocks/testingT.go . TestingT

type TestExpectationSuite struct {
	suite.Suite

	ctrl *gomock.Controller
	t    *mocks.MockTestingT
}

func (suite *TestExpectationSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())

	suite.t = mocks.NewMockTestingT(suite.ctrl)
	suite.t.EXPECT().Helper().AnyTimes()
}

func (suite *TestExpectationSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *TestExpectationSuite) TestJSONObject() {
	suite.t.EXPECT().Errorf(gomock.Any(), gomock.Any()).MinTimes(2)
	suite.t.EXPECT().FailNow().MinTimes(2)

	invalidJSON := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`"invalid json`))
	})
	arrayJSON := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`[1, 2, 3]`))
	})
	objectJSON := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`{"key": 1}`))
	})

	httpexpect.Get(suite.t, invalidJSON).JSONObject()
	httpexpect.Get(suite.t, arrayJSON).JSONObject()
	suite.InDelta(1, httpexpect.Get(suite.t, objectJSON).JSONObject().Number("key").Value(), 0.0)
}

func (suite *TestExpectationSuite) TestJSONArray() {
	suite.t.EXPECT().Errorf(gomock.Any(), gomock.Any()).MinTimes(2)
	suite.t.EXPECT().FailNow().MinTimes(2)

	invalidJSON := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`"invalid json`))
	})
	objectJSON := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`{"key"": 123}`))
	})
	arrayJSON := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`[1, 2, 3]`))
	})

	httpexpect.Get(suite.t, invalidJSON).JSONArray()
	httpexpect.Get(suite.t, objectJSON).JSONArray()
	suite.InDelta(3, httpexpect.Get(suite.t, arrayJSON).JSONArray().Len().Value(), 0.0)
}

func (suite *TestExpectationSuite) TestNoContent() {
	suite.t.EXPECT().Errorf(gomock.Any(), gomock.Any())

	called := false
	stubHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte{0x00})
		called = true
	})

	httpexpect.Get(suite.t, stubHandler).NoContent()

	suite.True(called)
}

func (suite *TestExpectationSuite) TestStatus() {
	suite.t.EXPECT().Errorf(gomock.Any(), gomock.Any())

	called := false
	stubHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		called = true
	})

	httpexpect.Get(suite.t, stubHandler).Status(http.StatusOK)

	suite.True(called)
}

func (suite *TestExpectationSuite) TestStatusRequired() {
	suite.t.EXPECT().Errorf(gomock.Any(), gomock.Any())
	suite.t.EXPECT().FailNow()

	called := false
	stubHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		called = true
	})

	httpexpect.Get(suite.t, stubHandler).Require().Status(http.StatusOK)

	suite.True(called)
}

func TestExpectation(t *testing.T) {
	suite.Run(t, new(TestExpectationSuite))
}
