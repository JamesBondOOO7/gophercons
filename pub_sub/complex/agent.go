package main

import (
	"sync"
)

type Agent struct {
	// mu protects access to the subs and closed fields using a mutex, a synchronization mechanism that allows only one
	//goroutine to access these fields at a time
	mu sync.Mutex

	// subs field is a map that holds a list of channels for each topic, allowing subscribers to receive messages
	//published to that topic
	subs map[string][]chan string

	// quit field is a channel that is closed when the agent is closed, allowing goroutines that are blocked on the
	//channel to unblock and exit
	quit chan struct{}

	// closed field is a flag that indicates whether the agent has been closed.
	closed bool
}

func NewAgent() *Agent {
	return &Agent{
		subs: make(map[string][]chan string),
		quit: make(chan struct{}),
	}
}

func (a *Agent) Publish(topic string, msg string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.closed {
		return
	}

	for _, subscriberChannel := range a.subs[topic] {
		subscriberChannel <- msg
	}
}

func (a *Agent) Subscribe(topic string) <-chan string {
	a.mu.Lock()
	defer a.mu.Unlock()

	//// THINK: Who should add the topics? For now any subscriber can add any topic
	//_, ok := a.subs[topic]
	//if !ok {
	//	// THINK: An unbuffered channel will block indefinitely, why??
	//	ch := make(chan string, 1)
	//	ch <- "Topic doesn't exist"
	//	close(ch)
	//	return ch
	//}

	// THINK: can we keep this as a buffered channel?
	newSubscriberChannel := make(chan string)
	a.subs[topic] = append(a.subs[topic], newSubscriberChannel)
	return newSubscriberChannel
}

func (a *Agent) Close() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.closed {
		return
	}

	a.closed = true
	close(a.quit)

	// close all channels for subscribers
	for _, subscriberChannelForATopic := range a.subs {
		for _, subscriberChannel := range subscriberChannelForATopic {
			close(subscriberChannel)
		}
	}
}
