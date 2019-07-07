package main

// import (
// 	"log"
// )

func main() {
	var p[3] process

	a := make(chan bool)
	// go func() {
	// 	p[0].startProcess("localhost:8090", 8090)
	// 	a <- true
	// }()

	go func() {
		p[1].startProcess("localhost:8091", 8091)
		a <- true
	}()

	// go func() {
	// 	p[2] = newProcess("localhost:8092", 2)
	// 	a <- true
	// }()

	// <-a
	<-a


	// for i := 0; i < 1; i++ {
	// //	log.Println("process ", p[i].id)
	// 	go p[i].runProcess()
	// }
	waiter := make(chan bool)
	<-waiter
}