package main

import (
	"net"
	"fmt"
	"log"
	"sync"
	"time"
	"encoding/gob"
)

type server struct {
	sync.Mutex
	queue [] string
}

func (s *server) handle_connection(c net.Conn) {
    
    log.Println("received a request")
   	var buffer string
    dec := gob.NewDecoder(c)
    err := dec.Decode(&buffer)

    if err != nil {
    	log.Fatal(err)
    	log.Fatal("Fail to Decode")
    }

    log.Println("decoded a request")
    log.Println("Received: ", buffer)
    s.Lock()
    s.queue = append(s.queue, buffer)
    s.Unlock()
    
    time.Sleep(3 * time.Second)

    enc := gob.NewEncoder(c)
    err = enc.Encode(s.queue)
    
    if err != nil {
    	log.Fatal("Fail to Encode")
    }
}

func (s *server) run_server() {
	

	// var wg sync.WaitGroup
	// wg.Add(1)

	l, err := net.Listen("tcp", "localhost:8080")
	
	if err != nil {
		log.Fatal("Fail to connect")
	}

	log.Println("Server Running")

	defer l.Close()

	// go func() { // time to run the server
	// 	defer wg.Done()
 //        time.Sleep(100 * time.Second)
 //    }()

    fmt.Println("Ready to receive requests")

	go func () { // 
		for {

		    c, err := l.Accept()

			if err != nil {
				log.Fatal("Fail to connect")
			}

			go s.handle_connection(c)
		}
	}()

	ch := make(chan int)
	<-ch
	//wg.Wait()

}

func main() {
	name_server := server{}
	name_server.run_server()

}