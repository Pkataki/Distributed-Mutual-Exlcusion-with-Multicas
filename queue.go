package main 

type messageQueue struct {
	buffer chan message
}

func (q * messageQueue ) newMessageQueue () {
	q.buffer = make(chan message)
}	

func (q * messageQueue) push(msg message) {
	go func() {
		q.buffer <- msg
	}()
}

func (q * messageQueue) pop() message {
	msg := <- q.buffer
	return msg
}

func (q * messageQueue) empty() bool {
	return len(q.buffer) == 0
}

func (q * messageQueue) size() int {
	return len(q.buffer)
}
