#!/usr/bin/env bash

docker pull namely/protoc-all

cp -r protobuf/* scripts/

docker run -v $(pwd)/protobuf:/defs namely/protoc-all -d ./ -l go -o ./

sudo rm scripts/filetransfer.proto
sudo mv protobuf/generated/* internal/generated
sudo rm -rf protobuf/generated