package main

import (
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/felixge/tcpkeepalive"
)

func keepAlive(conn net.Conn) {
	err := tcpkeepalive.SetKeepAlive(
		conn, 5*time.Second, 3, 1*time.Second)
	if err != nil {
		log.Fatal(err)
	}
}

func getenv(name string, defval string) string {
	value := os.Getenv(name)
	trimmed := strings.TrimSpace(value)
	if len(trimmed) > 0 {
		log.Println(name, value)
		return value
	}
	log.Println(name, defval)
	return defval
}
