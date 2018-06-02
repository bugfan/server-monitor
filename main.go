package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type Reader interface {
	Read(r chan []byte)
}
type Writer interface {
	Write(w chan string)
}
type LogProcess struct {
	readChan  chan []byte
	writeChan chan string
	read      Reader
	write     Writer
}
type ReadFrom struct {
	path string
}

type WriteTo struct {
	path string
}

func (s *ReadFrom) Read(r chan []byte) {
	file, err := os.Open(s.path)
	if err != nil {
		log.Fatal(err)
	}
	file.Seek(0, 2)
	rd := bufio.NewReader(file)
	for {
		bts, err := rd.ReadBytes('\n')
		if err == io.EOF {
			time.Sleep(500 * time.Millisecond)
			continue
		} else if err != nil {
			log.Panic("ReadBytes error:", err)
		}
		r <- bts[:len(bts)-1]
	}
}

func (s *WriteTo) Write(w chan string) {
	for data := range w {
		fmt.Println(data)
	}
}

func (s *LogProcess) Process() {
	// 处理数据
	for data := range s.readChan {
		s.writeChan <- strings.ToUpper(string(data))
	}
}

func main() {
	r := ReadFrom{
		path: "./a.txt",
	}
	w := WriteTo{
		path: "./b.txt",
	}
	lp := LogProcess{
		readChan:  make(chan []byte),
		writeChan: make(chan string),
		read:      &r,
		write:     &w,
	}
	go lp.read.Read(lp.readChan)
	go lp.Process()
	go lp.write.Write(lp.writeChan)

	<-(chan bool)(nil)
}
