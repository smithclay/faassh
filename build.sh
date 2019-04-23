#!/bin/zsh

mkdir -p layer/bin

GOOS=linux GOARCH=amd64 go build -o layer/bin/faassh main.go
