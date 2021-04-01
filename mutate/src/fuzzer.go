package main

import (
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"strconv"
	"os"
)

/*Fuzzer instance*/
type Fuzzer struct {
	count int
	id    string
}

/*Fuzz the file*/
func (fuzzer *Fuzzer) Fuzz(file string) {

	network := Network{}

	network.Read(file)
	network.mutation()
	network.Dumpto(AAGPATH)

	output, err :=
		exec.Command(AIGTOAIGPATH, AAGPATH, AIGPATH).CombinedOutput()
	if err != nil || string(output) != "" {
		log.Fatal(err, string(output))
	}

	//pdr := PDRABCcheck(AIGPATH)
	pdr,pdrOut := UseABCcheck(AIGPATH, "pdr")
	//inter := INTABCcheck(AIGPATH)
	inter,intOut := UseABCcheck(AIGPATH, "int")
	//bmc := BMCABCcheck(AIGPATH)
	bmc,bmcOut := UseABCcheck(AIGPATH, "bmc2")
	tempor,temporOut := UseABCcheck(AIGPATH, "tempor")
	dprove,dproveOut := UseABCcheck(AIGPATH, "dprove")
	//IC3 := UseIC3check(AIGPATH)
	IC3,IC3Out := UseIC3checkWithTimeout(AIGPATH)

	var checkList = [...]int{pdr, inter, bmc, tempor, dprove,IC3}
	var checkOut = [...]string{pdrOut,intOut,bmcOut,temporOut,dproveOut,IC3Out}
	var checkListName = [...]string{"pdr", "inter", "bmc", "tempor", "dprove","IC3"}
	//var checkList = [...]int{pdr, inter, bmc, tempor, dprove}
	//var checkListName = [...]string{"pdr", "inter", "bmc", "tempor", "dprove"}
	var bugType = [...]string{"sat", "unsat", "timeout", "crash", "ignore"}

	/*
		var compare=-1
		for index, check_info := range checkList{
			if(check_info != TIMEOUT && check_info != IGNORE){
				compare=check_info
				break
			}

		}
	*/
	for i, checkInfo := range checkList {
		print(checkListName[i], ":", bugType[checkInfo], "   ")
		if checkInfo == CRASH {
			filename := fmt.Sprintf("%s_%s_%s_%s", fuzzer.id, strconv.Itoa(fuzzer.count), "crash", checkListName[i])
			saveas(network, filename,checkOut[i])
			fuzzer.count++
			println(filename)
		}
		for j := range checkList {
			if notequal(checkList[i], checkList[j]) {
				filename := fmt.Sprintf("%s_%s_%s", fuzzer.id, strconv.Itoa(fuzzer.count), "incorrect")
				saveas(network, filename,checkOut[i])
				fuzzer.count++
				println(filename)
				return
			}
		}
	}
	println(" ")

	/*
		for i:=0;i<len(checkList)-1;i++{
			if(check_info == CRASH){
				filename := fmt.Sprintf("%s_%s_%s_%s", fuzzer.id, strconv.Itoa(fuzzer.count), "crash", "pdr")
				saveas(network, filename)
				fuzzer.count++
				println(filename)
			}
			if(notequal(checkList[i],checkList[j])){
				filename := fmt.Sprintf("%s_%s_%s", fuzzer.id, strconv.Itoa(fuzzer.count), "incorrect")
				saveas(network, filename)
				fuzzer.count++
				println(filename)
				return
			}
		}*/

	// switch {
	// case pdr == CRASH:
	// 	filename := fmt.Sprintf("%s_%s_%s_%s", fuzzer.id, strconv.Itoa(fuzzer.count), "crash", "pdr")
	// 	saveas(network, filename)
	// 	fuzzer.count++
	// 	println(filename)
	// 	fallthrough
	// case inter == CRASH:
	// 	filename := fmt.Sprintf("%s_%s_%s_%s", fuzzer.id, strconv.Itoa(fuzzer.count), "crash", "inter")
	// 	saveas(network, filename)
	// 	fuzzer.count++
	// 	println(filename)
	// 	fallthrough
	// case bmc == CRASH:
	// 	filename := fmt.Sprintf("%s_%s_%s_%s", fuzzer.id, strconv.Itoa(fuzzer.count), "crash", "bmc")
	// 	saveas(network, filename)
	// 	fuzzer.count++
	// 	println(filename)
	// case notequal(pdr, inter):
	// 	filename := fmt.Sprintf("%s_%s_%s", fuzzer.id, strconv.Itoa(fuzzer.count), "incorrect")
	// 	saveas(network, filename)
	// 	fuzzer.count++
	// 	println(filename)
	// case notequal(pdr, bmc):
	// 	filename := fmt.Sprintf("%s_%s_%s", fuzzer.id, strconv.Itoa(fuzzer.count), "incorrect")
	// 	saveas(network, filename)
	// 	fuzzer.count++
	// 	println(filename)
	// case notequal(bmc, inter):
	// 	filename := fmt.Sprintf("%s_%s_%s", fuzzer.id, strconv.Itoa(fuzzer.count), "incorrect")
	// 	saveas(network, filename)
	// 	fuzzer.count++
	// 	println(filename)
	// default:
	// }

}

func notequal(first int, second int) bool {
	if first != TIMEOUT && second != TIMEOUT &&
		first != IGNORE && second != IGNORE &&
		first != CRASH && second != CRASH &&
		first != second {
		return true
	}
	return false
}
/*
func saveas(network Network, filename string) {
	aagpath := BUGPATH + "/" + filename + ".aag"
	aigpath := BUGPATH + "/" + filename + ".aig"
	network.Dumpto(aagpath)
	output, err := exec.Command(AIGTOAIGPATH, aagpath, aigpath).CombinedOutput()
	if err != nil || string(output) != "" {
		log.Fatal(err, string(output))
	}
}
*/
func saveas(network Network, filename string,out string) {
	aagpath := BUGPATH + "/" + filename + ".aag"
	aigpath := BUGPATH + "/" + filename + ".aig"
	outpath := BUGPATH + "/" + filename + ".txt"
	network.Dumpto(aagpath)
	output, err := exec.Command(AIGTOAIGPATH, aagpath, aigpath).CombinedOutput()
	if err != nil || string(output) != "" {
		log.Fatal(err, string(output))
	}
	f , err1 := os.Create(outpath)
	if err1 != nil {
		log.Fatal(err1)
	}
	_, err2 := f.WriteString(out)
	if err2 != nil {
		log.Fatal(err2)
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


