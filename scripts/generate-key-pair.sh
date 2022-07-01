#!/usr/bin/env bash

openssl genrsa -out cmd/client/key 4096
openssl rsa -in cmd/client/key -pubout -out cmd/server/key.pub