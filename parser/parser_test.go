package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Error string start words
var errorPrefix = "SYNTAX ERROR:"


// Test Delete command where result is error ( Too long separator )
func TestDeleteCommandTooLongSepError(t *testing.T) {

	assert := assert.New(t)

	inputString := "delete Str...err ..."

	command := Parse(inputString)
	assert.Contains(command.(*PrintCmd).Text, errorPrefix)
}

// Test Print command where result is error ( Too many arguments )
func TestPrintCommandCountArgumentsError(t *testing.T) {
	assert := assert.New(t)
	inputStr := "print too many args"

	command := Parse(inputStr)

	assert.Contains(command.(*PrintCmd).Text, errorPrefix)
}

// Test Delete command where result is error ( Too many arguments )
func TestDeleteCommandCountArgumentsError(t *testing.T) {
	assert := assert.New(t)
	inputStr := "delete too many args"

	command := Parse(inputStr)

	assert.Contains(command.(*PrintCmd).Text, errorPrefix)
}

// Test error when command is unknown
func TestUnknownCommandError(t *testing.T) {
	assert := assert.New(t)
	inputStr := "qwery 2288"

	command := Parse(inputStr)

	assert.Contains(command.(*PrintCmd).Text, errorPrefix)
}

// Test Delete command where result is succesful
func TestDeleteCommandSuccessful(t *testing.T) {

	assert := assert.New(t)

	inputString := "delete man!delete!me!please !"
	expectedString := "man!delete!me!please"
	expectedDelSym := "!"

	command := Parse(inputString)
	assert.Equal(command.(*DeleteCmd).Text, expectedString)
	assert.Equal(command.(*DeleteCmd).DSymbol, expectedDelSym)
}

// Test Print command where result is successful
func TestPrintCommandSuccessful(t *testing.T) {

	assert := assert.New(t)

	inputString := "print hey!"
	expectedString := "hey!"

	command := Parse(inputString)
	assert.Equal(command.(*PrintCmd).Text, expectedString)
}


// Example of parsing with two commands
func Example() {
	inputPrint := "print your-string-here"
	inputDelete := "delete"

	commandPrint := Parse(inputPrint)
	commandDelete := Parse(inputDelete)

	var handler Handler
	handler.Post(commandPrint)
	handler.Post(commandDelete)
}
