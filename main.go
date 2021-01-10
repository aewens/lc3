package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/aewens/lc3/vm"
)

type FlagState struct {
	Program chan int
}

func cleanup() {
	r := recover()
	if r != nil {
		log.Fatal("[!]: ", r)
	}
}

func handleSigterm() {
	sigterm := make(chan os.Signal)
	signal.Notify(sigterm, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigterm
		cleanup()
	}()
}

func catch(err error) {
	if err != nil {
		panic(err)
	}
}

func any(xs ...bool) bool {
	for _, x := range xs {
		if x {
			return true
		}
	}

	return false
}

func lineToInts(line string, program chan int) {
	for _, field := range strings.Fields(line) {
		value, err := strconv.Atoi(field)
		catch(err)

		program <- value
	}
}

func readLines(path string, program chan int) {
	file, err := os.Open(path)
	catch(err)

	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		lineToInts(line, program)
	}

	err = scanner.Err()
	catch(err)
}

func parseFlags() *FlagState {
	state := &FlagState{
		Program: make(chan int),
	}

	fileFlag := flag.String("f", "", "Path to program to run")
	programFlag := flag.String("p", "", "Program to run")
	rawFlag := flag.String("r", "", "Provide raw instructions to run")

	hasFileFlag := len(*fileFlag) > 0
	hasProgramFlag := len(*programFlag) > 0
	hasRawFlag := len(*rawFlag) > 0

	if !any(hasFileFlag, hasProgramFlag, hasRawFlag) {
		log.Fatal("[!] Missing -p, -f, or -r flag")
		return state
	}

	if any(
		hasFileFlag && hasProgramFlag,
		hasFileFlag && hasRawFlag,
		hasProgramFlag && hasRawFlag,
	) {
		log.Fatal("[!] Cannot use -p, -f, and/or -r flags together")
		return state
	}

	if hasFileFlag {
		readLines(*fileFlag, state.Program)
	} else if hasProgramFlag {
		lines := strings.Split(*programFlag, "\n")
		for _, line := range lines {
			lineToInts(line, state.Program)
		}
	} else if hasRawFlag {
		lineToInts(*rawFlag, state.Program)
	}

	close(state.Program)
	return state
}

func main() {
	defer cleanup()
	handleSigterm()

	state := parseFlags()
	computer := vm.New()
	computer.LoadImage(state.Program)
	computer.Run()
}
