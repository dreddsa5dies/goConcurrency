package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

// save struct
type data struct {
	RWM sync.RWMutex
	bc  int    // the byte counts
	cc  int    // the character counts
	nc  int    // the newline counts
	wc  int    // the word counts
	fa  string // filepath
}

// counts - byte, character, newline & word counts
func (d *data) counts(f string) error {
	d.RWM.Lock()

	d.fa = f
	file, err := ioutil.ReadFile(d.fa)
	if err != nil {
		return err
	}
	d.bc = len(file)
	d.nc = strings.Count(string(file), "\n")
	d.cc = len([]rune(string(file)))
	d.wc = len(strings.Fields(string(file)))

	d.RWM.Unlock()
	return nil
}

func (d *data) print() {
	d.RWM.RLock()
	fmt.Fprintln(os.Stdout, d.nc, d.wc, d.bc, d.fa)
	defer d.RWM.RUnlock()
}

var allCounts data

func (d *data) allCount() {
	allCounts.RWM.Lock()

	allCounts.bc += d.bc
	allCounts.nc += d.nc
	allCounts.cc += d.cc
	allCounts.wc += d.wc

	defer allCounts.RWM.Unlock()
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println(os.Args[0], `- print newline, word, and byte counts for each FILE, and a total line if more than one FILE is specified. A word is a non-zero-length sequence of characters delimited by white space.`)
		os.Exit(1)
	}

	files := os.Args[1:]

	var waitGroup sync.WaitGroup
	d := data{}

	for i := range files {
		if _, err := os.Stat(files[i]); err == nil {
			waitGroup.Add(1)
			fileString := files[i]
			go func(s string) {
				defer waitGroup.Done()
				err = d.counts(fileString)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
				d.allCount()
				d.print()
			}(fileString)
		} else {
			fmt.Fprintln(os.Stdout, "No such file", files[i])
		}
	}

	waitGroup.Add(1)
	defer allCounts.print()
	waitGroup.Done()

	waitGroup.Wait()
}
