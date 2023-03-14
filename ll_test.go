// -*- tab-width: 2 -*-

package lll

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"
)

// Buffer type from stackoverflow for safe concurrency
type Buffer struct {
	b bytes.Buffer
	m sync.Mutex
}

func (b *Buffer) Read(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Read(p)
}
func (b *Buffer) Write(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Write(p)
}
func (b *Buffer) String() string {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.String()
}

func TestLa(t *testing.T) {

	var ml *Lll
	var buffer = new(Buffer)
	var modName = "TEST"
	var msgString = "hi"
	var numLogs = 1000

	SetWriter(buffer)
	ml = Init("TEST", "debug")
	for i := 0; i < numLogs; i++ {
		ml.La("hi" + fmt.Sprint(i))
	}
	time.Sleep(1100 * time.Millisecond)
	scanner := bufio.NewScanner(buffer)
	i := 0
	for scanner.Scan() {
		l := scanner.Text()
		sections := strings.Split(l, " ")
		if len(sections) != 4 {
			t.Fatal("Wrong format", l, len(sections))
			return
		}
		if modName != sections[2] {
			t.Fatal("Wrong module", l, sections[2])
			return
		}
		if !(msgString+fmt.Sprint(i) == sections[3]) {
			t.Fatal("Wrong log msg", l, sections[3])
			return
		}
		i++
	}

	if err := scanner.Err(); err != nil {
		t.Fatal(err)
		return
	}
	if i != numLogs {
		t.Fatal("Not enough logs got", i, "wanted", numLogs)
	}
	// all good
}

func TestLl(t *testing.T) {

	var ml *Lll
	var buffer = new(Buffer)
	var modName = "TEST"
	var msgString = "yo_this_is_a_longer_string_maybe_some_memory_thing_z"
	var numLogs = 10000
	SetWriter(buffer)
	ml = Init("TEST", "debug")
	ml.SetLevel("network")
	for i := 0; i < numLogs; i++ {
		go func(j int) {
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			ml.Ll(msgString + fmt.Sprint(j))
			ml.Ln(msgString + fmt.Sprint(j) + "whatevs")
		}(i)
	}
	ml.SetLevel("all")
	for i := 0; i < numLogs; i++ {
		go func(j int) {
			ml.Ll(msgString + fmt.Sprint(j) + "2")
			ml.Ln(msgString + fmt.Sprint(j) + "whatevs2")
		}(i)
	}
	time.Sleep(1100 * time.Millisecond)
	scanner := bufio.NewScanner(buffer)
	i := 0
	for scanner.Scan() {
		l := scanner.Text()
		fmt.Println("Got a line:", l)
		sections := strings.Split(l, " ")
		if len(sections) != 4 {
			t.Fatal("Wrong format", l, len(sections))
			return
		}
		if modName != sections[2] {
			t.Fatal("Wrong module", l, sections[2])
			return
		}
		if !(msgString == sections[3][0:len(msgString)]) {
			t.Fatal("Wrong log msg", l, sections[3])
			return
		}
		i++
	}

	if err := scanner.Err(); err != nil {
		t.Fatal(err)
		return
	}

	if i > numLogs/2 {
		t.Fatal("Too many logs got", i, "wanted a small fraction of", numLogs)
	}
	t.Log("Got", i)
	// all good
}
