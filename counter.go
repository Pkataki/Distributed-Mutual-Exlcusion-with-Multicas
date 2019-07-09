package main

import (
	"sync"
)

type replyCounter struct {
	mutex  sync.Mutex
	buffer map[int]bool
}

func NewReplyCounter() replyCounter {
	return replyCounter{buffer: make(map[int]bool)}
}

func (s *replyCounter) newReplyCounter() {
	s.buffer = make(map[int]bool)
}

func (s *replyCounter) addReply(msg message) {
	s.mutex.Lock()
	s.buffer[msg.Id] = true
	s.mutex.Unlock()
}

func (s *replyCounter) size() int {
	return len(s.buffer)
}

func (s *replyCounter) clear() {
	s.mutex.Lock()
	s.buffer = make(map[int]bool)
	s.mutex.Unlock()
}
