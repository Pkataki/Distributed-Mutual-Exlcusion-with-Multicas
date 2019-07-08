package main

import (
	"strconv"
	"log"
)

func main() {

	const NUMBER_PROCESSES = 5
	const INITIAL_PORT = 8090

	var p[NUMBER_PROCESSES] process

	a := make(chan bool)

	for i := 0; i < NUMBER_PROCESSES; i++ {
		port := INITIAL_PORT + i
		go func(i int) {
			s := strconv.Itoa(port)
			log.Println(port, "localhost:" + s )
			p[i].startProcess("localhost:" + s, i)
			a <- true
		}(i)
	}

	for i := 0; i < NUMBER_PROCESSES; i++ {
		<-a
	}

	for i := 0; i < NUMBER_PROCESSES; i++ {
		go p[i].runProcess()
	}
	waiter := make(chan bool)
	<-waiter
}