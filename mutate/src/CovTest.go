package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os/exec"
	"os"
	"time"
	"context"
	"strings"
)

func main() {
	
	rand.Seed(time.Now().UnixNano())
	path:=os.Args[1]

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	len := len(files)
	for {
		file := files[rand.Intn(len)]
		for i := 0; i < 10; i++ {
			println("Fuzzing:", path +"/"+ file.Name(), i)
			FuzzCovTest(path +"/"+ file.Name(),os.Args[2])
		}
	}
}

func FuzzCovTest(fileIn string,checker string) {

	t := fmt.Sprintf("%d_%v",os.Getpid(),time.Now().UnixNano())
	fileInAAG := "/mnt/d/softwareTesting/HWMCFuzz/Tmp/tempin"+string(t)+".aag"
	fileOutAAG := "/mnt/d/softwareTesting/HWMCFuzz/Tmp/tempout"+string(t)+".aag"
	fileOut := "/mnt/d/softwareTesting/HWMCFuzz/Tmp/TempOut/fuzz"+string(t)+".aig"
	output, err :=
		exec.Command(AIGTOAIGPATH, fileIn, fileInAAG).CombinedOutput()
	if err != nil || string(output) != "" {
		log.Fatal(err, string(output))
	}
	network := Network{}
	network.Read(fileInAAG)
	network.mutation()
	network.Dumpto(fileOutAAG)

	output2, err2 :=
		exec.Command(AIGTOAIGPATH, fileOutAAG, fileOut).CombinedOutput()
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

	res := 4
	if checker=="IC3" {
		IC3,_  := UseIC3check2Cov(fileOut)
		res=IC3
	}else if checker=="pdr" {
		pdr,_ := UseSampleABCcheck2Cov(fileOut, "pdr")
		res=pdr
	}else if checker=="bmc" {
		bmc,_ := UseSampleABCcheck2Cov(fileOut, "bmc3")
		res=bmc
	}else if checker=="int" {
		inter,_ := UseSampleABCcheck2Cov(fileOut, "int")
		res=inter
	}else if checker=="dprove" {
		dprove,_ := UseSampleABCcheck2Cov(fileOut, "dprove")
		res=dprove
	}else if checker=="tem or" {
		tempor,_ := UseSampleABCcheck2Cov(fileOut, "tempor")
		res=tempor
	}else if checker=="iimc" {
		res = UseIIMCcheck2Cov(fileOut)
	}else if checker=="avy"{
		res = UseAVYcheck2Cov(fileOut)
	}
	var bugType = [...]string{"sat", "unsat", "timeout", "crash", "ignore"}
	println(bugType[res])
	err4 := os.Remove(fileOut)
	if err4 != nil {
		log.Fatal(err4,"Delete temp file fail")
	}
	
}

func UseAVYcheck2Cov(file string) int {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	output, err :=
		exec.CommandContext(ctx, "/mnt/d/softwareTesting/HWMCFuzz/AVYCov/AVY2GCov/extavy/build/avy/src/avy", file).CombinedOutput()

	switch {
	case ctx.Err() == context.DeadlineExceeded:
		return TIMEOUT
	case err != nil && err.Error() == "exit status 1":
		fmt.Println(string(output))
		return SAT
	case err != nil:
		if !IsIgnore(output) {
			return CRASH
		}
		fmt.Println("Ignore crash")
		return IGNORE
	case strings.Contains(string(output), "0\nb0\n."):
		return UNSAT
	default:
		return CRASH
	}
}

func UseIC3check2Cov(file string) (int,string) {
	
	cmd := exec.Command("/mnt/d/softwareTesting/HWMCFuzz/IC3_Cov.sh", file)	
	output, err := cmd.CombinedOutput()
	timeout := fmt.Sprintf("%s",err);
	outMsg := string(output)	
	switch {
	case timeout=="exit status 124":
		return TIMEOUT,outMsg	
	case string(output)=="0\n":
		return UNSAT,outMsg
	case string(output)=="1\n":
		return SAT,outMsg
	case err != nil:		
		return CRASH,outMsg	
	default:
		return IGNORE,outMsg
	}			
}

func UseSampleABCcheck2Cov(file string, checker string) (int,string) {
	command := fmt.Sprintf("read %s;%s;", file, checker)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	output, err :=
		exec.CommandContext(ctx, "/mnt/d/softwareTesting/HWMCFuzz/abcGovTest/abc", "-c", command).CombinedOutput()
	outMsg := fmt.Sprintln(err)
	switch {
	case ctx.Err() == context.DeadlineExceeded:
		return TIMEOUT,outMsg
	case err != nil:
		if !IsIgnore(output) {
			if !strings.Contains(outMsg,"cannot allocate memory"){
				fmt.Println(string(output))
				return CRASH,outMsg
			}			
		}
		fmt.Println("Ignore crash")
		return IGNORE,outMsg
	case strings.Contains(string(output), "asserted") || strings.Contains(string(output), "SATISFIABLE"):
		return SAT,outMsg
	case strings.Contains(string(output), "Property proved") || strings.Contains(string(output), "UNSATISFIABLE") || strings.Contains(string(output),"No output failed"):
		return UNSAT,outMsg
	case strings.Contains(string(output),"Reading AIGER files with liveness properties is currently not supported"):
		return IGNORE,outMsg
	default:
		if !IsIgnore(output) {
			fmt.Println(string(output))
			return CRASH,outMsg
		}
		return IGNORE,outMsg
	}
}

func UseIIMCcheck2Cov(file string) int {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	output, err :=
		exec.CommandContext(ctx, "/mnt/d/softwareTesting/HWMCFuzz/iimcCovTest/iimc", file).CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return TIMEOUT
	}
	if err != nil {
		return CRASH
	}	
	outs := strings.Split(string(output),"\n")
	for _,line := range outs  {
		if strings.Compare(line,"0")==0 {
			return UNSAT
		} else if(strings.Compare(line,"1")==0){
			return SAT
		}
	}	
	return IGNORE
}


