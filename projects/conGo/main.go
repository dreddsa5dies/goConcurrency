package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

var readValue = make(chan data)
var writeValue = make(chan data)

// save struct
type data struct {
	bc int    // the byte counts
	cc int    // the character counts
	nc int    // the newline counts
	wc int    // the word counts
	fa string // filepath
}

// set - byte, character, newline & word counts
func set(f string) error {
	newValue := data{}
	newValue.fa = f
	file, err := ioutil.ReadFile(newValue.fa)
	if err != nil {
		return err
	}
	newValue.bc = len(file)
	newValue.nc = strings.Count(string(file), "\n")
	newValue.cc = len([]rune(string(file)))
	newValue.wc = len(strings.Fields(string(file)))
	writeValue <- newValue
	return nil
}

// print - byte, character, newline & word counts
func (d *data) print() {
	fmt.Fprintln(os.Stdout, d.nc, d.wc, d.bc, d.fa)
}

func controlGo() {
	var value data
	for {
		select {
		case newValue := <-writeValue:
			value = newValue
			value.print()
		case readValue <- value:
		}
	}
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println(os.Args[0], `- print newline, word, and byte counts for each FILE, and a total line if more than one FILE is specified. A word is a non-zero-length sequence of characters delimited by white space.`)
		os.Exit(1)
	}

	files := os.Args

	var w sync.WaitGroup
	go controlGo()

	for i := range files {
		if _, err := os.Stat(files[i]); err == nil {
			w.Add(1)
			fileString := files[i]
			go func(s string) {
				defer w.Done()
				err = set(s)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}(fileString)
		} else {
			fmt.Fprintln(os.Stdout, "No such file", files[i])
		}
	}

	w.Wait()
}
