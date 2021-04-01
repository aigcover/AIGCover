package main

import (
	"io/ioutil"
	"log"
	"fmt"
	"os/exec"
	"strings"
	"strconv"
	"os"
	"runtime"
	"time"
	"syscall"
)



func main(){
	DiffTestCheckConfigs()
	index := 0
	lastTime := int64(CreateTimeLimit)
	for _,dir := range checkDirs {
		files, _ := ioutil.ReadDir(dir)
		for _,file := range files {
			if(!file.IsDir()&&!strings.HasSuffix(file.Name(), "aig")){
				createTime := GetFileCreateTime(dir+"/"+file.Name())				
				if(createTime>CreateTimeLimit){
					//println(dir+"/"+file.Name())
					if(createTime>lastTime){
						lastTime=createTime
					}
					CopyFile(CaseDir+"/"+strconv.Itoa(index)+".aig",dir+"/"+file.Name())
					index=index+1
				}				
			}			
		}
	}
	files, err := ioutil.ReadDir(CaseDir)
	if err != nil {
		log.Fatal(err)
	}
	for _,f := range files {
		if strings.HasSuffix(f.Name(),".aig") {
			differentialTesting(CaseDir,f.Name())
			//GETVaild(path,f.Name())
		} 		
	}
	println("Final File CreateTime: "+strconv.FormatInt(lastTime,10))
}


func differentialTesting(path string,aigfile string){
	aigfilePath :=path+"/"+aigfile
	pdr,_ := UseSampleABCcheck(aigfilePath, "pdr")
	inter,_ := UseSampleABCcheck(aigfilePath, "int")
	//bmc,_ := UseSampleABCcheck(aigfilePath, "bmc3")
	//tempor,_ := UseSampleABCcheck(aigfilePath, "tempor")
	//dprove,_ := UseSampleABCcheck(aigfilePath, "dprove")
	//IC3,_ := UseIC3checkWithTimeout(aigfilePath)
	NuXmv := UsenuXmvcheck(aigfilePath)
	avy := AVYcheck(aigfilePath)
	//iimc := UseIIMCcheck(aigfilePath)

	//var checkList = [...]int{pdr, inter, tempor, dprove,IC3,NuXmv,avy,iimc}
	//var checkOut = [...]string{pdrOut,intOut,bmcOut,temporOut,dproveOut,IC3Out}
	var checkList = [...]int{pdr, inter, NuXmv,avy}
	//var checkListName = [...]string{"pdr", "inter", "tempor", "dprove","IC3","NuXmv","avy","iimc"}
	var checkListName = [...]string{"pdr", "inter", "NuXmv","avy"}
	var bugType = [...]string{"sat", "unsat", "timeout", "crash", "ignore"}
	println(aigfile)
	for i, checkInfo := range checkList {
		print(checkListName[i], ":", bugType[checkInfo], "   ")
		if checkInfo == CRASH {
			filename := fmt.Sprintf("%s_%s_%s", "crash", checkListName[i],aigfile)
			println(filename)
			CopyFile(bugPath+"/"+filename,aigfilePath)
		}
		for j := range checkList {
			if notequal(checkList[i], checkList[j]) {
				filename := fmt.Sprintf("%s_%s_%s_%s",  "incorrect",checkListName[i],checkListName[j],aigfile)
				println(filename)
				CopyFile(bugPath+"/"+filename,aigfilePath)
				return
			}
		}	
	}
	println()
}

func CopyFile (targetPath string,srcPath string){
	println("copy", srcPath, targetPath)
	output, err :=
		exec.Command("cp", srcPath, targetPath).CombinedOutput()
	if err != nil || string(output) != "" {
		log.Fatal(err, string(output))
	}
}

func GetFileCreateTime(path string) int64{
    osType := runtime.GOOS
    fileInfo, _ := os.Stat(path)
    if osType == "linux" {
        stat_t := fileInfo.Sys().(*syscall.Stat_t)
        tCreate := int64(stat_t.Ctim.Sec)
        return tCreate
    }
    return time.Now().Unix()
}



func GETNontimeout(path string,aigfile string){
	aigfilePath :=path+"/"+aigfile
	pdr,_ := UseABCcheck(aigfilePath, "pdr")
	inter,_ := UseABCcheck(aigfilePath, "int")
	bmc,_ := UseABCcheck(aigfilePath, "bmc2")
	tempor,_ := UseABCcheck(aigfilePath, "tempor")
	dprove,_ := UseABCcheck(aigfilePath, "dprove")
	IC3,_ := UseIC3checkWithTimeout(aigfilePath)

	var checkList = [...]int{pdr, inter, bmc, tempor, dprove,IC3}
	//var checkOut = [...]string{pdrOut,intOut,bmcOut,temporOut,dproveOut,IC3Out}
	var checkListName = [...]string{"pdr", "inter", "bmc", "tempor", "dprove","IC3"}
	var bugType = [...]string{"sat", "unsat", "timeout", "crash", "ignore"}
	println(aigfile)
	count := 0
	for i, checkInfo := range checkList {
		print(checkListName[i], ":", bugType[checkInfo], "   ")
		if checkInfo == TIMEOUT || checkInfo == IGNORE {
			count++
		}
	}
	println(" ")
	if count <= 4 {
		println("copy "+aigfile)
		CopyFile("/mnt/d/softwareTesting/HWMCFuzz/benchmark/benchmarkAIG2/"+aigfile,aigfilePath)
	}	
}

func GETVaild(path string,aigfile string){
	aigfilePath :=path+"/"+aigfile
	pdr,_ := UseABCcheck(aigfilePath, "pdr")
	inter,_ := UseABCcheck(aigfilePath, "int")
	bmc,_ := UseABCcheck(aigfilePath, "bmc2")
	tempor,_ := UseABCcheck(aigfilePath, "tempor")
	dprove,_ := UseABCcheck(aigfilePath, "dprove")

	var checkList = [...]int{pdr, inter, bmc, tempor, dprove}
	//var checkOut = [...]string{pdrOut,intOut,bmcOut,temporOut,dproveOut,IC3Out}
	var checkListName = [...]string{"pdr", "inter", "bmc", "tempor", "dprove"}
	var bugType = [...]string{"sat", "unsat", "timeout", "crash", "ignore"}
	println(aigfile)
	count := 0
	countCrash := 0 
	for i, checkInfo := range checkList {
		print(checkListName[i], ":", bugType[checkInfo], "   ")
		if checkInfo == TIMEOUT || checkInfo == IGNORE {
			count++
		}
		if checkInfo == CRASH { 
			countCrash++
		}
	}
	println(" ")
	if countCrash > 0 {
		CopyFile("/mnt/d/softwareTesting/benchmarkCrash/"+aigfile,aigfilePath)
		return
	}
	if count <= 3 {
		println("copy "+aigfile)
		CopyFile("/mnt/d/softwareTesting/benchmark2AFL/"+aigfile,aigfilePath)
	}	
}