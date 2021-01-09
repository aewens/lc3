package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/aewens/lc3/vm"
)

type FlagState struct {
	Program chan string
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

func readLines(path string, lines chan string) {
	file, err := os.Open(path)
	catch(err)

	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		lines <- line
	}

	close(lines)

	err = scanner.Err()
	catch(err)
}

func parseFlags() *FlagState {
	state := &FlagState{
		Program: make(chan string),
	}

	fileFlag := flag.String("f", "", "Path to program to run")
	programFlag := flag.String("p", "", "Program to run")

	hasFileFlag := len(*fileFlag) > 0
	hasProgramFlag := len(*programFlag) > 0

	if !hasFileFlag && !hasProgramFlag {
		log.Fatal("[!] Missing -p or -f flag")
		return state
	}

	if hasFileFlag && hasProgramFlag {
		log.Fatal("[!] Cannot use both -p and -f flags together")
		return state
	}

	if hasFileFlag {
		readLines(*fileFlag, state.Program)
	} else if hasProgramFlag {
		lines := strings.Split(*programFlag, "\n")
		for _, line := range lines {
			state.Program <- line
		}
	}

	return state
}

func main() {
	defer cleanup()
	handleSigterm()

	state := parseFlags()
	computer := vm.New()
	output := computer.Run(state.Program)

	log.Print("Output:")
	for _, out := range output {
		log.Printf(" %d", out)
	}
	log.Println()
}
