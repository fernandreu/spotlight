package main

import (
	"fmt"
	"github.com/fernandreu/spotlight/pkg"
	"github.com/fernandreu/spotlight/pkg/gui"
	"strings"
)

func main() {
	options := app.ParseFlags()
	if options.ShowGui {
		gui.LaunchGui(options)
		return
	}

	result := app.ProcessFiles(options.Source, options.Destination)
	if !options.Prompt {
		return
	}

	if len(result.NewFiles) > 0 {
		fmt.Print("Do you want to open the copied pictures (y/[n])? ")
		var answer string
		_, err := fmt.Scanf("%s", &answer)
		if err == nil && strings.ToUpper(answer) == "Y" {
			for _, file := range result.NewFiles {
				app.OpenFile(file.FullPath())
			}
		}
	} else {
		fmt.Print("No new pictures found. Press any key to exit the program...")
		app.Pause()
	}
}
