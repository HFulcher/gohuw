#!/usr/bin/env bash

GOOS=windows GOARCH=amd64 go build -o bin/gohuw-amd64.exe .
GOOS=windows GOARCH=386 go build -o bin/gohuw-386.exe .
GOOS=darwin GOARCH=amd64 go build -o bin/gohuw-amd64-darwin .
GOOS=darwin GOARCH=386 go build -o bin/gohuw-386-darwin .
GOOS=linux GOARCH=amd64 go build -o bin/gohuw-amd64-linux .
GOOS=linux GOARCH=386 go build -o bin/gohuw-386-linux .