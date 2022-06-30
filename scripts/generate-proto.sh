#!/usr/bin/env bash

docker pull namely/protoc-all

docker run -v $1:/defs namely/protoc-all -d ./ -l go -o ./

rm -rf `pwd`/internal/generated
mv `pwd`/protobuf/generated `pwd`/internal