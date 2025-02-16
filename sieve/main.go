package main

import "fmt"

// generate : send the sequence 2, 3, 4 ... to channel ch
func generate(ch chan<- int) {
	for i := 2; ; i++ {
		ch <- i
	}
}

// filter : copy the values from channel `in` to channel `out`
// removing those divisible by `prime`
func filter(in <-chan int, out chan<- int, prime int) {
	for {
		// receive the value `i` from channel `in`
		i := <-in
		if i%prime != 0 {
			out <- i // send `i` to channel out
		}
	}
}

func sieve() {
	// ch <- (2, 3, 4, 5, 6, 7, 8, 9, 10, 11, ...)
	// iteration 1 : prime = 2
	//    ch -- supplying numbers to filter() -- throwing numbers to -- ch1
	//    ch1 <- (3, 5, 7, 9, 11, ...)
	//    ch now references to ch1 but the above flow is ongoing continuously
	// iteration 2 : prime = 3
	//    ch -- supplying numbers to filter() -- throwing numbers to -- ch1
	//    ch1 <- (3, 5, 7, 11, ...)
	//    ch now references to ch1 but the above flow is ongoing continuously
	// ...
	ch := make(chan int)
	go generate(ch)
	for {
		prime := <-ch
		fmt.Println(prime)
		ch1 := make(chan int)
		go filter(ch, ch1, prime)
		ch = ch1
	}
}

func main() {
	sieve()
}
