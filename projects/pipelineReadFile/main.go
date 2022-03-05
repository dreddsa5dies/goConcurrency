package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	filepath := flag.String("f", "/home/dreddsa/go/src/github.com/dreddsa5dies/Mastering-Go-Second-Edition-master/", "Filepath")
	phrase := flag.String("p", "main", "Phrase")
	flag.Parse()

	A := make(chan string)
	B := make(chan int)

	go first(*filepath, A)
	go second(B, A, *phrase)
	third(B)
}

func read(file, phrase string) int {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)

	count := 0
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(err)
			break
		}
		count += strings.Count(line, phrase)
	}

	return count
}

func first(source string, out chan<- string) {
	err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		out <- path
		return nil
	})
	if err != nil {
		log.Fatalf("walk error [%v]\n", err)
	}
	close(out)
}

func second(out chan<- int, in <-chan string, phrase string) {
	for x := range in {
		out <- read(x, phrase)
	}
	fmt.Println()
	close(out)
}

func third(in <-chan int) {
	var sum int
	sum = 0
	for x2 := range in {
		sum = sum + x2
	}
	fmt.Printf("The sum of phrase %d.\n", sum)
}
