package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	"github.com/samuelventura/go-tree"
)

func main() {
	os.Setenv("GOTRACEBACK", "all")
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetOutput(os.Stdout)

	ctrlc := make(chan os.Signal, 1)
	signal.Notify(ctrlc, os.Interrupt)

	log.Println("starting...")
	log.Println("pid", os.Getpid())
	defer log.Println("exit")
	node := root()
	defer node.WaitDisposed()
	defer node.Close()
	err := run(node)
	if err != nil {
		log.Fatal(err)
	}

	stdin := make(chan interface{})
	go func() {
		defer close(stdin)
		ioutil.ReadAll(os.Stdin)
	}()
	select {
	case <-node.Closed():
		log.Println("root closed")
	case <-ctrlc:
		log.Println("ctrlc interrupt")
	case <-stdin:
		log.Println("stdin closed")
	}
}

func root() tree.Node {
	node := tree.NewRoot(nil)
	node.SetValue("endpoint", getenv("ECHO_ENDPOINT", "127.0.0.1:31653"))
	node.SetValue("errors", getenv("ECHO_ERRORS", "false") == "true")
	return node
}

func run(node tree.Node) error {
	return echo(node)
}
