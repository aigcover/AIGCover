package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {

	if len(os.Args) != 3 {
		fmt.Println("Usage: ./reducer [interestingness_test] [file_to_reduce]")
		fmt.Println("[file_to_reduce]:	The aag file that will being reduced")
		fmt.Println(`[interestingness_test]: The interestingness test is an
			executable program (usually a shell script) that returns 0 when a 
			partially reduced file is interesting (a candidate for further   
			reduction) and returns non-zero when a partially reduced file is not
			interesting (not a candidate for further reduction -- all
			uninteresting files are discarded).`)
		return
	}

	interestingnessTest := os.Args[1]
	file := os.Args[2]
	if !strings.Contains(interestingnessTest, "\\") {
		interestingnessTest = "./" + interestingnessTest
	}

	fmt.Println("Start Dry run...")
	output, err :=
		exec.Command(interestingnessTest).CombinedOutput()
	if err != nil {
		fmt.Println("Dry run failed!")
		log.Fatal(err, string(output))
	}
	fmt.Println("Dry run success, start reducing...")

	network := Network{}
	network.Read(file)
	network.reduce(interestingnessTest, file)
}

func (N *Network) reduce(interestingnessTest string, file string) {

	reduced := Network{}
	reduced.Read(file)

	fileStat, err := os.Stat(file)
	if err != nil {
		log.Fatal(err)
	}
	beforeSize := fileStat.Size()
	smaller := true

	for smaller {

		smaller = false
		for _, i := range N.inputs {
			succ := N.tryRemovePin(i, interestingnessTest, file)
			if succ {
				reduced = Network{}
				reduced.Read(file)
				smaller = true
			} else {
				reduced.Dumpto(file)
				N = &Network{}
				N.Read(file)
			}
		}

		for _, l := range N.latches {
			succ := N.tryRemovePin(l.output, interestingnessTest, file)
			if succ {
				reduced = Network{}
				reduced.Read(file)
				smaller = true
			} else {
				reduced.Dumpto(file)
				N = &Network{}
				N.Read(file)
			}

		}

		for _, o := range N.outputs {
			succ := N.tryRemovePin(o, interestingnessTest, file)
			if succ {
				reduced = Network{}
				reduced.Read(file)
				smaller = true
			} else {
				reduced.Dumpto(file)
				N = &Network{}
				N.Read(file)
			}
		}

		for _, p := range N.properties {
			succ := N.tryRemovePin(p, interestingnessTest, file)
			if succ {
				reduced = Network{}
				reduced.Read(file)
				smaller = true
			} else {
				reduced.Dumpto(file)
				N = &Network{}
				N.Read(file)
			}
		}

		for _, in := range N.invariants {
			succ := N.tryRemovePin(in, interestingnessTest, file)
			if succ {
				reduced = Network{}
				reduced.Read(file)
				smaller = true
			} else {
				reduced.Dumpto(file)
				N = &Network{}
				N.Read(file)
			}
		}

		for _, g := range N.gates {
			succ := N.tryRemovePin(g.output, interestingnessTest, file)
			if succ {
				reduced = Network{}
				reduced.Read(file)
				smaller = true
			} else {
				reduced.Dumpto(file)
				N = &Network{}
				N.Read(file)
			}
		}

		fileStat, err = os.Stat(file)
		if err != nil {
			log.Fatal(err)
		}

		reduced.Dumpto(file)
		N = &Network{}
		N.Read(file)

		afterSize := fileStat.Size()

		ratio := fmt.Sprintf("%.2f", (1-float64(afterSize)/float64(beforeSize))*100)

		fmt.Println("("+ratio, "%,", afterSize, "bytes)")
	}
}

func (N *Network) tryRemovePin(P pin, interestingnessTest string, file string) bool {

	N.removePin(P)

	N.Dumpto(file)

	output, err :=
		exec.Command(interestingnessTest).CombinedOutput()
	if err != nil {
		//log.Fatal(err, string(output))
		fmt.Println("unsucc")
		return false
	}

	fmt.Println(string(output))
	return true
}

func (N *Network) removePin(P pin) {

	tmp := N.inputs[:0]
	var pins []pin
	for _, i := range N.inputs {
		if !P.SameIDAs(i) {
			tmp = append(tmp, i)
			pins = append(pins, i)
		} else {
			N.nInputs--
		}
	}
	N.inputs = tmp
	for _, p := range pins {
		N.removePin(p)
	}

	tmpl := N.latches[:0]
	pins = pins[:0]
	for _, l := range N.latches {
		if !(P.SameIDAs(l.output)) {
			tmpl = append(tmpl, l)
		}
		if (P.SameIDAs(l.next) || P.SameIDAs(l.init)) && !P.SameIDAs(l.output) {
			pins = append(pins, l.output)
		}
	}
	N.nLatches = N.nLatches - (len(N.latches) - len(tmpl))
	N.latches = tmpl
	for _, p := range pins {
		N.removePin(p)
	}

	tmp = N.outputs[:0]
	for _, o := range N.outputs {
		if !P.SameIDAs(o) {
			tmp = append(tmp, o)
		} else {
			N.nOutputs--
		}
	}
	N.outputs = tmp

	tmp = N.properties[:0]
	for _, p := range N.properties {
		if !P.SameIDAs(p) {
			tmp = append(tmp, p)
		} else {
			N.nProperties--
		}
	}
	N.properties = tmp

	tmp = N.invariants[:0]
	for _, in := range N.invariants {
		if !P.SameIDAs(in) {
			tmp = append(tmp, in)
		} else {
			N.nInvariant--
		}
	}
	N.invariants = tmp

	tmpg := N.gates[:0]
	pins = pins[:0]
	for _, g := range N.gates {
		if !P.SameIDAs(g.output) {
			tmpg = append(tmpg, g)
		}
		if (P.SameIDAs(g.first) || P.SameIDAs(g.second)) && !P.SameIDAs(g.output) {
			pins = append(pins, g.output)
		}
	}
	N.nGates = N.nGates - (len(N.gates) - len(tmpg))
	N.gates = tmpg
	for _, p := range pins {
		N.removePin(p)
	}
}
