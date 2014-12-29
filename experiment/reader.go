package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var data = `1415926
5358979
582097
`

type loggingReader struct {
	r io.Reader
}

func (r loggingReader) Read(buf []byte) (int, error) {
	n, err := r.r.Read(buf)
	if err != nil {
		fmt.Printf("read -> %q, error %v\n", buf[0:n], err)
	} else {
		fmt.Printf("read -> %q\n", buf[0:n])
	}
	return n, err
}

func main() {

	r := loggingReader{bufio.NewReader(os.Stdin)}
	readNumbers(r, fscanf)
	// fmt.Printf("-- with Fscanln:\n")
	// readNumbers(loggingReader{strings.NewReader(data)}, fscanln)
	// fmt.Printf("\n-- with Fscanf:\n")
	// readNumbers(loggingReader{strings.NewReader(data)}, fscanf)
}

func fscanf(r io.Reader, ip *int) (int, error) {
	return fmt.Fscanf(r, "%d\n", ip)
}

func fscanln(r io.Reader, ip *int) (int, error) {
	return fmt.Fscanln(r, ip)
}

func readNumbers(r io.Reader, scan func(r io.Reader, ip *int) (int, error)) {
	var i int
	for {
		_, err := scan(r, &i)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Printf("You have entered %d.\n", i)
	}
}
