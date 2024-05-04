#!/usr/bin/env bash

GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o bin/gohuw-amd64.exe .
GOOS=windows GOARCH=386 go build -ldflags "-s -w" -o bin/gohuw-386.exe .
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o bin/gohuw-amd64-darwin .
GOOS=darwin GOARCH=386 go build -ldflags "-s -w" -o bin/gohuw-386-darwin .
GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o bin/gohuw-arm64-darwin .
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/gohuw-amd64-linux .
GOOS=linux GOARCH=386 go build -ldflags "-s -w" -o bin/gohuw-386-linux .

upx --ultra-brute --lzma bin/gohuw-amd64.exe bin/gohuw-386.exe bin/gohuw-amd64-darwin bin/gohuw-386-darwin bin/gohuw-amd64-linux bin/gohuw-386-linux
