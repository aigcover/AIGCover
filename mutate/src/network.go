package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func test() {
	fmt.Printf("File contents")
}

/*Network for modeling the model*/
type Network struct {
	nIndex      int /* maximum variable index                 */
	nInputs     int /* number of inputs                       */
	nLatches    int /* number of latches                      */
	nOutputs    int /* number of outputs                      */
	nGates      int /* number of AND gates                    */
	nProperties int /* number of "bad state" properties       */
	nInvariant  int /* number of AND invariant constraints    */
	//nJustice    int;    /* number of AND justice properties       */
	//nFairness   int;    /* number of AND fairness constraints     */
	inputs     []pin   /* inputs of the network                */
	outputs    []pin   /* outputs of the network               */
	latches    []latch /* latches of the network               */
	properties []pin   /* properties of the network            */
	invariants []pin   /* invariants of the network            */
	gates      []gate  /* gates of the network                 */
}

type latch struct {
	output pin
	next   pin
	init   pin
}

type gate struct {
	output pin
	first  pin
	second pin
}

type pin struct {
	id     int
	status bool
}

func (P *pin) Read(val int) pin {
	if val%2 == 1 {
		P.id = val - 1
		P.status = false
	} else {
		P.id = val
		P.status = true
	}
	return *P
}

func (P *pin) Flip() {
	P.status = !P.status
}

func (P pin) Eval() (val int) {
	if P.status == false {
		val = P.id + 1
	} else {
		val = P.id
	}
	return
}

func (P pin) SameIDAs(P2 pin) bool {
	return P.id == P2.id
}

func (N *Network) Read(file string) {

	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(content), "\n")

	title := strings.Split(lines[0], " ")
	if title[0] != "aag" || len(title) < 6 || len(title) > 8 {
		log.Fatal(file, " is not a valid aag file.")
	}

	N.nIndex, err = strconv.Atoi(title[1])
	N.nInputs, err = strconv.Atoi(title[2])
	N.nLatches, err = strconv.Atoi(title[3])
	N.nOutputs, err = strconv.Atoi(title[4])
	N.nGates, err = strconv.Atoi(title[5])
	if len(title) > 6 {
		N.nProperties, err = strconv.Atoi(title[6])
	}
	if len(title) > 7 {
		N.nInvariant, err = strconv.Atoi(title[7])
	}
	if err != nil {
		log.Fatal(err)
	}

	start := 1
	end := start + N.nInputs
	for _, line := range lines[start:end] {
		tokens := strings.Split(line, " ")
		if len(tokens) != 1 {
			log.Fatal("input", tokens, " Incorrect token numbers.")
		}
		input, err := strconv.Atoi(tokens[0])
		if err != nil {
			log.Fatal(err)
		}
		inputPin := pin{}
		N.inputs = append(N.inputs, inputPin.Read(input))
	}

	start = end
	end = end + N.nLatches
	for _, line := range lines[start:end] {
		tokens := strings.Split(line, " ")
		if len(tokens) > 3 || len(tokens) < 2 {
			log.Fatal("latch", tokens, " Incorrect token numbers.")
		}
		output, err := strconv.Atoi(tokens[0])
		next, err := strconv.Atoi(tokens[1])
		init := 0
		if len(tokens) == 3 {
			init, err = strconv.Atoi(tokens[2])
			if init != 1 {
				init = 0
			}
		}
		if err != nil {
			log.Fatal(err)
		}
		outputPin := pin{}
		nextPin := pin{}
		initPin := pin{}
		latch := latch{outputPin.Read(output),
			nextPin.Read(next),
			initPin.Read(init)}
		N.latches = append(N.latches, latch)
	}

	start = end
	end = end + N.nOutputs
	for _, line := range lines[start:end] {
		tokens := strings.Split(line, " ")
		if len(tokens) != 1 {
			log.Fatal("output", tokens, " Incorrect token numbers.")
		}
		output, err := strconv.Atoi(tokens[0])
		if err != nil {
			log.Fatal(err)
		}
		outputPin := pin{}
		N.outputs = append(N.outputs, outputPin.Read(output))
	}

	start = end
	end = end + N.nProperties
	for _, line := range lines[start:end] {
		tokens := strings.Split(line, " ")
		if len(tokens) != 1 {
			log.Fatal("property", tokens, " Incorrect token numbers.")
		}
		property, err := strconv.Atoi(tokens[0])
		if err != nil {
			log.Fatal(err)
		}
		propertyPin := pin{}
		N.properties = append(N.properties, propertyPin.Read(property))
	}

	start = end
	end = end + N.nInvariant
	for _, line := range lines[start:end] {
		tokens := strings.Split(line, " ")
		if len(tokens) != 1 {
			log.Fatal("invariant", tokens, " Incorrect token numbers.")
		}
		invariant, err := strconv.Atoi(tokens[0])
		if err != nil {
			log.Fatal(err)
		}
		invariantPin := pin{}
		N.invariants = append(N.invariants, invariantPin.Read(invariant))
	}

	start = end
	end = end + N.nGates
	for _, line := range lines[start:end] {
		tokens := strings.Split(line, " ")
		if len(tokens) != 3 {
			log.Fatal("gate", tokens, " Incorrect token numbers.")
		}
		output, err := strconv.Atoi(tokens[0])
		first, err := strconv.Atoi(tokens[1])
		second, err := strconv.Atoi(tokens[2])
		if err != nil {
			log.Fatal(err)
		}
		outputPin := pin{}
		firstPin := pin{}
		secondPin := pin{}
		gate := gate{outputPin.Read(output),
			firstPin.Read(first),
			secondPin.Read(second)}
		N.gates = append(N.gates, gate)
	}
}

