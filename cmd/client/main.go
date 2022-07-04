package main

import (
	_ "embed"
	"flag"
	"fmt"
	"fransfer/internal/auth"
	"fransfer/internal/filetransfer"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/TwiN/go-color"
	"github.com/sirupsen/logrus"
)

var (
	//go:embed key
	privateKey []byte
)

func main() {
	address := flag.String("address", "localhost:50051", "grpc service address")
	filePath := flag.String("file", "", "path to file you want to transfer")
	flag.Parse()

	if *filePath == "" {
		fmt.Println(color.InRed("File can't be empty"))
		flag.Usage()
		os.Exit(1)
	}

	fileContent, err := ioutil.ReadFile(*filePath)
	if err != nil {
		fmt.Println(color.InRed(fmt.Sprintf("Error occurred during reading file; %+v", err)))
		os.Exit(1)
	}

	jwt, err := auth.NewJWTIssuer(privateKey)
	if err != nil {
		logrus.Fatal(err)
	}

	c, err := filetransfer.NewClient(*address, jwt)
	if err != nil {
		fmt.Println(color.InRed(fmt.Sprintf("Error occurred during client init; %+v", err)))
		os.Exit(1)
	}
	defer c.Close()

	if err := c.Send(*filePath, fileContent); err != nil {
		fmt.Println(color.InRed(fmt.Sprintf("Error occurred during file delivering; %+v", err)))
		os.Exit(1)
	}

	log.Printf("File [%s] transferred to host [%s] successfully", filepath.Base(*filePath), *address)
}
