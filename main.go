package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
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
type Message struct {
	TimeLocal time.Time
	Time      int64
	Value     int64
	Name      string
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
	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  MyDB,
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a point and add to batch
	tags := map[string]string{"cpu": "cpu-total"}
	fields := map[string]interface{}{
		"idle":   10.1,
		"system": 53.3,
		"user":   46.6,
	}

	pt, err := client.NewPoint("cpu_usage", tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := c.Write(bp); err != nil {
		log.Fatal(err)
	}

	// Close client resources
	if err := c.Close(); err != nil {
		log.Fatal(err)
	}
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
