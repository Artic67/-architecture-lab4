package loop

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/Artic67/architecture-lab4/parser"
	"github.com/stretchr/testify/assert"
)

// Function that takes out string
func captureOutString() func() (string, error) {

	r, w, err := os.Pipe()

	// Check error
	if err != nil {
		panic(err)
	}

	made := make(chan error, 1)

	saveStdOut := os.Stdout
	os.Stdout = w

	var buf strings.Builder

	go func() {
		_, error := io.Copy(&buf, r)
		r.Close()
		made <- error
	}()

	return func() (string, error) {
		os.Stdout = saveStdOut
		w.Close()
		error := <-made
		return buf.String(), error
	}
}

// Function to test loop for multiple awaits
func TestManyAwaitsAtOneTime(t *testing.T) {

	loop := new(EventLoop)
	loop.Start()

	loop.Post(&parser.PrintCmd{Text: "?@!/,. s ddd!"})
	loop.Post(&parser.PrintCmd{Text: "I am here!"})
	loop.Post(&parser.DeleteCmd{
		Text:    "M E R G E D W O R D ?",
		DSymbol: " ",
	})

	go loop.AwaitClose()
	go loop.AwaitClose()
	go loop.AwaitClose()

	loop.AwaitClose()

}

// Function to test post command after loop was stopped
func TestPostCommandAfterLoopWasStopped(t *testing.T) {

	loop := new(EventLoop)
	loop.Start()
	loop.Post(&parser.PrintCmd{Text: "Yeah I am going forward into event loop :)"})

	loop.AwaitClose()

	loop.Post(&parser.PrintCmd{Text: "Oh no, I will be executed in stopped loop!"})
	loop.Post(&parser.PrintCmd{Text: "Me too, bro :("})
	loop.Post(&parser.DeleteCmd{
		Text:    "Yall stupid go home!",
		DSymbol: " ",
	})

}

// Function to test Print command input
func TestExecutionOfPrintCommand(t *testing.T) {

	assert := assert.New(t)
	loop := new(EventLoop)

	input := "TEST MESSAGE"
	expected := "TEST MESSAGE\n"
	command := parser.PrintCmd{Text: input}

	loop.Start()
	getString := captureOutString()

	loop.Post(&command)
	loop.AwaitClose()

	capOut, err := getString()
	if err != nil {
		panic(err)
	}

	assert.Equal(capOut, expected)
}

// Function to test Delete command input
func TestExecutionOfDeleteCommand(t *testing.T) {

	assert := assert.New(t)
	loop := new(EventLoop)

	inputString := "delete,me,please"
	inputDelSym := ","
	expected := "deletemeplease\n"

	command := parser.DeleteCmd{
		Text:    inputString,
		DSymbol: inputDelSym,
	}

	loop.Start()
	getStr := captureOutString()

	loop.Post(&command)
	loop.AwaitClose()

	capturedOutput, err := getStr()

	if err != nil {
		panic(err)
	}

	assert.Equal(capturedOutput, expected)
}

// Function to test stoping of empty loop
func TestStopEmptyLoop(t *testing.T) {

	loop := new(EventLoop)
	loop.Start()
	loop.AwaitClose()

}

// Example function of event loop
func ExampleEL() {

	loop := new(EventLoop)

	loop.Start()
	loop.Post(&parser.PrintCmd{Text: "Hello, I am event loop!"})
	loop.AwaitClose()

	// Output:
	// Hello, I am event loop!
}
