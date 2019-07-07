package main

import (
	"net"
	"log"
	"encoding/gob"
	"math/rand"
  	"time"
  	"strings"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
				"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func make_random_string(length int, charset string) string {
  	
  	b := make([]byte, length)
  	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func generate_random_string(length int) string {
	return make_random_string(length, charset)
}


func handle_connection(c net.Conn) {
    
   // log.Println("received a request")

   	var buffer string
    dec := gob.NewDecoder(c)
    err := dec.Decode(&buffer)

    if err != nil {
    	log.Fatal(err)
    	log.Fatal("Fail to Decode")
    }

    info := strings.Split(buffer,"|")
    log.Println(info)
    log.Println(info[0], " entered in critical region with timestamp ", info[1])

    enc := gob.NewEncoder(c)
    err = enc.Encode(generate_random_string(10))
    
    if err != nil {
    	log.Fatal("Fail to Encode")
    }
}

func run_server() {
	
	l, err := net.Listen("tcp", "localhost:8081")
	
	if err != nil {
		log.Fatal("Fail to connect")
	}

	log.Println("Server of Critical Region is running")

	defer l.Close()

   
	for {

	    c, err := l.Accept()
		if err != nil {
			log.Fatal("Fail to connect")
		}

		go handle_connection(c)
	}
}

func main() {
	c := make(chan int)
	go run_server()
	<-c
}