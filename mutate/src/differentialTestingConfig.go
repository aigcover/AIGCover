package main

import (
	"log"
	"math/rand"
	"os"
	"time"
)

var checkDirs = [...]string{"/mnt/d/softwareTesting/AFLoutPDR/queue",
						"/mnt/d/softwareTesting/AFLoutINT/queue",
						//"/mnt/d/softwareTesting/AFLoutTEMPOR/queue",
						"/mnt/d/softwareTesting/AVYout/queue",
						"/mnt/d/softwareTesting/AVYout/crashes"}
const CaseDir = "/mnt/d/softwareTesting/TestCase"

const bugPath = "/mnt/d/softwareTesting/bugs4"

//const CreateTimeLimit = 1610344098
const CreateTimeLimit = 1610946566



func DiffTestCheckConfigs() {

	rand.Seed(time.Now().UnixNano())
	for _,dir := range checkDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			log.Fatal("Cannot find folder: " +dir)
		}
	}

	if _, err := os.Stat(CaseDir); os.IsNotExist(err) {
		err := os.Mkdir(CaseDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	} else if !os.IsNotExist(err){
		if !isEmpty(CaseDir) {
			log.Fatal("Case folder already existed: " + CaseDir)
		}
	}

	if _, err := os.Stat(bugPath); os.IsNotExist(err) {
		err := os.Mkdir(bugPath, 0755)
		if err != nil {
			log.Fatal(err)
		}
	} else if !os.IsNotExist(err) {
		if !isEmpty(bugPath) {
			log.Fatal("Bug folder already existed: " + bugPath)
		}
	}

	
}