/*Print the info of the network*/
func (N Network) Print() {
	fmt.Println("nIndex:      ", N.nIndex)
	fmt.Println("nInputs:     ", N.nInputs)
	fmt.Println("nLatches:    ", N.nLatches)
	fmt.Println("nOutputs:    ", N.nOutputs)
	fmt.Println("nGates:      ", N.nGates)
	fmt.Println("nProperties: ", N.nProperties)
	fmt.Println("nInvariant:  ", N.nInvariant)
	fmt.Println("inputs:      ", N.inputs)
	fmt.Println("outputs:     ", N.outputs)
	fmt.Println("latches:     ", N.latches)
	fmt.Println("properties:  ", N.properties)
	fmt.Println("invariants:  ", N.invariants)
	fmt.Println("gates:       ", N.gates)

}

/*Dumpto the network to the file*/
func (N Network) Dumpto(file string) {
	writer, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Println(err)
	}
	defer writer.Close()

	title := fmt.Sprintf("%s %d %d %d %d %d %d %d",
		"aag", N.nIndex, N.nInputs, N.nLatches, N.nOutputs,
		N.nGates, N.nProperties, N.nInvariant)
	fmt.Fprintln(writer, title)

	for i := 0; i < N.nInputs; i++ {
		line := fmt.Sprintf("%d", N.inputs[i].Eval())
		fmt.Fprintln(writer, line)
	}

	for i := 0; i < N.nLatches; i++ {
		line := fmt.Sprintf("%d %d %d", N.latches[i].output.Eval(), N.latches[i].next.Eval(), N.latches[i].init.Eval())
		fmt.Fprintln(writer, line)
	}

	for i := 0; i < N.nOutputs; i++ {
		line := fmt.Sprintf("%d", N.outputs[i].Eval())
		fmt.Fprintln(writer, line)
	}

	for i := 0; i < N.nProperties; i++ {
		line := fmt.Sprintf("%d", N.properties[i].Eval())
		fmt.Fprintln(writer, line)
	}

	for i := 0; i < N.nInvariant; i++ {
		line := fmt.Sprintf("%d", N.invariants[i].Eval())
		fmt.Fprintln(writer, line)
	}

	for i := 0; i < N.nGates; i++ {
		line := fmt.Sprintf("%d %d %d", N.gates[i].output.Eval(), N.gates[i].first.Eval(), N.gates[i].second.Eval())
		fmt.Fprintln(writer, line)
	}
}

/*IsCyclic : check if lhs dependency relation to rhs is cyclic */
func (N Network) IsCyclic(lhs pin, rhs pin) bool {
	var rhsstack []pin
	var checkedpinstack []pin
	rhsstack = append(rhsstack, rhs)
	for len(rhsstack) != 0 {
		checkedpin := rhsstack[len(rhsstack)-1]
		checkedpinstack = append(checkedpinstack, checkedpin)
		rhsstack = rhsstack[:len(rhsstack)-1]
		for i := 0; i < N.nGates; i++ {
			if N.gates[i].output.SameIDAs(checkedpin) {
				if N.gates[i].first.SameIDAs(lhs) || N.gates[i].second.SameIDAs(lhs) {
					return true
				}
				firstflag := true
				secondflag := true
				for j := 0; j < len(checkedpinstack); j++ {
					if N.gates[i].first.SameIDAs(checkedpinstack[j]) {
						firstflag = false
					}
					if N.gates[i].second.SameIDAs(checkedpinstack[j]) {
						secondflag = false
					}
				}
				if firstflag {
					rhsstack = append(rhsstack, N.gates[i].first)
				}
				if secondflag {
					rhsstack = append(rhsstack, N.gates[i].second)
				}
			}
		}
	}
	return false
}
