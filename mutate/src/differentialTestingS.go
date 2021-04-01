package main

import (
	"os"
	"io/ioutil"
	"log"
	"fmt"
	"os/exec"
	"strings"
)



func main(){
	path:=os.Args[1]
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _,f := range files {
		if strings.HasSuffix(f.Name(),".aig") {
			getBMC3(path,f.Name())
			//sampleDifferentialTesting(path,f.Name())
			//GETVaild(path,f.Name())
		} 		
	}
	
}

func getBMC3(path string,aigfile string){
	aigfilePath :=path+"/"+aigfile
	pdr,_ := UseSampleABCcheck(aigfilePath, "pdr")
	bmc,_ := UseSampleABCcheck(aigfilePath, "bmc3")
	var bugType = [...]string{"sat", "unsat", "timeout", "crash", "ignore"}
	println(aigfile)
	println("pdr:"+bugType[pdr]+" bmc3:"+bugType[bmc])
	if pdr == SAT && (bmc == CRASH || bmc == TIMEOUT) {
		println("find"+aigfilePath)
		CopyFile(BUGPATH+"/"+aigfile,aigfilePath)
	}
	
	println()		
}


func sampleDifferentialTesting(path string,aigfile string){
	aigfilePath :=path+"/"+aigfile
	pdr,_ := UseSampleABCcheck(aigfilePath, "pdr")
	inter,_ := UseSampleABCcheck(aigfilePath, "int")
	//bmc,_ := UseSampleABCcheck(aigfilePath, "bmc3")
	tempor,_ := UseSampleABCcheck(aigfilePath, "tempor")
	dprove,_ := UseSampleABCcheck(aigfilePath, "dprove")
	IC3,_ := UseIC3checkWithTimeout(aigfilePath)
	NuXmv := UsenuXmvcheck(aigfilePath)
	avy := AVYcheck(aigfilePath)
	iimc := UseIIMCcheck(aigfilePath)

	var checkList = [...]int{pdr, inter, tempor, dprove,IC3,NuXmv,avy,iimc}
	//var checkOut = [...]string{pdrOut,intOut,bmcOut,temporOut,dproveOut,IC3Out}
	var checkListName = [...]string{"pdr", "inter", "tempor", "dprove","IC3","NuXmv","avy","iimc"}
	var bugType = [...]string{"sat", "unsat", "timeout", "crash", "ignore"}
	println(aigfile+"(sample cmd line)")
	for i, checkInfo := range checkList {
		print(checkListName[i], ":", bugType[checkInfo], "   ")
		if checkInfo == CRASH {
			filename := fmt.Sprintf("%s_%s_%s.aig", aigfile, "crash", checkListName[i])
			println(filename)
			CopyFile(BUGPATH+"/"+filename,aigfilePath)
		}
		for j := range checkList {
			if notequal(checkList[i], checkList[j]) {
				filename := fmt.Sprintf("%s_%s_%s_%s.aig", aigfile, "incorrect",checkListName[i],checkListName[j])
				println(filename)
				CopyFile(BUGPATH+"/"+filename,aigfilePath)
				return
			}
		}	
	}
	println()		
}

func sampleDifferentialTesting2(path string,aigfile string){
	aigfilePath :=path+"/"+aigfile
	pdr,_ := UseSampleABCcheck(aigfilePath, "pdr")
	inter,_ := UseSampleABCcheck(aigfilePath, "int")
	bmc,_ := UseSampleABCcheck(aigfilePath, "bmc2")
	tempor,_ := UseSampleABCcheck(aigfilePath, "tempor")
	dprove,_ := UseSampleABCcheck(aigfilePath, "dprove")
	IC3,_ := UseIC3checkWithTimeout(aigfilePath)

	var checkList = [...]int{pdr, inter, bmc, tempor, dprove,IC3}
	//var checkOut = [...]string{pdrOut,intOut,bmcOut,temporOut,dproveOut,IC3Out}
	var checkListName = [...]string{"pdr", "inter", "bmc", "tempor", "dprove","IC3"}
	var bugType = [...]string{"sat", "unsat", "timeout", "crash", "ignore"}
	println(aigfile+"(sample cmd line)")
	count := 0
	for i, checkInfo := range checkList {
		print(checkListName[i], ":", bugType[checkInfo], "   ")
		if checkInfo == CRASH {
			filename := fmt.Sprintf("%s_%s_%s.aig", aigfile, "crash", checkListName[i])
			println(filename)
			CopyFile(BUGPATH+"/"+filename,aigfilePath)
		}
		for j := range checkList {
			if notequal(checkList[i], checkList[j]) {
				filename := fmt.Sprintf("%s_%s.aig", aigfile, "incorrect")
				println(filename)
				CopyFile(BUGPATH+"/"+filename,aigfilePath)
				return
			}
		}
		if checkInfo == TIMEOUT || checkInfo == IGNORE {
			count++
		}
		
	}
	if count <= 4 {
		println("\ncopy "+aigfile)
		CopyFile("/mnt/d/softwareTesting/AFLbenchmark2Sample/"+aigfile,aigfilePath)
	}else{
		println("")
	}		
}

func CopyFile (targetPath string,srcPath string){
	output, err :=
		exec.Command("cp", srcPath, targetPath).CombinedOutput()
	if err != nil || string(output) != "" {
		log.Fatal(err, string(output))
	}
}

