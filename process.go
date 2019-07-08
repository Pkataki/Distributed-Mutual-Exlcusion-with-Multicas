package main

import (
	"encoding/gob"
    "net"
    "log"
    "sync"
    "strconv"
    "time"
    //"bufio"
)

type process struct {
	mutex sync.Mutex
	address string
	id int
	state int
	timestamp int
	requestTimestamp int
	processesAddresses [] string
	s replyCounter
	q messageQueue
	channels [] chan message
	receivedAllReplies chan bool
	channelIndex map[string]int

}

func (p * process ) numberOfProcesses() int {
	return len(p.processesAddresses)
}

func (p * process ) addReply() {
	p.s.addReply()
	if p.s.size() == p.numberOfProcesses() - 1 {
		go func() {
			p.receivedAllReplies<-true
		}()		
	}
}

func (p * process ) getState() string {
	state := ""
	if p.state == WANTED {
		state = "WANTED"
	}

	if p.state == RELEASED {
		state = "RELEASED"
	}

	if p.state == HELD {
		state = "HELD"
	}
	return state
}

func (p * process ) clearReplySet() {
	p.s.clear()
}

func (p * process ) enqueueMessage(msg message) {
	p.q.push(msg)
}


func (p * process ) getMessageQueueTopRequest() string {
	msg := p.q.pop()
	return msg.Address
}

func (p * process ) isMessageQueueEmpty() bool {
	return p.q.empty() == true
}

func (p * process ) updateTimestamp(timestamp int) {
	p.mutex.Lock()
	p.timestamp = max(p.timestamp, timestamp) + 1
	p.mutex.Unlock()
}

func (p * process ) changeState(state int) {
	p.mutex.Lock()
	p.state = state
	p.mutex.Unlock()
	p.updateTimestamp(p.timestamp)
}

func (p * process ) updateRequestTimestamp() {
	p.requestTimestamp = p.timestamp
}

func (p * process ) getIndexFromAddress(address string) int {
	return p.channelIndex[address]
}

func (p * process) startProcess(address string, id int)  {
	
	p.s.newReplyCounter()
	p.q.newMessageQueue()
	p.channelIndex = make(map[string]int)
	p.receivedAllReplies = make(chan bool)
	p.address = address
	p.id = id
	
	p.changeState(RELEASED)
	if err := p.startListenPort(); err != nil {
		log.Fatal("Error on startListenPort")
	}

	log.Println("Process: ", p.id, " startListenPort: OK")

	if err := p.getOtherProcessesAddresses(); err != nil {
		log.Fatal("Error on getOtherProcessesAddresses")
	}


	log.Println("Process: ", p.id, " getOtherProcessesAddresses: OK")

	if err := p.openAllProcessesTCPConnections(); err != nil {
		log.Fatal("Error on openAllProcessesTCPConnections")
	}

	p.sendPermissionToAllProcesses()

	log.Println("Process: ", p.id, " is ready")
}

func (p * process) sendPermissionToAllProcesses( ) { 
	p.doMulticast(PERMISSION)
	p.waitAllProcessesReplies()
	p.clearReplySet()
}

func (p * process) openTCPConnection(address string) error {
	
	// connection, err := net.Dial("tcp", address)
	// if err != nil {
	// 	return err
	// }

	go func() {
		var msg message

		//defer connection.Close()
		for {
			msg = <-p.channels[p.getIndexFromAddress(address)]
			connection, err := net.Dial("tcp", address)
			if err != nil {
			//	return err
			}
			
			if err := msg.encodeAndSendMessage(connection); err != nil {
				log.Println(err)
			}
			connection.Close()
		}
	}()
	return nil
}

func (p * process) openAllProcessesTCPConnections() error {
	p.channels = make([]chan message, p.numberOfProcesses())
	for i := range p.processesAddresses {
    	p.channelIndex[p.processesAddresses[i]] = i
		if p.processesAddresses[i] != p.address { 

    		p.channels[i] = make(chan message)
    		
    		if err := p.openTCPConnection(p.processesAddresses[i]); err != nil {
    			return err
    		}
    	}
    }
    return nil

}

func (p * process) getOtherProcessesAddresses() error {

	conn, err := net.Dial("tcp", REGISTER_SERVER_ADDRESS)
	if err != nil {
		return err
	}
	defer conn.Close()

	enc := gob.NewEncoder(conn)
    err = enc.Encode(p.address)

    if err != nil {
    	return err
    }

	dec := gob.NewDecoder(conn)
    err = dec.Decode(&p.processesAddresses)

    if err != nil {
    	return err
    }


    if err != nil {
   		return err
    }

    return nil
}

