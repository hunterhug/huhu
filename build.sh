#!/usr/bin/env bash

mkdir ./bin
mkdir ./bin/data
cp main.go ./bin/
echo "" > ./bin/cookie.txt
cd ./bin


GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -x -o zhihu_linux_amd64 main.go
GOOS=linux GOARCH=386 go build -ldflags "-s -w" -x -o zhihu_linux_386 main.go
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -x -o zhihu_windows_amd64.exe main.go
GOOS=windows GOARCH=386 go build -ldflags "-s -w" -x -o zhihu_windows_386.exe main.go
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -x -o zhihu_mac main.go

cd -
rm -rf ./bin/main.go
