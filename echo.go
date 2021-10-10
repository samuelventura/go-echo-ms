package main

import (
	"io"
	"log"
	"net"

	"github.com/samuelventura/go-tree"
)

func echo(node tree.Node) error {
	endpoint := node.GetValue("endpoint").(string)
	listen, err := net.Listen("tcp", endpoint)
	if err != nil {
		return err
	}
	node.AddCloser("listen", listen.Close)
	port := listen.Addr().(*net.TCPAddr).Port
	log.Println("port", port)
	node.SetValue("port", port)
	node.Go("listen", func() {
		defer node.Close()
		id := NewId("echo-" + listen.Addr().String())
		for {
			conn, err := listen.Accept()
			if err != nil {
				log.Println(err)
				return
			}
			addr := conn.RemoteAddr().String()
			cid := id.Next(addr)
			child := node.AddChild(cid)
			if child == nil {
				conn.Close()
				continue
			}
			child.AddCloser("conn", conn.Close)
			go handleConnection(child, conn)
		}
	})
	return nil
}

func handleConnection(node tree.Node, conn net.Conn) {
	defer node.Close()
	err := keepAlive(conn)
	if err != nil {
		log.Println(err)
		return
	}
	node.Go("copy1", func() {
		defer node.Close()
		_, err := io.Copy(conn, conn)
		if err != nil {
			log.Println(node.Name(), err)
		}
	})
	node.Go("copy2", func() {
		defer node.Close()
		_, err := io.Copy(conn, conn)
		if err != nil {
			log.Println(node.Name(), err)
		}
	})
	<-node.Closed()
}
