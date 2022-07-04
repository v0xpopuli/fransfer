# fransfer
[![.github/workflows/tests.yml](https://github.com/v0xpopuli/fransfer/actions/workflows/tests.yml/badge.svg?branch=development&event=push)](https://github.com/v0xpopuli/fransfer/actions/workflows/tests.yml)
![Coverage](https://img.shields.io/badge/Coverage-73.6%25-brightgreen)

Simple p2p file transfer that works over gRPC. \
If you want to transfer files to your own machine over web, \
just forward needed port, add firewall rule, or simply use tools that will do all of this instead of you (ngrok etc.)

##  How to run
```
1. Run make genproto for generating protobuf
2. Run make genkeypair for generation ssh keypair for JWT signing
3. Run make build-both-for-${windows or macos for now} 
4. Built binaries located in /build/${os} directory 
```

### Usage
```
Usage of server:
  -address
        grpc service address
  -output
        directory where files to be saved
  -debug
        debug mode (means logging level)
        
```
```
Usage of client:
 -address
        grpc service address
 -file
        path to file you want to transfer (can't be empty)
```