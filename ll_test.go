// -*- tab-width: 2 -*-

package lll

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestLa(t *testing.T) {

	var ml Lll
	var buffer = new(bytes.Buffer)
	var modName = "TEST"
	var msgString = "hi"

	SetWriter(buffer)
	ml = Init("TEST", "debug")
	for i := 0; i < 1000; i++ {
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
	// all good
}

func TestLl(t *testing.T) {

	var ml Lll
	var buffer = new(bytes.Buffer)
	var modName = "TEST"
	var msgString = "hi"

	SetWriter(buffer)
	ml = Init("TEST", "debug")
	for i := 0; i < 1000; i++ {
		ml.Ll("hi" + fmt.Sprint(i))
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
	// all good
}
