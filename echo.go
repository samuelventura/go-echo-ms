package main

import (
	"io"
	"net"

	"github.com/samuelventura/go-state"
	"github.com/samuelventura/go-tree"
)

func echo(node tree.Node) error {
	endpoint := node.GetValue("endpoint").(string)
	log := node.GetValue("log").(*state.Log)
	listen, err := net.Listen("tcp", endpoint)
	if err != nil {
		return err
	}
	node.AddCloser("listen", listen.Close)
	port := listen.Addr().(*net.TCPAddr).Port
	log.Info("port", port)
	node.SetValue("port", port)
	node.AddProcess("listen", func() {
		id := NewId("echo")
		for {
			conn, err := listen.Accept()
			if err != nil {
				log.Warn(err)
				return
			}
			addr := conn.RemoteAddr().String()
			cid := id.Next(addr)
			child := node.AddChild(cid)
			if child == nil {
				conn.Close()
				continue
			}
			log.Info("open", addr)
			child.AddCloser("conn", conn.Close)
			child.AddProcess("loop", func() {
				defer log.Info("close", addr)
				handleConnection(child, conn)
			})
		}
	})
	return nil
}

func handleConnection(node tree.Node, conn net.Conn) {
	log := node.GetValue("log").(*state.Log)
	err := keepAlive(conn)
	if err != nil {
		log.Warn(err)
		return
	}
	node.AddProcess("copy1", func() {
		_, err := io.Copy(conn, conn)
		if err != nil {
			log.Warn(err)
		}
	})
	node.AddProcess("copy2", func() {
		_, err := io.Copy(conn, conn)
		if err != nil {
			log.Warn(err)
		}
	})
	node.WaitClosed()
}
