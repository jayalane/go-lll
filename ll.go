// -*- tab-width:2 -*-

package lll

import (
	"github.com/lestrrat-go/file-rotatelogs"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"os/user"
	"strings"
	"sync/atomic"
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
	log    *log.Logger
	N      uint64 // should be a map but whatevs
	level  int32
}

// Init is called once if you want a non-default
// file name

var initOnceDone int64
var theWriter io.Writer
var theLogPath string = "."

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
		logPathTemplate = theLogPath + "/" + item + ".log.%Y%m%d"
	}
	if theWriter == nil {
		// init rotating logs
		t, err := rotatelogs.New(
			logPathTemplate,
			rotatelogs.WithRotationSize(8*1024*1024),
		)
		if err != nil {
			log.Panic("Can't open rotating logs")
		}
		theWriter = t
	}
	log.SetOutput(theWriter)
}

// SetLogPath sets the path for the template for non-root use.  It must be called before any Init()
func SetLogPath(logPath string) {
	if atomic.LoadInt64(&initOnceDone) == 1 {
		panic("SetLogPath called after Init!")
	}
	theLogPath = logPath
}

// SetWriter takes a writer and sets it up so the logger when inited
// will use this writer level
func SetWriter(writer io.Writer) {
	theWriter = writer
}

// Init takes a module name and a level string and returns a logger
func Init(modName string, level string) Lll {
	initOnce()
	if len(modName) > 50 {
		log.Panic("Init lll called with giant module name", modName)
	}
	l := log.New(theWriter, modName, 0)
	l.SetPrefix(modName + " ")
	l.SetOutput(theWriter)
	l.SetFlags(log.Ldate + log.Ltime + log.Lmsgprefix + log.Lmicroseconds)
	res := Lll{module: modName, log: l, level: all}
	res.SetLevel(level)
	return res
}

// SetLevel global needed for backward compatibility
func SetLevel(l *Lll, level string) {
	l.SetLevel(level)
}

// GetLevel needed for go-globals tests
func GetLevel(l *Lll) int {
	return int(atomic.LoadInt32(&l.level))
}

// GetLevel returns the level for this specific logger instance
func (ll *Lll) GetLevel() int {
	return int(atomic.LoadInt32(&ll.level))
}

// SetLevel takes a low level logger and a level string and resets the
// log level
func (ll *Lll) SetLevel(level string) {
	var theLev int32
	if level == "network" {
		theLev = network
	} else if level == "none" {
		theLev = none
	} else if level == "state" {
		theLev = state
	} else {
		theLev = all
	}
	atomic.StoreInt32(&ll.level, theLev)
}

// Ln is Log Network - most volumunous
func (ll Lll) Ln(ls ...interface{}) {
	if ll.level > network {
		return
	}
	ll.log.Println(ls...)
}

// Ls is Log State - TCP reads/writes (but not what), accept/close
func (ll Lll) Ls(ls ...interface{}) {
	if ll.level > state {
		return
	}
	ll.log.Println(ls...)
}

// La is Log Always - Listens, serious errors, etc.
func (ll Lll) La(ls ...interface{}) {
	if ll.level > all {
		return
	}
	ll.log.Println(ls...)
}

// Ll is Log Lots ays - starts out logging all then 1/N then every 60,000th
func (ll *Lll) Ll(ls ...interface{}) {
	if ll.level > all {
		return
	}
	numLoggeds := atomic.AddUint64(&ll.N, 1)
	atomic.StoreUint64(&ll.N, numLoggeds)
	if numLoggeds%50000 == 0 ||
		rand.Float64() < 1.0/float64(math.Pow(math.Log(1+float64(numLoggeds)), 2)) {
		ll.log.Println(ls...)
	}
}
