package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kardianos/service"
)

//different filename same extension
func relativeSibling(sibling string) string {
	exe := executablePath()
	dir := filepath.Dir(exe)
	base := filepath.Base(exe)
	ext := filepath.Ext(base) //includes .
	file := sibling + ext
	return filepath.Join(dir, file)
}

//same file name different extension
func relativeExtension(ext string) string {
	path := executablePath()
	return changeExtension(path, ext)
}

func changeExtension(path string, next string) string {
	ext := filepath.Ext(path) //includes .
	npath := strings.TrimSuffix(path, ext)
	return npath + next
}

func executablePath() string {
	exe, err := os.Executable()
	panicIfError(err)
	return exe
}

func environFromFile(log service.Logger) {
	path := relativeExtension(".config")
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	log.Infof("loading config %s", path)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			log.Warningf("invalid config line %s", line)
			continue
		}
		log.Infof("setting config %s", line)
		os.Setenv(parts[0], parts[1])
	}
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func millis(ms int) time.Duration {
	return time.Duration(ms) * time.Millisecond
}

func traceRecover(log service.Logger) {
	r := recover()
	if r != nil {
		log.Warningf("recover %v", r)
	}
}
