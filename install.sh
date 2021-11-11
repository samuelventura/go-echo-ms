#!/bin/bash -x

if [[ "$OSTYPE" == "linux"* ]]; then
    BIN=go-echo-ms
    SRC=$HOME/go/bin
    DST=/usr/local/bin
    go install
    DAEMON=echo
    curl -X POST http://127.0.0.1:31600/api/daemon/stop/echo
    curl -X POST http://127.0.0.1:31600/api/daemon/uninstall/echo
    sudo cp $SRC/$BIN $DST
    curl -X POST http://127.0.0.1:31600/api/daemon/install/echo?path=$DST/$BIN
    curl -X POST http://127.0.0.1:31600/api/daemon/enable/echo
    curl -X POST http://127.0.0.1:31600/api/daemon/start/echo
    curl -X GET http://127.0.0.1:31600/api/daemon/info/echo
    curl -X GET http://127.0.0.1:31600/api/daemon/env/echo
fi
