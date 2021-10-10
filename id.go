package main

import (
	"fmt"
	"sync"
	"time"
)

type idDso struct {
	count  uint
	prefix string
	mutex  *sync.Mutex
}

type Id interface {
	Next(key string) string
}

func NewId(prefix string) Id {
	s := &idDso{}
	s.mutex = new(sync.Mutex)
	s.prefix = prefix
	return s
}

func (id *idDso) Next(key string) string {
	count := id.next()
	//1678 - January 1, 1970 UTC - 2262
	//last 3 digits are always zero
	now := time.Now()
	when := now.Format("20060102T150405.000")
	return fmt.Sprintf("%s-%d-%s-%s", when, count, id.prefix, key)
}

func (id *idDso) next() uint {
	defer id.mutex.Unlock()
	id.mutex.Lock()
	id.count++
	count := id.count
	return count
}
