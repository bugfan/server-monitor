package main

import (
	"log"
	"strings"
	"time"
)

type Reader interface {
	Read(r chan string)
}
type Writer interface {
	Write(w chan string)
}
type LogProcess struct {
	readChan  chan string
	writeChan chan string
	read      Reader
	write     Writer
}
type ReadFrom struct {
	path string
}

func (s *ReadFrom) Read(r chan string) {
	r <- "read chan string"
}

type WriteTo struct {
	path string
}

func (s *WriteTo) Write(w chan string) {
	log.Println(<-w)
}

func (s *LogProcess) Process() {
	data := <-s.readChan
	s.writeChan <- strings.ToUpper(data)
}

func main() {
	r := ReadFrom{
		path: "./a.txt",
	}
	w := WriteTo{
		path: "./a.txt",
	}
	lp := LogProcess{
		readChan:  make(chan string),
		writeChan: make(chan string),
		read:      &r,
		write:     &w,
	}
	go lp.read.Read(lp.readChan)
	go lp.Process()
	go lp.write.Write(lp.writeChan)

	time.Sleep(1e9)
}
