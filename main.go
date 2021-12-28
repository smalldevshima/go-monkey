package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/smalldevshima/go-monkey/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Monkey programming language REPL!\n", user.Username)
	fmt.Printf("Feel free to type in some code!\n")
	repl.Start(os.Stdin, os.Stdout)
}
