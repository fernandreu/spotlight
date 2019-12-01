package main

import (
	"os"
	"spotlight/pkg"
)

func main() {
	if len(os.Args) > 1 {
		app.ProcessFiles(os.Args[1])
	} else {
		app.ProcessFiles(".\\")
	}

	app.Pause()
}
