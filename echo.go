package main

import (
	"io"
	"log"
	"net"

	"github.com/samuelventura/go-tree"
)

func echo(node tree.Node) {
	endpoint := node.GetValue("endpoint").(string)
	listen, err := net.Listen("tcp", endpoint)
	if err != nil {
		log.Fatal(err)
	}
	node.AddCloser("listen", listen.Close)
	port := listen.Addr().(*net.TCPAddr).Port
	log.Println("port", port)
	node.SetValue("port", port)
	node.AddProcess("listen", func() {
		id := NewId("echo")
		for {
			conn, err := listen.Accept()
			if err != nil {
				log.Fatal(err)
				return
			}
			addr := conn.RemoteAddr().String()
			cid := id.Next(addr)
			child := node.AddChild(cid)
			log.Println("open", addr)
			child.AddCloser("conn", conn.Close)
			child.AddProcess("loop", func() {
				defer log.Println("close", addr)
				handleConnection(child, conn)
			})
		}
	})
}

func handleConnection(node tree.Node, conn net.Conn) {
	err := keepAlive(conn)
	if err != nil {
		log.Fatal(err)
		return
	}
	node.AddProcess("copy1", func() {
		_, err := io.Copy(conn, conn)
		if err != nil {
			log.Println(err)
		}
	})
	node.AddProcess("copy2", func() {
		_, err := io.Copy(conn, conn)
		if err != nil {
			log.Println(err)
		}
	})
	node.WaitClosed()
}
