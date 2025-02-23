package main

import (
	"fmt"
	"time"
)

// https://medium.com/globant/pub-sub-in-golang-an-introduction-8be4c65eafd4

func main() {
	// channel to publish messages to
	msgChannel := make(chan string)

	// goroutine to publish messages
	go publishingMessage("Hello from Globant", msgChannel)

	// goroutine to receive messages
	go receivingMessage(msgChannel)

	time.Sleep(1 * time.Second)
	fmt.Println("main goroutine exit...")
}

// function to publish messages to the channel
func publishingMessage(message string, msgChannel chan string) {
	msgChannel <- message
}

// function to receive messages from the channel
// As this is an open-unbuffered channel, this function will always block indefinitely until the channel is closed
func receivingMessage(msgChannel chan string) {
	for {
		fmt.Println("I am blocked")
		msg := <-msgChannel
		fmt.Println("Received message:", msg)
	}
	fmt.Println("I am unblocked")
}
