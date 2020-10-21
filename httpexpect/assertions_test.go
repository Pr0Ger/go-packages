package httpexpect_test

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

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

	called := false
	invalidJSON := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`"invalid json`))
		called = true
	})
	arrayJSON := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`[1, 2, 3]`))
		called = true
	})

	httpexpect.Get(suite.t, invalidJSON).JSONObject()
	suite.True(called)

	httpexpect.Get(suite.t, arrayJSON).JSONObject()
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
