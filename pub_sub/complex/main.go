package main

// https://medium.com/globant/pub-sub-in-golang-an-introduction-8be4c65eafd4

import (
	"fmt"
)

func main() {
	// Create a new agent
	agent := NewAgent()

	// Subscribe to a topic
	sub := agent.Subscribe("foo")

	// Publish a message to the topic
	go agent.Publish("foo", "hello world")

	// Print the message
	fmt.Println(<-sub)

	// Close the agent
	agent.Close()
}
