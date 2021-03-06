#!/bin/bash
protoc  -I=proto \
  --go_out=proto/grpc/v1 \
  --go_opt=paths=source_relative \
  --go-grpc_out=proto/grpc/v1 \
  --go-grpc_opt=paths=source_relative \
  proto/manager.proto