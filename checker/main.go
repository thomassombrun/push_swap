package main

import (
	"bufio"
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

	st := stack.State{A: a}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		instr := scanner.Text()
		if instr == "" {
			continue
		}
		if !st.Exec(instr) {
			fmt.Fprintln(os.Stderr, "Error")
			os.Exit(1)
		}
	}

	if st.IsSorted() {
		fmt.Println("OK")
	} else {
		fmt.Println("KO")
	}
}
