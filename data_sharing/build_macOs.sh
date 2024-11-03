#!/bin/sh
CLIENT_FILES="./main/client_main.go"

version=""
buildDate=$(date "+%Y-%m-%dT%H:%M:%S")
ldflags="-X 'rtk-cross-share/buildConfig.Version=$version' -X 'rtk-cross-share/buildConfig.BuildDate=$buildDate'"

cd client
go build -ldflags "$ldflags" -o client_mac $CLIENT_FILES
cd -
