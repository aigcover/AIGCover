package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
	"os"
	"io"
)

const (
	/*SAT found counterexample*/
	SAT = iota
	/*UNSAT model is safe*/
	UNSAT = iota
	/*TIMEOUT checking timeout*/
	TIMEOUT = iota
	/*CRASH checker crash*/
	CRASH = iota
	/*IGNORE checker crash*/
	IGNORE = iota
)

/*ABCcheck check aig network file by abc*/
func ABCcheck(file string) int {
	command := fmt.Sprintf("read %s;pdr -g;write_cex -a out.cex", file)

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	output, err :=
		exec.CommandContext(ctx, ABCPATH, "-c", command).CombinedOutput()

	switch {
	case ctx.Err() == context.DeadlineExceeded:
		return TIMEOUT
	case err != nil:
		if !IsIgnore(output) {
			return CRASH
		}
		fmt.Println("Ignore crash")
		return IGNORE
	case strings.Contains(string(output), "asserted"):
		return SAT
	case strings.Contains(string(output), "successful"):
		return UNSAT
	default:
		return CRASH
	}
}

/*AVYcheck check aig network file by abc*/
func AVYcheck(file string) int {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	output, err :=
		exec.CommandContext(ctx, AVYPATH, file).CombinedOutput()

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

/*PDRABCcheck check aig network file by abc with */
func PDRABCcheck(file string) int {
	command := fmt.Sprintf("read %s;logic;undc;strash;zero;fold;pdr;", file)

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	output, err :=
		exec.CommandContext(ctx, ABCPATH, "-c", command).CombinedOutput()

	switch {
	case ctx.Err() == context.DeadlineExceeded:
		return TIMEOUT
	case err != nil:
		if !IsIgnore(output) {
			return CRASH
		}
		fmt.Println("Ignore crash")
		return IGNORE
	case strings.Contains(string(output), "asserted"):
		return SAT
	case strings.Contains(string(output), "Property proved"):
		return UNSAT
	default:
		if !IsIgnore(output) {
			return CRASH
		}
		return IGNORE
	}
}

/*BMCABCcheck check aig network file by abc with safe commands sequence suggested by developer*/
func BMCABCcheck(file string) int {
	command := fmt.Sprintf("read %s;logic;undc;strash;zero;fold;bmc2;", file)

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	output, err :=
		exec.CommandContext(ctx, ABCPATH, "-c", command).CombinedOutput()

	switch {
	case ctx.Err() == context.DeadlineExceeded:
		return TIMEOUT
	case err != nil:
		if !IsIgnore(output) {
			return CRASH
		}
		fmt.Println("Ignore crash")
		return IGNORE
	case strings.Contains(string(output), "asserted"):
		return SAT
	case strings.Contains(string(output), "Property proved"):
		return UNSAT
	default:
		if !IsIgnore(output) {
			return CRASH
		}
		return IGNORE
	}
}

/*INTABCcheck check aig network file by abc with safe commands sequence suggested by developer*/
func INTABCcheck(file string) int {
	command := fmt.Sprintf("read %s;logic;undc;strash;zero;fold;int;", file)

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	output, err :=
		exec.CommandContext(ctx, ABCPATH, "-c", command).CombinedOutput()

	switch {
	case ctx.Err() == context.DeadlineExceeded:
		return TIMEOUT
	case err != nil:
		if !IsIgnore(output) {
			return CRASH
		}
		fmt.Println("Ignore crash")
		return IGNORE
	case strings.Contains(string(output), "asserted"):
		return SAT
	case strings.Contains(string(output), "Property proved"):
		return UNSAT
	default:
		if !IsIgnore(output) {
			fmt.Println(string(output))
			return CRASH
		}
		return IGNORE
	}
}
/*
func UseABCcheck(file string, checker string) int {
	command := fmt.Sprintf("read %s;logic;undc;strash;zero;fold;%s;", file, checker)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	output, err :=
		exec.CommandContext(ctx, ABCPATH, "-c", command).CombinedOutput()

	switch {
	case ctx.Err() == context.DeadlineExceeded:
		return TIMEOUT
	case err != nil:
		if !IsIgnore(output) {
			fmt.Println(string(output))
			return CRASH
		}
		fmt.Println("Ignore crash")
		return IGNORE
	case strings.Contains(string(output), "Property proved") || strings.Contains(string(output), "UNSATISFIABLE"):
		return UNSAT
	case strings.Contains(string(output), "asserted") || strings.Contains(string(output), "SATISFIABLE"):
		return SAT
	default:
		if !IsIgnore(output) {
			fmt.Println(string(output))
			return CRASH
		}
		return IGNORE
	}
}
*/
func UseABCcheck(file string, checker string) (int,string) {
	command := fmt.Sprintf("read %s;logic;undc;strash;zero;fold;%s;", file, checker)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	output, err :=
		exec.CommandContext(ctx, ABCPATH, "-c", command).CombinedOutput()
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
	case strings.Contains(string(output), "Property proved") || strings.Contains(string(output), "UNSATISFIABLE"):
		return UNSAT,outMsg
	case strings.Contains(string(output), "asserted") || strings.Contains(string(output), "SATISFIABLE"):
		return SAT,outMsg
	default:
		if !IsIgnore(output) {
			fmt.Println(string(output))
			return CRASH,outMsg
		}
		return IGNORE,outMsg
	}
}

func UseSampleABCcheck(file string, checker string) (int,string) {
	command := fmt.Sprintf("read %s;%s;", file, checker)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	output, err :=
		exec.CommandContext(ctx, ABCPATH, "-c", command).CombinedOutput()
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


func UseIC3checkWithTimeout(file string) (int,string) {
	
	cmd := exec.Command(IC3BashWithTimeoutPATH, file)	
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


func UseIC3check(file string) int {
	
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	output, err :=
		exec.CommandContext(ctx, IC3BashPATH,file).CombinedOutput()
		
	switch {
	case ctx.Err() == context.DeadlineExceeded:
		return TIMEOUT	
	case string(output)=="0\n":
		return UNSAT
	case string(output)=="1\n":
		return SAT
	case err != nil:		
		if !IsIgnore(output) {
			fmt.Println(string(output))
			return CRASH
		}
		fmt.Println("Ignore crash")
		return IGNORE	
	default:
		if !IsIgnore(output) {
			fmt.Println(string(output))
			return CRASH
		}
		return IGNORE
	}			
}


func UsenuXmvcheck(file string) int {

	f, err0 := os.OpenFile(nuXmvCmdPATH, os.O_RDWR | os.O_TRUNC | os.O_CREATE, 0666)
    if err0 != nil {
        fmt.Println(err0.Error())
    }
	_, err1 := io.WriteString(f, "read_aiger_model -i "+file+"\ngo\ncheck_property\nquit\n")
	if err1 != nil {
        fmt.Println(err1.Error())
    }
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	output, err :=
		exec.CommandContext(ctx, nuXmvPATH, "-source", nuXmvCmdPATH).CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return TIMEOUT
	}
	if err != nil {
		return CRASH
	}	
	outs := strings.Split(string(output),"\n")
	for _,line := range outs  {
		if strings.Contains(line,"invariant") {
			if strings.Contains(line,"true") {
				return UNSAT
			} else {
				return SAT
			}
		}
	}
	
	return IGNORE
}

func UseIIMCcheck(file string) int {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	output, err :=
		exec.CommandContext(ctx, iimcPATH, file).CombinedOutput()
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