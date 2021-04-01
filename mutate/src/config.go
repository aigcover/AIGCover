package main

import (
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

//const SCRATCHPATH = "/home/zhangche/HWMC/HWMCFuzz/scratch"
const SCRATCHPATH = "/home/aigcover/AIGCover/HWMCFuzzer/scratch"

//const SEEDPATH = "/home/zhangche/HWMC/HWMCFuzz/models/benchmarks/simple_aag/"

const SEEDPATH = "/mnt/d/softwareTesting/HWMCFuzz/benchmark/benchmarkAAG/"
//const SEEDPATH = "/mnt/d/softwareTesting/HWMCFuzz/benchmark/test/"

/*AVYPATH to avy binary*/
//const AVYPATH = "/home/zhangche/HWMC/HWMCFuzz/extavy/bin/avy"
const AVYPATH = "/mnt/d/softwareTesting/HWMCFuzz/AVY/extavy/bin/avy"

/*ABCPATH to abc binary*/
//const ABCPATH = "/home/zhangche/HWMC/abc/abc"
const ABCPATH = "/home/aigcover/AIGCover/abcCov/abc"

/*IC3ref bash path*/
const IC3BashPATH = "/mnt/d/softwareTesting/HWMCFuzz/IC3.sh"

/*IC3ref bash with timeout path*/
const IC3BashWithTimeoutPATH = "/mnt/d/softwareTesting/HWMCFuzz/IC3_Timeout8.sh"

const nuXmvPATH = "/mnt/d/softwareTesting/HWMCFuzz/nuXmv-2.0.0-Linux/bin/nuXmv"
const nuXmvCmdPATH = "/mnt/d/softwareTesting/HWMCFuzz/nuXmv-2.0.0-Linux/bin/cmd"

const iimcPATH = "/mnt/d/softwareTesting/HWMCFuzz/iimc/iimc"

/*AIGTOAIGPATH to abc binary*/
//const AIGTOAIGPATH = "/home/zhangche/HWMC/aiger-1.9.9/aigtoaig"
const AIGTOAIGPATH = "/home/aigcover/AIGCover/aiger-1.9.9/aigtoaig"

/*BUGPATH to folder for saving the bugs*/
//const BUGPATH = "/home/zhangche/HWMC/HWMCFuzz/bugs"
const BUGPATH = "/home/aigcover/AIGCover/bugs"

/*IGNOREINFO error meaasges ignored*/
var IGNOREINFO = [5]string{"Does not work for combinational networks.",
	"The leading sequence has length",
	"Reached local conflict limit",
	"The current network is combinational"}

var RANDID = ""

/*AAGPATH to temp aag file for fuzzing*/
var AAGPATH = ""

/*AIGPATH to temp aig file for fuzzing*/
var AIGPATH = ""

/*IsIgnore return true if the output would like to be ignored*/
func IsIgnore(output []byte) bool {
	for _, info := range IGNOREINFO {
		if strings.Contains(string(output), info) {
			return true
		}
	}
	return false
}

/*CheckConfigs of the fuzzer*/
func CheckConfigs() {

	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 5
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}

	RANDID = b.String()
	AAGPATH = SCRATCHPATH + "/" + RANDID + ".aag"
	AIGPATH = SCRATCHPATH + "/" + RANDID + ".aig"

	if _, err := os.Stat(SEEDPATH); os.IsNotExist(err) {
		log.Fatal("Cannot find seed folder: " + SEEDPATH)
	}

	if _, err := os.Stat(BUGPATH); os.IsNotExist(err) {
		err := os.Mkdir(BUGPATH, 0755)
		if err != nil {
			log.Fatal(err)
		}
	} else if !os.IsNotExist(err) {
		if !isEmpty(BUGPATH) {
			log.Fatal("Bug folder already existed: " + BUGPATH)
		}
	}
}

func isEmpty(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true
	}
	return false
}
