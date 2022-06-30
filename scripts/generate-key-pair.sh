#!/usr/bin/env bash

mkdir -p "configs/keys"

openssl genrsa -out configs/keys/rsa 4096
openssl rsa -in configs/keys/rsa -pubout -out configs/keys/rsa.pub