package filetransfer

import (
	"fransfer/internal/generated"
	"io"
	"io/fs"
	"io/ioutil"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type (
	Service interface {
		generated.FileTransferServiceServer
	}

	service struct {
		outputDirectory string
	}
)

func NewService(outputDirectory string) Service {
	return service{outputDirectory: outputDirectory}
}

func (g service) Send(stream generated.FileTransferService_SendServer) error {
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
			return err
		}

		switch f := req.File.(type) {
		case *generated.FileSendRequest_Name:
			name = f.Name
		case *generated.FileSendRequest_Chunk:
			content = append(content, f.Chunk...)
		}
	}

	if err := ioutil.WriteFile(filepath.Join(g.outputDirectory, name), content, fs.ModePerm); err != nil {
		logrus.WithError(err).WithField("file_name", name).
			Error("Error occurred during attempt to save file")
		return err
	}

	logrus.WithField("file_name", name).Debug("File successfully saved")
	return stream.SendAndClose(&generated.FileSendResponse{})
}
