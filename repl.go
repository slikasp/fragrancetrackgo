package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func runREPL() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			fmt.Println("\ninput closed, exiting")
			return
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		fmt.Printf("you typed: %s\n", input)
	}
}
