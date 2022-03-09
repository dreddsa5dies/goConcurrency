package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// save struct
type data struct {
	bc int    // the byte counts
	cc int    // the character counts
	nc int    // the newline counts
	wc int    // the word counts
	fa string // filepath
}

// counts - byte, character, newline & word counts
func (d *data) counts() error {
	file, err := ioutil.ReadFile(d.fa)
	if err != nil {
		return err
	}
	d.bc = len(file)
	d.nc = strings.Count(string(file), "\n")
	d.cc = len([]rune(string(file)))
	d.wc = len(strings.Fields(string(file)))

	return nil
}

func (d *data) print() {
	switch {
	case *c:
		fmt.Fprintln(os.Stdout, d.bc, d.fa)
	case *m:
		fmt.Fprintln(os.Stdout, d.cc, d.fa)
	case *l:
		fmt.Fprintln(os.Stdout, d.nc, d.fa)
	case *w:
		fmt.Fprintln(os.Stdout, d.wc, d.fa)
	default:
		fmt.Fprintln(os.Stdout, d.nc, d.wc, d.bc, d.fa)
	}
}

var (
	c, m, l, w *bool
)

func init() {
	c = flag.Bool("c", false, "print the byte counts")
	m = flag.Bool("m", false, "print the character counts")
	l = flag.Bool("l", false, "print the newline counts")
	w = flag.Bool("w", false, "print the word counts")
	flag.Parse()
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println(os.Args[0], `- print newline, word, and byte counts for each FILE, and a total line if more than one FILE is specified. A word is a non-zero-length sequence of characters delimited by white space.
Flags:`)
		flag.PrintDefaults()
		os.Exit(1)
	}

	files := os.Args[1:]

	dChan := make(chan data, 2)
	counter := len(files)

	for i := range files {
		if _, err := os.Stat(files[i]); err == nil {
			var tmp data
			tmp.fa = files[i]
			err = tmp.counts()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			select {
			case dChan <- tmp:
			default:
				fmt.Println("Not enough space for buffer", files[i])
			}
		} else {
			fmt.Fprintln(os.Stdout, "No such file")
		}
	}

	for i := 0; i < counter; i++ {
		select {
		case wc := <-dChan:
			wc.print()
		default:
			fmt.Fprintln(os.Stdout, "Nothing more to be done!", files[i])
		}
	}
}
