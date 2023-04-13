#!/usr/bin/env bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o uploadFileForLinux main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o uploadFileForMacos main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o uploadFileForM1 main.go
# 如果出现 invalid char '\r' 警告 可以使用dos2unix工具转换