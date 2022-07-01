#!/usr/bin/env bash

go install github.com/matryer/moq@latest

moq -out internal/generated/file_transfer_client_mock.go internal/generated FileTransferClient
moq -out internal/generated/file_transfer_send_client_mock.go internal/generated FileTransfer_SendClient