#!/bin/bash

mkdir -p ./proto/golang && mkdir -p ./proto/python && cd ./proto/golang
protoc --go_out=. --proto_path=../ --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ../real_estate_service.proto


cd ../python
python3 -m grpc.tools.protoc --proto_path=../ --python_out=. --grpc_python_out=. ../real_estate_service.proto
