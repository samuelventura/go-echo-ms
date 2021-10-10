package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/kardianos/service"
)

func logfp(fp string) string {
	dir := os.Getenv("ECHO_LOGS")
	if len(dir) > 0 {
		//filepath handles Windows/Unix separators
		return filepath.Join(dir, filepath.Base(fp))
	}
	return fp
}

func daemon(log service.Logger, sibling string, exit chan bool) chan bool {
	done := make(chan bool)
	path := relativeSibling(sibling)
	outp := changeExtension(path, ".out.log")
	ff := os.O_APPEND | os.O_WRONLY | os.O_CREATE
	outf, err := os.OpenFile(logfp(outp), ff, 0644)
	panicIfError(err)
	errp := changeExtension(path, ".err.log")
	errf, err := os.OpenFile(logfp(errp), ff, 0644)
	panicIfError(err)
	go func() {
		defer log.Infof("exited %s", path)
		defer traceRecover(log)
		defer close(done)
		defer outf.Close()
		defer errf.Close()
		run := func() {
			defer traceRecover(log)
			cmd := exec.Command(path)
			cmd.Env = os.Environ()
			cmd.Stdout = outf
			cmd.Stderr = errf
			sin, err := cmd.StdinPipe()
			panicIfError(err)
			defer sin.Close()
			err = cmd.Start()
			panicIfError(err)
			go func() {
				defer traceRecover(log)
				defer sin.Close()
				select {
				case <-exit:
				case <-done:
				}
			}()
			err = cmd.Wait()
			panicIfError(err)
		}
		count := 0
		for {
			if count > 0 {
				time.Sleep(millis(2000))
			}
			log.Infof("%d %s", count, path)
			run()
			count++
			select {
			case <-exit:
				return
			default:
				continue
			}
		}
	}()
	return done
}
