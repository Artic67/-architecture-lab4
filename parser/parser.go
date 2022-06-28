package parser

import (
	"fmt"
	"strings"
)

// Interface to commands
type Command interface {
	Execute(handler Handler)
}

// Interface to Handler
type Handler interface {
	Post(cmd Command)
}


// ------------------------ START Commands (Print, Delete)
// Struction for print command
type PrintCmd struct {
	Text string
}

// Method for print command
func (pc *PrintCmd) Execute(handler Handler) {
	fmt.Println(pc.Text)
}

// Struction for delete command
type DeleteCmd struct {
	Text      string
	DSymbol   string
}

// Method for Delete command
func (dc *DeleteCmd) Execute(handler Handler) {
	resultString := strings.ReplaceAll(dc.Text, dc.DSymbol, "")
	handler.Post(&PrintCmd{Text: resultString})
}
// ------------------------ END Commands (Print, Delete)

// Type for functions that makes struct proto from command args
type cmdProc func(args []string) (Command, error)

// Functions that makes struct proto from Print command args
func funcPrint(args []string) (Command, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("wrong number of arguments."+
			"Expected 1, got %d instead", len(args))
	}

	return &PrintCmd{args[0]}, nil
}

// Functions that makes struct proto from Delete command args
func funcDelete(args []string) (Command, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("Wrong number of args. "+
			"Expected 2, got %d", len(args))
	}

	if len(args[1]) != 1 {
		return nil, fmt.Errorf("Deletion symbol length should be 1 character long!")
	}

	return &DeleteCmd{args[0], args[1]}, nil
}

// List of existing comands
var cmdList = map[string]cmdProc{
	"print": funcPrint,
	"delete": funcDelete,
}

// Function to find command function to give args
func findCommand(commandStr string) (cmdProc, error) {
	for cmd, fn := range cmdList {
		if cmd == commandStr {
			return fn, nil
		}
	}

	return nil, fmt.Errorf("Unknown command")
}

// Main Parse funcion
func Parse(text string) Command {
	values := strings.Split(text, " ")
	errPrefix := "SYNTAX ERROR: "
	cmdFunc, error := findCommand(values[0])
	if error != nil {
		return &PrintCmd{
			Text: errPrefix + error.Error(),
		}
	}

	command, error := cmdFunc(values[1:])
	if error != nil {
		return &PrintCmd{
			Text: errPrefix + error.Error(),
		}
	}

	return command
}
