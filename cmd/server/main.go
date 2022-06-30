package main

import (
	_ "embed"
	"flag"
	"fransfer/internal/auth"
	"fransfer/internal/filetransfer"
	"fransfer/internal/server"

	"github.com/sirupsen/logrus"
)

var (
	//go:embed ..\..\configs\keys\rsa
	privateKey []byte
	//go:embed ..\..\configs\keys\rsa.pub
	publicKey []byte
)

func main() {
	address := flag.String("address", "localhost:50051", "grpc service address")
	outputDir := flag.String("output", "", "directory where files to be saved")
	debug := flag.Bool("debug", true, "debug mode")
	flag.Parse()

	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	jwt, err := auth.NewJWT(privateKey, publicKey)
	if err != nil {
		logrus.Fatal(err)
	}

	s := server.New(*address, jwt)
	s.Register(filetransfer.NewService(*outputDir))

	logrus.Info("Server up and running...")
	logrus.Fatal(s.Run())
}
