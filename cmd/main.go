package main

import (
	"fmt"
	"github.com/fernandreu/spotlight/pkg"
	"os"
)

func main() {
	var files []string
	if len(os.Args) > 1 {
		files = app.ProcessFiles(app.GetDefaultSpotlightFolder(), os.Args[1])
	} else {
		files = app.ProcessFiles(app.GetDefaultSpotlightFolder(), ".\\")
	}

	if len(files) > 0 {
		for _, file := range files {
			app.OpenFile(file)
		}
	}

	fmt.Print("Press any key to exit the program...")
	app.Pause()
}
