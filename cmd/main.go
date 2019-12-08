package main

import (
	"fmt"
	"github.com/fernandreu/spotlight/pkg"
	"os"
	"strings"
)

func main() {
	var files []string
	if len(os.Args) > 1 {
		files = app.ProcessFiles(app.GetDefaultSpotlightFolder(), os.Args[1])
	} else {
		files = app.ProcessFiles(app.GetDefaultSpotlightFolder(), ".\\")
	}

	if len(files) > 0 {
		var answer string
		fmt.Print("Do you want to open the copied pictures (y/[n])? ")
		_, err := fmt.Scanf("%s", &answer)
		if err != nil && strings.ToUpper(answer) == "Y" {
			for _, file := range files {
				app.OpenFile(file)
			}
		}
	} else {
		fmt.Print("No new pictures found. Press any key to exit the program...")
		app.Pause()
	}
}
