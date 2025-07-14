// This is a dummy program that just echos user input on both stdout and stderr. Used in client testing.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		os.Stdout.WriteString(input + "\n")
		os.Stderr.WriteString(input + "\n")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
	}
}
