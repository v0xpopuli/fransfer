package filetransfer

import (
	"errors"
	"fransfer/internal/generated"
	"io"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServerTestSuite struct {
	suite.Suite

	defaultError error

	service service
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (s *ServerTestSuite) SetupSuite() {
	s.defaultError = errors.New("something went wrong")

	s.service = NewService("").(service)
}

func (s *ServerTestSuite) TestSend() {
	s.Run("EOF reached", func() {
		streamMock := s.eofReachedBehaviourMock()

		s.Nil(s.service.Send(streamMock))
	})

	s.Run("error during file stream", func() {
		streamMock := s.errorDuringFileStreamBehaviourMock()

		s.Equal(ErrFileStream, s.service.Send(streamMock))
	})
}

func (s *ServerTestSuite) eofReachedBehaviourMock() *generated.FileTransfer_SendServerMock {
	return &generated.FileTransfer_SendServerMock{
		RecvFunc: func() (*generated.FileSendRequest, error) {
			return nil, io.EOF
		},
		SendAndCloseFunc: func(_ *generated.FileSendResponse) error {
			return nil
		},
	}
}

func (s *ServerTestSuite) errorDuringFileStreamBehaviourMock() *generated.FileTransfer_SendServerMock {
	return &generated.FileTransfer_SendServerMock{
		RecvFunc: func() (*generated.FileSendRequest, error) {
			return nil, s.defaultError
		},
		SendAndCloseFunc: func(_ *generated.FileSendResponse) error {
			return nil
		},
	}
}
