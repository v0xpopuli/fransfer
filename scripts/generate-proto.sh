#!/usr/bin/env bash

docker pull namely/protoc-all

docker run -v $1:/defs namely/protoc-all -d ./ -l go -o ./

mv `pwd`/protobuf/generated/* `pwd`/internal/generated/
rm -rf `pwd`/protobuf/generated
