#!/bin/sh
version="1.0.0"
buildDate=$(date "+%Y-%m-%dT%H:%M:%S")
platform="android"
ldflags="-X 'rtk-cross-share/buildConfig.Version=$version' -X 'rtk-cross-share/buildConfig.BuildDate=$buildDate' -X 'rtk-cross-share/buildConfig.Platform=$platform'"

cd client/platform/libp2p_clipboard
gomobile bind -target=android -androidapi 21 -ldflags "$ldflags"
cd -

