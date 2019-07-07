package main

import (
	"sync"
)

type replySet struct {
	mutex sync.Mutex
	buffer map[int]bool
}

func (s * replySet ) newReplySet() {	
	s.buffer = make(map[int]bool)
}	

func (s * replySet) insertReply(msg message) {
	s.mutex.Lock()
	s.buffer[msg.Id] = true
	s.mutex.Unlock()
}

func (s * replySet) deleteReply(msg message) {
	s.mutex.Lock()
	delete(s.buffer, msg.Id)
	s.mutex.Unlock()
}

func (s * replySet) size() int {
	return len(s.buffer)
}

func (s * replySet) clear() {
	s.mutex.Lock()
	s.buffer = make(map[int]bool)
	s.mutex.Unlock()
}
