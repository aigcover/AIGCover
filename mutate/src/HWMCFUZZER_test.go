package main

import (
	"testing"
	"os"
	"log"
	"context"
	"fmt"
	"os/exec"
	"time"
)

func TestChecker(t *testing.T) {
	SATExample := "/mnt/d/softwareTesting/HWMCFuzz/benchmark/example/SAT.aig"
	UNSATExample := "/mnt/d/softwareTesting/HWMCFuzz/benchmark/example/UNSAT.aig"	
    pdr, _ := UseABCcheck(SATExample, "pdr")
	inter, _ := UseABCcheck(SATExample, "int")
	bmc, _ := UseABCcheck(SATExample, "bmc2")
	tempor, _ := UseABCcheck(SATExample, "tempor")
	dprove, _ := UseABCcheck(SATExample, "dprove")
	IC3, _:= UseIC3checkWithTimeout(SATExample)
	var checkList = [...]int{pdr, inter, bmc, tempor, dprove, IC3}
	for _, checkInfo := range checkList {
		if checkInfo != 0 && checkInfo != 4 && checkInfo != 2 {
			t.Errorf("Check was incorrect, got: %d, want: %d or %d or %d.", checkInfo, 0, 2, 4)
		 } 
	}
	pdr, _= UseABCcheck(UNSATExample, "pdr")
	inter, _ = UseABCcheck(UNSATExample, "int")
	bmc, _ = UseABCcheck(UNSATExample, "bmc2")
	tempor, _ = UseABCcheck(UNSATExample, "tempor")
	dprove, _ = UseABCcheck(UNSATExample, "dprove")
	IC3, _ = UseIC3checkWithTimeout(UNSATExample)
	checkList = [...]int{pdr, inter, bmc, tempor, dprove, IC3}
	for _, checkInfo := range checkList {
		if checkInfo != 1 && checkInfo != 4 && checkInfo != 2 {
			t.Errorf("Check was incorrect, got: %d, want: %d or %d or %d.", checkInfo, 1, 2, 4)
		 } 
	}
}

func TestIC3Crash(t *testing.T){
	IC3CRASHExample := "/mnt/d/softwareTesting/HWMCFuzz/benchmark/example/IC3CRASH.aig"
	IC3, _ := UseIC3checkWithTimeout(IC3CRASHExample)
	if(IC3==CRASH){
		println("IC3CRASH Happened")
	}
}

func TestComparer(t *testing.T)  {
	var checkList = [...]int{0, 1, 0, 1, 2, 4}
	var checkListName = [...]string{"pdr", "inter", "bmc", "tempor", "dprove","IC3"}
	var bugType = [...]string{"sat", "unsat", "timeout", "crash", "ignore"}

	for i, checkInfo := range checkList {
		print(checkListName[i], ":", bugType[checkInfo], "   ")
		if checkInfo == CRASH {
			if checkInfo != 3 {
				t.Errorf("Crash was incorrect, got: %d, want: %d.", checkInfo, 3)
			}
		}
		for j := range checkList {
			if notequal(checkList[i], checkList[j]) {
				if (checkList[i] == 0 || checkList[i] == 1)	&& (checkList[j] == 0 || checkList[j] == 1){
					if(checkList[i]==checkList[j]){
						t.Errorf("Compare err,want not equal but get %d and %d.",checkList[i],checkList[j])
					}
				}			
			}
		}
	}
}
func TestSave (t *testing.T)  {
	CRASHExample := "/mnt/d/softwareTesting/HWMCFuzz/benchmark/example/CRASH.aig"
	_, bmcout := UseABCcheck(CRASHExample, "bmc2")
	//pdr, pdrout := UseABCcheck(CRASHExample, "pdr")
	outpath := "/mnt/d/softwareTesting/HWMCFuzz/benchmark/example/CRASH.txt"
	println(bmcout)
	//println(pdrout)
	f , err1 := os.Create(outpath)
	if err1 != nil {
		log.Fatal(err1)
	}
	_, err2 := f.WriteString("Crash Log:\n"+bmcout)
	if err2 != nil {
		log.Fatal(err2)
	}
}

func TestOut(t *testing.T){
	file := "/mnt/d/softwareTesting/HWMCFuzz/benchmark/example/CRASH.aig"
	checker := "bmc3"
	command := fmt.Sprintf("read %s;logic;undc;strash;zero;fold;%s;", file, checker)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	output, err :=
		exec.CommandContext(ctx, ABCPATH, "-c", command).CombinedOutput()
	fmt.Println(string(output))
	out := fmt.Sprintln(err)
	fmt.Println(out)
}

func TestMutate(t *testing.T){
	fileIn := "/mnt/d/softwareTesting/HWMCFuzz/benchmark/example/SAT.aig"
	fileOut := "/mnt/d/softwareTesting/HWMCFuzz/Tmp/SATMutate.aig"
	aigmutate(fileIn,fileOut)
}

func TestNuXmv(t *testing.T){
	var bugType = [...]string{"sat", "unsat", "timeout", "crash", "ignore"}
	fileIn := "/mnt/d/softwareTesting/HWMCFuzz/benchmark/example/UNSAT.aig"
	s := UsenuXmvcheck(fileIn)
	println(bugType[s])
}

func TestIIMC(t *testing.T){
	var bugType = [...]string{"sat", "unsat", "timeout", "crash", "ignore"}
	fileIn := "/mnt/d/softwareTesting/HWMCFuzz/benchmark/example/SAT.aig"
	s := UseIIMCcheck(fileIn)
	println(bugType[s])
}

func TestAVY(t *testing.T){
	var bugType = [...]string{"sat", "unsat", "timeout", "crash", "ignore"}
	fileIn := "/mnt/d/softwareTesting/HWMCFuzz/benchmark/example/UNSAT.aig"
	s := AVYcheck(fileIn)
	println(bugType[s])
}


