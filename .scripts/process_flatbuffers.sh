#!/bin/bash

# Process the flatbuffers files
go run ./cmd/tools/fbsprocessor/main.go -dir ./go/flatdata -lang go
