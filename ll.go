// -*- tab-width:2 -*-

package lll

import (
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"os/user"
	"strings"
	"sync/atomic"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

const (
	network = iota
	state
	all
	none
)

const (
	eightMB = 8 * 1024 * 1024
	// MaxModLen is the maximum length for the module name.
	MaxModLen = 50
	two       = 2
)

// Lll is a low level logger.
type Lll struct {
	module string
	n      uint64 // should be a map but whatevs
	log    *log.Logger
	level  int64
}

// Init is called once if you want a non-default
// file name

var (
	initOnceDone int64
	theWriter    io.Writer
	theLogPath   = "."
)

// initOnce nees to be called to get log rotation going.
func initOnce() {
	swapped := atomic.CompareAndSwapInt64(&initOnceDone, 0, 1)
	if !swapped {
		return
	}

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
			rotatelogs.WithRotationSize(eightMB),
		)
		if err != nil {
			log.Panic("Can't open rotating logs")
		}

		theWriter = t
	}

	log.SetOutput(theWriter)
}

// SetLogPath sets the path for the template for non-root use.  It must be called before any Init().
func SetLogPath(logPath string) {
	if atomic.LoadInt64(&initOnceDone) == 1 {
		panic("SetLogPath called after Init!")
	}

	theLogPath = logPath
}

// SetWriter takes a writer and sets it up so the logger when inited
// will use this writer level.
func SetWriter(writer io.Writer) {
	theWriter = writer
}

// Init takes a module name and a level string and returns a logger.
func Init(modName string, level string) *Lll {
	initOnce()

	if len(modName) > MaxModLen {
		log.Panic("Init lll called with giant module name", modName)
	}

	l := log.New(theWriter, modName, 0)

	l.SetPrefix(modName + " ")
	l.SetOutput(theWriter)
	l.SetFlags(log.Ldate + log.Ltime + log.Lmsgprefix + log.Lmicroseconds)

	res := Lll{module: modName, log: l, level: all}
	res.SetLevel(level)

	return &res
}

// SetLevel global needed for backward compatibility.
func SetLevel(l *Lll, level string) {
	l.SetLevel(level)
}

// GetLevel needed for go-globals tests.
func GetLevel(l *Lll) int {
	return int(atomic.LoadInt64(&l.level))
}

// GetLevel returns the level for this specific logger instance.
func (ll *Lll) GetLevel() int {
	return int(atomic.LoadInt64(&ll.level))
}

// SetLevel takes a low level logger and a level string and resets the
// log level.
func (ll *Lll) SetLevel(level string) {
	var theLev int64

	switch level {
	case "network":
		theLev = network
	case "none":
		theLev = none
	case "state":
		theLev = state
	default:
		theLev = all
	}

	atomic.StoreInt64(&ll.level, theLev)
	fmt.Println("Set level to", atomic.LoadInt64(&ll.level), theLev)
}

// Ln is Log Network - most volumunous.
func (ll *Lll) Ln(ls ...interface{}) {
	if atomic.LoadInt64(&ll.level) > network {
		return
	}

	ll.log.Println(ls...)
}

// Ls is Log State - TCP reads/writes (but not what), accept/close.
func (ll *Lll) Ls(ls ...interface{}) {
	if atomic.LoadInt64(&ll.level) > state {
		return
	}

	ll.log.Println(ls...)
}

// La is Log Always - Listens, serious errors, etc.
func (ll *Lll) La(ls ...interface{}) {
	if atomic.LoadInt64(&ll.level) > all {
		return
	}

	ll.log.Println(ls...)
}

// Ll is Log Lots ays - starts out logging all then 1/N then every 60,000th.
func (ll *Lll) Ll(ls ...interface{}) {
	if atomic.LoadInt64(&ll.level) > all {
		return
	}

	numLoggeds := atomic.AddUint64(&ll.n, 1)

	if numLoggeds%50000 == 0 ||
		rand.Float64() < 1.0/float64(math.Pow(math.Log(1+float64(numLoggeds)), two)) { //nolint:gosec
		ll.log.Println(ls...)
	}
}
