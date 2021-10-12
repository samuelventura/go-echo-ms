package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	"github.com/samuelventura/go-state"
	"github.com/samuelventura/go-tree"
)

func run(launch func(tree.Node)) {
	os.Setenv("GOTRACEBACK", "all")
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetOutput(os.Stdout)

	slog := &state.Log{}
	slog.Info = func(args ...interface{}) {
		log.Println("info", args)
	}
	slog.Warn = func(args ...interface{}) {
		log.Println("warn", args)
	}
	slog.Recover = func(ss string, args ...interface{}) {
		log.Println("recover", args, ss)
	}

	slog.Info("start", os.Getpid())
	defer slog.Info("exit")

	ctrlc := make(chan os.Signal, 1)
	signal.Notify(ctrlc, os.Interrupt)

	root := tree.NewRoot("root", &slog.Log)
	root.SetValue("log", slog)
	defer root.WaitDisposed()
	defer slog.Info("disposing")
	defer root.Recover()
	//async launcher must close root on error
	//and cleanup on root closed channel.
	go launch(root)

	stdin := make(chan interface{})
	go func() {
		defer close(stdin)
		ioutil.ReadAll(os.Stdin)
	}()
	select {
	case <-root.Closed():
	case <-ctrlc:
	case <-stdin:
	}
}
