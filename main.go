package main

import (
	"bufio"
	"io"
	"log"
	"math/rand"
	"os"
	"server-monitor/influxdb"
	"strings"
	"time"
	// client "github.com/influxdata/influxdb/client/v2"
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

type People struct {
	Id   int64
	Name string
	Age  string
}

func (s *WriteTo) Write(w chan string) {
	for data := range w {
		// fmt.Println(data)
		c := influxdb.Cli{
			Addr:      "http://localhost:8086",
			Username:  "admin",
			Password:  "123456",
			MyDB:      "mydb",
			Precision: "ms",
		}
		c.InitHttp()
		err := c.WriteDB("table1", map[string]string{
			"name": "zxy",
		}, map[string]interface{}{"time": time.Now().Unix(), "value": rand.Intn(100), "str": data})
		if err != nil {
			log.Println(err)
		}
		// log.Println(c.QueryDB("select *from table1"))
		c.Session.Close()
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
		path: "./log.txt",
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
