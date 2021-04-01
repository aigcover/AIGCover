package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {

	output, err :=
		exec.Command(AIGTOAIGPATH, "original.aag", "original.aig").CombinedOutput()
	if err != nil || string(output) != "" {
		log.Fatal(err, string(output))
		os.Exit(1)
	}

	pdr := PDRABCcheck("original.aig")
	bmc := BMCABCcheck("original.aig")

	println(pdr)
	println(bmc)

	switch {
	case pdr == CRASH:
		os.Exit(1)
	case bmc == CRASH:
		os.Exit(0)
	default:
		os.Exit(1)
	}
	os.Exit(1)
}
