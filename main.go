package main

import (
	"bufio"
	"flag"
	"log"
	"os"

	parser "github.com/Artic67/al4/parser"
	loop "github.com/Artic67/al4/loop"
)

// File name argument input
var examplesFileName = flag.String("f", "./src/examples.txt",
	"Write a file name with commands ")

// Main functions
func main() {

	flag.Parse()
	file, error := os.Open(*examplesFileName)

	if error != nil {
		log.Fatal(error)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	el := new(loop.EventLoop)
	el.Start()

	for scanner.Scan() {
		commandLine := scanner.Text()
		cmd := parser.Parse(commandLine)
		el.Post(cmd)
	}

	el.AwaitClose()
}
