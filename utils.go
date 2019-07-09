package main

//servers addresses
const (
	REGISTER_SERVER_ADDRESS        = "localhost:8080"
	CRITICAL_REGION_SERVER_ADDRESS = "localhost:8081"
)

// types of message
const (
	REPLY = iota
	REQUEST
	PERMISSION
)

// states of process
const (
	WANTED = iota
	RELEASED
	HELD
)

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func less(p *process, msg message) bool {
	return checkProcessIsMinorThanMessage(p.requestTimestamp, p.id, msg.RequestTimestamp, msg.Id)

}

func checkProcessIsMinorThanMessage(timestamp_process, id_process, timestamp_message, id_message int) bool {
	if timestamp_process < timestamp_message {
		return true
	} else if (timestamp_process == timestamp_message) && (id_process < id_message) {
		return true
	} else {
		return false
	}
}
