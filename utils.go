package main

import (
	"net"
	"os"
	"strings"
	"time"

	"github.com/felixge/tcpkeepalive"
)

func keepAlive(conn net.Conn) error {
	return tcpkeepalive.SetKeepAlive(
		conn, 5*time.Second, 3, 1*time.Second)
}

func getenv(name string, defval string) string {
	value := os.Getenv(name)
	if len(strings.TrimSpace(value)) > 0 {
		return value
	}
	return defval
}
