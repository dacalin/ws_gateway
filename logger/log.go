package _logger

import (
	"log"
	"sync"
)

type Log struct {
	debug bool
}

var lock = &sync.Mutex{}
var instance *Log

func Instance() *Log {
	return instance
}

func New(debug bool) *Log {
	lock.Lock()
	defer lock.Unlock()
	if instance == nil {
		instance = &Log{
			debug: debug,
		}
	}
	return instance
}

func (self *Log) Printf(format string, v ...any) {
	if self.debug {
		log.Printf(format, v)
	}
}

func (self *Log) Println(v ...any) {
	if self.debug {
		log.Println(v)
	}
}
