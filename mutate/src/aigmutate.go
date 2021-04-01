package main

import (
	"os/exec"
	"log"
	"os"
	"math/rand"
	"time"
	"fmt"
	"path/filepath"
)

func main()  {
	rand.Seed(time.Now().UnixNano())	
	fileIn := os.Args[1]
	fileOut := os.Args[2]
	dir , err := os.Executable()
	if err != nil {
		fmt.Println(err)
	}
	exPath := filepath.Dir(dir)	
	aigmutate(exPath,fileIn,fileOut)
}

func aigmutate(exPath, fileIn, fileOut string) {
	t := fmt.Sprintf("%d_%v",os.Getpid(),time.Now().UnixNano())
	fileInAAG := exPath+"/Temp/tempin"+string(t)+".aag"
	fileOutAAG := exPath+"/Temp/tempout"+string(t)+".aag"
	aigtoaig := exPath+"/aiger-1.9.9/aigtoaig"
	output, err :=
		exec.Command(aigtoaig, fileIn, fileInAAG).CombinedOutput()
	if err != nil || string(output) != "" {
		log.Fatal(err, string(output))
	}
	network := Network{}
	network.Read(fileInAAG)
	network.mutation()
	network.Dumpto(fileOutAAG)

	output2, err2 :=
		exec.Command(aigtoaig, fileOutAAG, fileOut).CombinedOutput()
	if err2 != nil || string(output2) != "" {
		log.Fatal(err2, string(output2))
	}
	err3 := os.Remove(fileInAAG) 
	if err3 != nil {
		log.Fatal(err3,"Delete temp file fail")
	} 
	err3 = os.Remove(fileOutAAG)
	if err3 != nil {
		log.Fatal(err3,"Delete temp file fail")
	}
}

func (N *Network) mutation() {
	const (
		outputs    = iota
		latches    = iota
		properties = iota
		invariants = iota
		gates      = iota
	)

	for count := 0; N.nIndex>0&&count < rand.Intn(N.nIndex)%100; count++ {
		newPin := pin{}
		randomid := rand.Intn(N.nIndex)
		if randomid < 2 {
			randomid = 2
		}
		newPin.Read(randomid)
		rInt := rand.Intn(5)
		switch rInt {
		case outputs:
			if len := len(N.outputs); len != 0 {
				N.outputs[rand.Intn(len)] = newPin
			}
		case latches:
			if len := len(N.latches); len != 0 {
				switch rand.Intn(2) {
				case 0:
					N.latches[rand.Intn(len)].next = newPin
				case 1:
					N.latches[rand.Intn(len)].init.Flip()
				}
			}

		case properties:
			if len := len(N.properties); len != 0 {
				N.properties[rand.Intn(len)] = newPin
			}
		case invariants:
			if len := len(N.invariants); len != 0 {
				N.invariants[rand.Intn(len)] = newPin
			}
		case gates:
			if len := len(N.gates); len != 0 {
				index := rand.Intn(len)
				if newPin.id != N.gates[index].output.id {
					if !N.IsCyclic(N.gates[index].output, newPin) {
						N.gates[index].first = newPin
						switch rand.Intn(2) {
						case 0:
							N.gates[index].first = newPin
						case 1:
							N.gates[index].second = newPin
						}
					} else {
						//println("Is cyclic definition")
					}

				}
			}
		}
	}
}