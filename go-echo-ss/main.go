package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kardianos/service"
)

var exit chan bool
var single chan bool
var logger service.Logger

type program struct{}

func (p *program) Start(s service.Service) (err error) {
	exit = make(chan bool)
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("recover %v", r)
		}
	}()
	single = daemon(logger, "go-echo-ms", exit)
	return nil
}

func (p *program) Stop(s service.Service) error {
	close(exit)
	select {
	case <-single:
	case <-time.After(3 * time.Second):
	}
	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetOutput(os.Stdout)
	//-service install, uninstall, start, stop, restart
	svcFlag := flag.String("service", "", "Control the system service.")
	flag.Parse()
	svcConfig := &service.Config{
		Name:        "GoEchoMs",
		DisplayName: "GoEchoMs Service",
		Description: "GoEchoMs https://github.com/samuelventura/go-echo-ms",
	}
	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}
	// use = for logger not to be nil
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	//after logger created
	environFromFile(logger)
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
