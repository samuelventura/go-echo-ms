#!/bin/bash -x

if [[ "$OSTYPE" == "linux"* ]]; then
    SRC=$HOME/go/bin
    DST=/usr/local/bin
    if [[ -f "$DST/go-echo-ss" ]]; then
        sudo systemctl stop GoEchoMs
        sudo $DST/go-echo-ss -service uninstall
        sleep 3
    fi
    go install
    (cd go-echo-ss; go install)
    sudo cp $SRC/go-echo-ms $DST
    sudo cp $SRC/go-echo-ss $DST
    sudo $DST/go-echo-ss -service install
    sudo systemctl restart GoEchoMs
fi
