package main

import (
	"github.com/samuelventura/go-state"
	"github.com/samuelventura/go-tree"
)

func main() {
	run(func(root tree.Node) {
		enode := root.AddChild("echo")
		snode := root.AddChild("state")
		if enode == nil || snode == nil {
			return //root already closed
		}

		log := root.GetValue("log").(*state.Log)
		endpoint := getenv("ECHO_ENDPOINT", "127.0.0.1:31653")
		log.Info("endpoint", endpoint)

		enode.AddAction("root", root.Close)
		enode.SetValue("endpoint", endpoint)
		err := echo(enode) //async
		if err != nil {
			root.Close()
			log.Warn(err)
			return
		}
		path := state.SingletonPath("/tmp")
		log.Info("path", path)
		mux := state.NewMux()
		state.AddPProfHandlers(mux)
		state.AddNodeHandlers(mux, root)
		state.AddNodeHandlers(mux, enode)
		state.AddNodeHandlers(mux, snode)
		state.AddEnvironHandlers(mux)
		snode.AddAction("root", root.Close)
		snode.SetValue("mux", mux)
		snode.SetValue("path", path)
		err = state.Serve(snode) //async
		if err != nil {
			root.Close()
			log.Warn(err)
			return
		}
	})
}
