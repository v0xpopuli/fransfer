package filetransfer

import (
	"errors"
	"fransfer/internal/generated"
	"io"
	"io/fs"
	"io/ioutil"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var (
	ErrFileStream = errors.New("error occurred during file stream")
	ErrFileSave   = errors.New("error occurred during file save")
)

type (
	Service interface {
		generated.FileTransferServer
	}

	service struct {
		outputDirectory string
	}
)

func NewService(outputDirectory string) Service {
	return service{outputDirectory: outputDirectory}
}

func (s service) Send(stream generated.FileTransfer_SendServer) error {
	var (
		name    string
		content []byte
	)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			logrus.Warn("EOF reached")
			break
		}
		if err != nil {
			logrus.WithError(err).WithField("file_name", name).
				Error("Error occurred during file content streaming")
			return ErrFileStream
		}

		switch f := req.File.(type) {
		case *generated.FileSendRequest_Name:
			name = f.Name
		case *generated.FileSendRequest_Chunk:
			content = append(content, f.Chunk...)
		}
	}

	if name != "" && len(content) != 0 {
		if err := ioutil.WriteFile(filepath.Join(s.outputDirectory, name), content, fs.ModePerm); err != nil {
			logrus.WithError(err).WithField("file_name", name).
				Error("Error occurred during attempt to save file")
			return ErrFileSave
		}
		logrus.WithField("file_name", name).Debug("File successfully saved")
	}

	return stream.SendAndClose(&generated.FileSendResponse{})
}
