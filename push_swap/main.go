package main

import (
	"fmt"
	"os"
	"strings"

	"swap/stack"
)

func main() {
	if len(os.Args) < 2 {
		return
	}

	a, ok := stack.Parse(strings.Join(os.Args[1:], " "))
	if !ok {
		fmt.Fprintln(os.Stderr, "Error")
		os.Exit(1)
	}
	if len(a) == 0 {
		return
	}
	for _, instr := range stack.Sort(a) {
		fmt.Println(instr)
	}
}
