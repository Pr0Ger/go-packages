package httpexpect

import (
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"go.pr0ger.dev/x/httpexpect/mocks"
)

type JSONValueSuite struct {
	suite.Suite

	ctrl        *gomock.Controller
	t           *mocks.MockTestingT
	expectation *Expectation
}

func (s *JSONValueSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.t = mocks.NewMockTestingT(s.ctrl)
	s.t.EXPECT().Helper().AnyTimes()

	s.expectation = &Expectation{t: s.t}
}

func (s *JSONValueSuite) TearDownTest() {
	s.ctrl.Finish()
}
