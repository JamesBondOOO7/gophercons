package main

import "testing"

const creations = 20_000_000

func TestCopyIt(t *testing.T) {
	for i := 0; i < creations; i++ {
		_ = CreateCopy()
	}
}

func TestPointerIt(t *testing.T) {
	for i := 0; i < creations; i++ {
		_ = CreatePointer()
	}
}
