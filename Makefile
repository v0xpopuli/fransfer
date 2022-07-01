genmocks:
	go install github.com/matryer/moq@latest
	moq -out internal/generated/file_transfer_client_mock.go internal/generated FileTransferClient
	moq -out internal/generated/file_transfer_send_client_mock.go internal/generated FileTransfer_SendClient

tests: genmocks
	go test -v ./... -covermode=count -coverprofile=coverage.out

genproto:
	docker pull namely/protoc-all
	docker run -v $p:/defs namely/protoc-all -d ./ -l go -o internal

genkeypair:
	mkdir -p "configs/keys"
	openssl genrsa -out configs/keys/rsa 4096
	openssl rsa -in configs/keys/rsa -pubout -out configs/keys/rsa.pub

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
