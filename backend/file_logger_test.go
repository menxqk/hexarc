package backend

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestNewFileTransactionLogger(t *testing.T) {
	const filename = "/tmp/new_ftl_temp.txt"
	defer os.Remove(filename)

	tl, err := NewFileTransactionLogger(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer tl.Close()

	_, ok := tl.(*FileTransactionLogger)
	if !ok {
		panic(fmt.Errorf("type mismatch for FileTransactionLogger: %v", reflect.TypeOf(tl)))
	}
}

func TestFileReadEvents(t *testing.T) {
	const filename = "/tmp/read4_ftl_temp.txt"
	defer os.Remove(filename)

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	err = write4FileEvents(file)
	if err != nil {
		t.Fatalf("could not write 4 events: %v", err)
	}

	tl1, err := NewFileTransactionLogger(filename)
	if err != nil {
		t.Fatal(err)
	}
	chEv, chErr := tl1.ReadEvents()
	count, ok := 0, true
	for ok && err == nil {
		select {
		case err, ok = <-chErr:
		case _, ok = <-chEv:
			if ok {
				count++
			}
		}
	}
	if count != 4 {
		t.Errorf("4 events were expected, events read: %d", count)
	}
	tl1.Run()
	tl1.WritePut("quatro", "quarta") // Event 5
	tl1.WriteDelete("quatro")        // Event 6
	tl1.WritePut("cinco", "quinta")  // Event 7
	err = tl1.Close()
	if err != nil {
		t.Error(err)
	}

	tl2, err := NewFileTransactionLogger(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer tl2.Close()
	chEv, chErr = tl2.ReadEvents()
	count, ok = 0, true
	for ok && err == nil {
		select {
		case err, ok = <-chErr:
		case _, ok = <-chEv:
			if ok {
				count++
			}
		}
	}
	if count != 7 {
		t.Errorf("7 events were expected, events read: %d", count)
	}

	if tl2.LastSequence() != 7 {
		t.Errorf("Lastsequence should be 7 not %d", count)
	}
}

func write4FileEvents(f *os.File) error {
	events := []string{
		"1	1	uma	primeira",
		"2	1	duas	segunda",
		"3	1	três	terceira",
		"4	2	uma	",
	}

	for _, e := range events {
		_, err := f.WriteString(e + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
