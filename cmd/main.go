package main

import (
	"github.com/fernandreu/spotlight/pkg"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		app.ProcessFiles(app.GetDefaultSpotlightFolder(), os.Args[1])
	} else {
		app.ProcessFiles(app.GetDefaultSpotlightFolder(), ".\\")
	}

	app.Pause()
}
