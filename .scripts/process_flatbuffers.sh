#!/bin/bash

# Process the flatbuffers files
go run github.com/arisu-archive/bluearchive-fbs-generator@latest -i ./.schema/ -o ./go/flatdata -p flatdata
go run ./cmd/tools/fbsprocessor/main.go -dir ./go/flatdata -lang go
go mod tidy
