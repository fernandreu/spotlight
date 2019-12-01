package main

import (
	"bufio"
	"fmt"
	"os"
	"spotlight/pkg"
)

func main() {
	if len(os.Args) > 1 {
		app.ProcessFiles(os.Args[1])
	} else {
		app.ProcessFiles(".\\")
	}

	fmt.Print("Press 'Enter' to continue...")
	_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
}
