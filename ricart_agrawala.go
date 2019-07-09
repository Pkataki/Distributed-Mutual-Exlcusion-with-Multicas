package main

import (
	"log"
	"net"
)

func (p *process) waitAllProcessesReplies() {
	<-p.receivedAllReplies
}

func (p *process) replyAllEnqueuedRequests() {
	for p.isMessageQueueEmpty() == false {
		address := p.getMessageQueueTopRequest()
		//log.Println(p.timestamp, " Process ", p.id, " sendind enqueue REQUEST as a REPLY to ", address)
		p.sendMessage(REPLY, address)
	}
}

func (p *process) startListenPort() error {
	listener, err := net.Listen("tcp", p.address)
	if err != nil {
		return err
	}

	go func(listener net.Listener) {
		defer listener.Close()
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println(err)
			}
			go p.handleRequest(conn)
		}
	}(listener)
	return nil
}

func (p *process) handleRequest(connection net.Conn) {
	defer connection.Close()
	var msg message

	if err := msg.receiveAndDecodeMessage(connection); err != nil {
		log.Println("error on receiveAndDecodeMessage")
		log.Fatal(err)
	}

	log.Println(msg.Timestamp, " Process ", p.id, " on state ", p.getState(), " received a ", msg.getType(), " from ", msg.Id, " address: ", msg.Address, "with timestamp ", msg.Timestamp)

	p.updateTimestamp(msg.Timestamp)

	if msg.TypeMessage == REPLY || msg.TypeMessage == PERMISSION {
		p.incrementReply(msg)
	} else {
		if p.state == HELD || (p.state == WANTED && less(p, msg)) {
			log.Println(p.timestamp, " Process ", p.id, " on state ", p.getState(), " enqueued because ", p.requestTimestamp, " is less than ", msg.RequestTimestamp)
			p.enqueueMessage(msg)
			log.Println(p.timestamp, " Process ", p.id, " size counter: ", p.s.size(), " queue size: ", p.q.size())
		} else {
			p.sendMessage(REPLY, msg.Address)
		}
	}

}

func (p *process) enterOnCriticalSection() {
	//described on book
	p.changeState(WANTED)
	p.updateRequestTimestamp()
	p.doMulticast(REQUEST)
	p.waitAllProcessesReplies()
	p.changeState(HELD)

	//really enter in critical region
	p.getRandomString()

}

func (p *process) leaveCriticalSection() {
	//described on book
	p.changeState(RELEASED)
	p.replyAllEnqueuedRequests()
	p.clearReplyCounter()
}

func (p *process) runProcess() {
	log.Println(p.timestamp, " Process ", p.id, " is running ", p.address)

	for {
		log.Println(p.timestamp, " Process ", p.id, " IS TRYING TO ENTER IN CRITICAL REGION")

		p.enterOnCriticalSection()
		log.Println(p.timestamp, " Process ", p.id, " ENTERED CRITICAL REGION")

		p.leaveCriticalSection()
		log.Println(p.timestamp, " Process ", p.id, " LEFT CRITICAL REGION")
	}
}
