#!/usr/bin/env bash

go install github.com/matryer/moq@latest

moq -out internal/generated/file_transfer_client_mock.go internal/generated FileTransferClient
moq -out internal/generated/file_transfer_send_client_mock.go internal/generated FileTransfer_SendClient

moq -out internal/generated/file_transfer_server_mock.go internal/generated FileTransferServer
moq -out internal/generated/file_transfer_send_server_mock.go internal/generated FileTransfer_SendServer