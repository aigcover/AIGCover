package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
)

func main() {
	CheckConfigs()

	files, err := ioutil.ReadDir(SEEDPATH)
	if err != nil {
		log.Fatal(err)
	}

	fuzzer := Fuzzer{id: RANDID}

	len := len(files)
	for {
		file := files[rand.Intn(len)]
		for i := 0; i < 100; i++ {
			fmt.Println("Fuzzing:", SEEDPATH+file.Name(), i)
			fuzzer.Fuzz(SEEDPATH + file.Name())
		}
	}
}
