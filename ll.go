// -*- tab-width:2 -*-

package lll

import (
	"github.com/lestrrat-go/file-rotatelogs"
	"log"
	"os"
	"os/user"
	"strings"
	"sync/atomic"
	"time"
)

const (
	network = iota
	state
	all
	none
)

// Lll is a low level logger
type Lll struct {
	module string
	level  int
}

// Init is called once if you want a non-default
// file name

var initOnceDone int64

// initOnce nees to be called to get log rotation going
func initOnce() {
	if atomic.LoadInt64(&initOnceDone) == 1 {
		return
	}
	atomic.StoreInt64(&initOnceDone, 1) // defer?  But rather miss some than
	// have this run twice
	binaryFilename, err := os.Executable()
	if err != nil {
		panic(err)
	}
	p := strings.Split(binaryFilename, "/")
	item := p[len(p)-1]
	logPathTemplate := "/var/log/" + item + ".log.%Y%m%d"
	u, err := user.Current()
	if err != nil {
		log.Panic("Can't check user id")
	}
	if u.Uid != "0" {
		logPathTemplate = "./proxy.log.%Y%m%d"
	}
	// init rotating logs
	r1, err := rotatelogs.New(
		logPathTemplate,
		rotatelogs.WithMaxAge(time.Hour*168),
	)
	if err != nil {
		log.Panic("Can't open rotating logs")
	}
	log.SetOutput(r1)
}

// SetLevel takes a low level logger and a level string and resets the log
// level
func SetLevel(res *Lll, level string) {
	var theLev int
	if level == "network" {
		theLev = network
	} else if level == "none" {
		theLev = none
	} else if level == "state" {
		theLev = state
	} else {
		theLev = all
	}
	res.level = theLev
}

// Init takes a module name and a level string and returns a logger
func Init(modName string, level string) Lll {
	initOnce()
	if len(modName) > 50 {
		log.Panic("Init lll called with giant module name", modName)
	}
	res := Lll{module: modName, level: all}
	SetLevel(&res, level)
	return res
}

// Ln is Log Network - most volumunous
func (ll Lll) Ln(ls ...interface{}) {
	if ll.level > network {
		return
	}
	log.Println(ls...)
}

// Ls is Log State - TCP reads/writes (but not what), accept/close
func (ll Lll) Ls(ls ...interface{}) {
	if ll.level > state {
		return
	}
	log.Println(ls...)
}

// La is Log Always - Listens, serious errors, etc.
func (ll Lll) La(ls ...interface{}) {
	if ll.level > all {
		return
	}
	log.Println(ls...)
}
