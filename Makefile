genmocks:
	go install github.com/matryer/moq@latest
	moq -out internal/generated/file_transfer_client_mock.go internal/generated FileTransferClient
	moq -out internal/generated/file_transfer_send_client_mock.go internal/generated FileTransfer_SendClient

tests: genmocks
	go test -v ./... -covermode=count -coverprofile=coverage.out

genproto:
ifeq ($(OS), Windows_NT)
	scripts/generate-proto.sh `pwd -W`/protobuf
else
	scripts/generate-proto.sh `pwd`/protobuf
endif

genkeypair:
	scripts/generate-key-pair.sh

build-server-for-windows:
	GOOS=windows GOARCH=386 go build -o ./build/windows/server.exe ./cmd/server

build-client-for-windows:
	GOOS=windows GOARCH=386 go build -o ./build/windows/client.exe ./cmd/client

build-server-for-macos:
	GOOS=darwin GOARCH=amd64 go build -o ./build/macos/server ./cmd/server

build-client-for-macos:
	GOOS=darwin GOARCH=amd64 go build -o ./build/macos/client ./cmd/client

build-both-for-windows: build-server-for-windows build-client-for-windows

build-both-for-macos: build-server-for-macos build-client-for-macos
