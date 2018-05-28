package main

import (
	"flag"
	"log"
)

func main(){
	var key string
	flag.StringVar(&key,"name","eye","bee")	 //解析命令参数
	flag.Parse()
	log.Println("end:",key)
}
