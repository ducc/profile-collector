#!/bin/bash

rm -rf protos/
protoc -I=. --go_out=plugins=grpc:. *.proto
mv github.com/ducc/profile-collector/protos/*.pb.go .
rm -r github.com/
