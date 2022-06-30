package filetransfer

import (
	"context"
	"fransfer/internal/auth"
	"fransfer/internal/generated"
	"path/filepath"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type (
	Client interface {
		Send(string, []byte) error
		Close()
	}

	client struct {
		jwt           auth.JWT
		conn          *grpc.ClientConn
		serviceClient generated.FileTransferServiceClient
	}
)

func NewClient(addr string, jwt auth.JWT) (Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return client{
		jwt:           jwt,
		conn:          conn,
		serviceClient: generated.NewFileTransferServiceClient(conn),
	}, nil
}

func (c client) Close() {
	c.conn.Close()
}

func (c client) Send(filePath string, content []byte) error {
	ctx, err := c.getContextWithToken()
	if err != nil {
		return err
	}

	stream, err := c.serviceClient.Send(ctx)
	if err != nil {
		return err
	}

	if err := stream.Send(c.buildFileSendRequestWithName(filePath)); err != nil {
		return err
	}
	for _, chunk := range c.split(content, 512) {
		if err := stream.Send(c.buildFileSendRequestWithChunk(chunk)); err != nil {
			return err
		}
	}

	if _, err = stream.CloseAndRecv(); err != nil {
		return err
	}
	return nil
}

func (c client) buildFileSendRequestWithName(filePath string) *generated.FileSendRequest {
	return &generated.FileSendRequest{File: &generated.FileSendRequest_Name{Name: filepath.Base(filePath)}}
}

func (c client) buildFileSendRequestWithChunk(chunk []byte) *generated.FileSendRequest {
	return &generated.FileSendRequest{File: &generated.FileSendRequest_Chunk{Chunk: chunk}}
}

func (c client) getContextWithToken() (context.Context, error) {
	token, err := c.jwt.CreateWithTTL(5 * time.Minute)
	if err != nil {
		return context.Background(), err
	}

	md := metadata.New(map[string]string{auth.HeaderApiKey: token})
	return metadata.NewOutgoingContext(context.Background(), md), nil
}

func (c client) split(buf []byte, lim int) [][]byte {
	chunk, chunks := make([]byte, 0), make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:])
	}
	return chunks
}
