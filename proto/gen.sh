#!/usr/bin/env bash
protoc --go_out=../src/rpc/ --go_opt=paths=source_relative --go-grpc_out=../src/rpc/ --go-grpc_opt=paths=source_relative *.proto