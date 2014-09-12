package main

import (
	"fmt"
	"net/http"
)

// Empty struct to attach our endpoint methods to.
type MessengerService struct {
}

// Message just wraps a string in a struct.  We could add other fields here if we wanted.
type Message struct {
	Content string
}

// Echo repeats the message you send back to you.
func (ms *MessengerService) Echo(r *http.Request, m *Message, resp *Message) error {
	resp.Content = "Recieved message: " + m.Content
	return nil
}

func main() {
	// Create instance of messenger service.
	service := MessengerService{}

	// Create message to send.
	m1 := Message{Content: "Hello!"}

	// Message to receive
	var m2 Message

	service.Echo(nil, &m1, &m2)
	fmt.Println(m2)
}