func (p * process) sendMessage(typeMessage int, address string) {
	msg := message{
			Timestamp: p.timestamp,
			RequestTimestamp: p.requestTimestamp,
			TypeMessage: typeMessage,
			Address: p.address,
			Id: p.id,
		}

	go func() {
		//log.Printf("%d Channel of %s on Process %d on state %s sending message to %s",p.timestamp,address,p.id,p.getState(),address)
		p.channels[p.getIndexFromAddress(address)] <- msg
	}()
}

func (p * process) handleRequest(connection net.Conn) {
	var msg message

	if err := msg.receiveAndDecodeMessage(connection); err != nil {
		log.Println("error on receiveAndDecodeMessage")
		log.Fatal(err)
	}

	log.Println(msg.Timestamp," Process: ", p.id, " on state ", p.getState(), " received a ", msg.getType() ," from ", msg.Id, " address: ", msg.Address)
	p.updateTimestamp(msg.Timestamp)
    if msg.TypeMessage == REPLY || msg.TypeMessage == PERMISSION {
    	p.addReply()
    } else {
    	if p.state == HELD || (p.state == WANTED && less(p,msg)) {
    		log.Println(p.timestamp, " In Process ", p.id , " on state ", p.getState(), " enqueued because ",p.requestTimestamp, " is less than ", msg.RequestTimestamp)
    		p.enqueueMessage(msg)
    		log.Println(p.timestamp, "In Process ", p.id , " size set: " , p.s.size(), " queue size: ", p.q.size())
    	} else {
    		p.sendMessage(REPLY, msg.Address)
    	}
    }
}

func (p * process) startListenPort() error {
	listener, err := net.Listen("tcp", p.address)
	if err != nil {
		return err
	}
	go func() {
		defer listener.Close()
		for {
			conn, err := listener.Accept()
			//log.Println(p.timestamp, " Process ", p.id, " received message")
			if err != nil {
				log.Println(err)
			}
			go p.handleRequest(conn)	
		}
	}()
	return nil
}

func (p * process) doMulticast(typeMessage int) {
	log.Println(p.timestamp, " Process ", p.id, " in doMulticast is on state ", p.getState())
	for i := range p.processesAddresses {
		if p.processesAddresses[i] != p.address {
			p.sendMessage(typeMessage, p.processesAddresses[i])
			time.Sleep(200 * time.Millisecond)
		}
	}
}



func (p *process) getRandomString() {
	
	conn, err := net.Dial("tcp", CRITICAL_REGION_SERVER_ADDRESS)

	if err != nil {
		log.Println("Dial on ", p.address)
		log.Fatal(err)
	}

	defer conn.Close()

	msg := p.address + "|" + strconv.Itoa(p.timestamp)

	enc := gob.NewEncoder(conn)
	err = enc.Encode(msg)

	dec := gob.NewDecoder(conn)
    err = dec.Decode(&msg)

    log.Println(p.timestamp," PROCESS ", p.id, " received ", msg)

	if err != nil {
		log.Fatal("ERROR")
	}
}

func (p * process) waitAllProcessesReplies() {
	<-p.receivedAllReplies
}

func (p * process) enterOnCriticalSection() {
	p.changeState(WANTED)
	p.doMulticast(REQUEST)
	p.updateRequestTimestamp()
	p.waitAllProcessesReplies()
	p.changeState(HELD)

	//really enter in critical region 
	p.getRandomString()
	time.Sleep(1 * time.Second)
}

func (p * process) replyAllEnqueuedRequests() {
	for p.isMessageQueueEmpty() == false {
		address := p.getMessageQueueTopRequest()
		p.sendMessage(REPLY, address)
	}
} 

func (p * process) leaveCriticalSection() {
	p.changeState(RELEASED)
	p.replyAllEnqueuedRequests()
	p.clearReplySet()
}

func (p * process) runProcess() {
	log.Println(p.timestamp," Process ", p.id, " is running ", p.address)
	go func() {
		for {
			log.Println(p.timestamp, " Process: ", p.id, " IS TRYING TO ENTER IN CRITICAL REGION");
			p.enterOnCriticalSection();
			if p.state != HELD {
				log.Fatal(p.timestamp," Process: ", p.id, " is not on HELD state ")
			}
			log.Println(p.timestamp, " Process: ", p.id, " ENTERED CRITICAL REGION");
			time.Sleep(1 * time.Second)
			p.leaveCriticalSection()
			log.Println(p.timestamp, " Process: ", p.id, " LEFT CRITICAL REGION");
		}
	}()

	waiter := make(chan bool)
	<-waiter
}

