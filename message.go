package main

import (
    "encoding/gob"
    "net"
)

type message struct {
	RequestTimestamp int
	Timestamp int
	TypeMessage int
	Id int
	Address string
}

func (m * message) getType() string {
	s := ""
	if m.TypeMessage == REPLY {
		s = "REPLY"
	} else if m.TypeMessage == REQUEST {
		s = "REQUEST"
	} else {
		s = "PERMISSION"
	}
    return s
}

func (m * message) receiveAndDecodeMessage(connection net.Conn) error {
	dec := gob.NewDecoder(connection)
    err := dec.Decode(m)
    return err
}

func (m * message) encodeAndSendMessage(connection net.Conn) error {
	enc := gob.NewEncoder(connection)
    err := enc.Encode(m)
    return err
}
