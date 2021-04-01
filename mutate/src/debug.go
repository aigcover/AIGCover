package main

import "os"

func main() {

	avy := AVYcheck(os.Args[1])
	if avy == SAT {
		os.Exit(0)
	}
}
