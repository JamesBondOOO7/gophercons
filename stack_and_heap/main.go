package main

// https://medium.com/eureka-engineering/understanding-allocations-in-go-stack-heap-memory-9a2631b5035d

type BigStruct struct {
	A, B, C int
	D, E, F string
	G, H, I bool
}

//go:noinline
func CreateCopy() BigStruct {
	return BigStruct{
		A: 123, B: 456, C: 789,
		D: "ABC", E: "DEF", F: "HIJ",
		G: true, H: true, I: true,
	}
}

//go:noinline
func CreatePointer() *BigStruct {
	return &BigStruct{
		A: 123, B: 456, C: 789,
		D: "ABC", E: "DEF", F: "HIJ",
		G: true, H: true, I: true,
	}
}
