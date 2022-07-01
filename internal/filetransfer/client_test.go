package filetransfer

import (
	"context"
	"errors"
	"fransfer/internal/auth"
	"fransfer/internal/generated"
	"fransfer/internal/util"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

type ClientTestSuite struct {
	suite.Suite

	jwt auth.JWT

	name         string
	content      []byte
	defaultError error

	client client
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

func (s *ClientTestSuite) SetupSuite() {
	var err error

	privateKeyBytes, publicKeyBytes, err := util.GenerateKeyPair()
	s.NoError(err)

	s.jwt, err = auth.NewJWT(privateKeyBytes, publicKeyBytes)
	s.NoError(err)

	s.defaultError = errors.New("something went wrong")
	s.name = "file.txt"
	s.content = []byte{1, 2, 3}
	s.client = client{jwt: s.jwt}
}

func (s *ClientTestSuite) TestSend() {
	s.Run("file send successfully", func() {
		s.client.serviceClient = s.fileSendSuccessfullyBehaviourMock()

		s.Nil(s.client.Send(s.name, s.content))
	})

	s.Run("failed to init a stream", func() {
		s.client.serviceClient = s.failedToInitAStreamBehaviourMock()

		s.Equal(ErrOccurredDuringFileSend, s.client.Send(s.name, s.content))
	})

	s.Run("failed to send file name", func() {
		s.client.serviceClient = s.failedToSendFileNameBehaviourMock()

		s.Equal(ErrOccurredDuringFileSend, s.client.Send(s.name, s.content))
	})

	s.Run("failed to send file content", func() {
		s.client.serviceClient = s.failedToSendFileContentBehaviourMock()

		s.Equal(ErrOccurredDuringFileSend, s.client.Send(s.name, s.content))
	})

	s.Run("failed to close stream", func() {
		s.client.serviceClient = s.failedToCloseStreamBehaviourMock()

		s.Equal(ErrOccurredDuringFileSend, s.client.Send(s.name, s.content))
	})
}

func (s *ClientTestSuite) fileSendSuccessfullyBehaviourMock() *generated.FileTransferClientMock {
	return &generated.FileTransferClientMock{
		SendFunc: func(ctx context.Context, opts ...grpc.CallOption) (generated.FileTransfer_SendClient, error) {
			return &generated.FileTransfer_SendClientMock{
				SendFunc:         func(_ *generated.FileSendRequest) error { return nil },
				CloseAndRecvFunc: func() (*generated.FileSendResponse, error) { return nil, nil },
			}, nil
		},
	}
}

func (s *ClientTestSuite) failedToInitAStreamBehaviourMock() *generated.FileTransferClientMock {
	return &generated.FileTransferClientMock{
		SendFunc: func(ctx context.Context, opts ...grpc.CallOption) (generated.FileTransfer_SendClient, error) {
			return nil, s.defaultError
		},
	}
}

func (s *ClientTestSuite) failedToSendFileNameBehaviourMock() *generated.FileTransferClientMock {
	return &generated.FileTransferClientMock{
		SendFunc: func(ctx context.Context, opts ...grpc.CallOption) (generated.FileTransfer_SendClient, error) {
			return &generated.FileTransfer_SendClientMock{
				SendFunc: func(_ *generated.FileSendRequest) error { return s.defaultError },
			}, nil
		},
	}
}

func (s *ClientTestSuite) failedToSendFileContentBehaviourMock() *generated.FileTransferClientMock {
	return &generated.FileTransferClientMock{
		SendFunc: func(ctx context.Context, opts ...grpc.CallOption) (generated.FileTransfer_SendClient, error) {
			return &generated.FileTransfer_SendClientMock{
				SendFunc: func(fr *generated.FileSendRequest) error {
					if _, ok := fr.File.(*generated.FileSendRequest_Chunk); ok {
						return s.defaultError
					}
					return nil
				},
			}, nil
		},
	}
}

func (s *ClientTestSuite) failedToCloseStreamBehaviourMock() *generated.FileTransferClientMock {
	return &generated.FileTransferClientMock{
		SendFunc: func(ctx context.Context, opts ...grpc.CallOption) (generated.FileTransfer_SendClient, error) {
			return &generated.FileTransfer_SendClientMock{
				SendFunc:         func(fr *generated.FileSendRequest) error { return nil },
				CloseAndRecvFunc: func() (*generated.FileSendResponse, error) { return nil, s.defaultError },
			}, nil
		},
	}
}
