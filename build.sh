#!/bin/zsh

GOOS=linux GOARCH=amd64 go build -o lambda/faassh main.go
