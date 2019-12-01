package app

import (
	"bufio"
	"fmt"
	"github.com/fernandreu/spotlight/pkg/common"
	"github.com/fernandreu/spotlight/pkg/img"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Reads a file with "name : hash" entries. Returns both a name->hash and a hash->name map
func readCheckSumFile(path string) (map[string]string, map[string]string) {
	filesByHash := make(map[string]string)
	filesByName := make(map[string]string)
	if file, err := os.Open(path); err == nil {
		defer common.TryCloseFile(file)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			parts := strings.Split(scanner.Text(), ":")
			if len(parts) < 2 {
				continue
			}
			for i := range parts {
				parts[i] = strings.TrimSpace(parts[i])
			}
			filesByName[parts[0]] = parts[1]
			filesByHash[parts[1]] = parts[0]
		}
	}
	return filesByName, filesByHash
}

func writeCheckSumFile(path string, filesByName map[string]string) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer common.TryCloseFile(file)

	w := bufio.NewWriter(file)
	defer common.TryFlush(w)
	for key, value := range filesByName {
		_, _ = fmt.Fprintf(w, "%s : %s\n", key, value)
	}
}

func GetDefaultSpotlightFolder() string {
	const SubPath = "\\AppData\\Local\\Packages\\Microsoft.Windows.ContentDeliveryManager_cw5n1h2txyewy\\LocalState\\Assets\\"
	return filepath.Join(os.Getenv("USERPROFILE"), SubPath)
}

func ProcessFiles(originFolder string, destFolder string) {
	fmt.Printf("Copying all available Windows Spotlight pictures\n\n")
	fmt.Printf("Origin folder:\n%s\n\n", originFolder)
	fmt.Printf("Destination folder:\n%s\n\n", destFolder)

	// Read the Hash checksum of the existing files from a "CheckSums.txt" file in the same folder,
	// if it exists. File is simply stored as "filename : hash" on each line
	checkSumFile := filepath.Join(destFolder, "CheckSums.txt")
	filesByName, filesByHash := readCheckSumFile(checkSumFile)

	// Scan all files in the destination folder in case there is a mismatch between the stored
	// hashes and the actual files in the folder
	for _, file := range img.FindAll(destFolder) {
		_, ok := filesByName[file.Name]
		if ok {
			continue
		}

		hash, _, err := file.Hash()
		if err != nil {
			continue
		}

		filesByName[file.Name] = hash
		filesByHash[hash] = file.Name
	}

	// Scan the origin folder now. Make sure the pictures are at least 1024x768, as there might be thumbnails or
	// different versions of the same picture too
	count := 0
	for _, file := range img.FindAll(originFolder) {

		hash, size, err := file.Hash()
		if err != nil {
			continue
		}

		if _, ok := filesByHash[hash]; !ok {
			destName := file.Name + "." + file.Format
			err = file.SaveAs(filepath.Join(destFolder, destName))
			if err != nil {
				log.Print(err)
				continue
			}

			fmt.Printf("File copied: %s (%dx%d, %d kb)\n", destName, file.Width, file.Height, size/1024)
			count += 1

			filesByHash[hash] = destName
			filesByName[destName] = hash
		}
	}

	fmt.Printf("Finished copying; %d new pictures found in total\n", count)

	// Write the full list of files and hashes back to the text file
	writeCheckSumFile(checkSumFile, filesByName)
}

/**
Uses the standard pause command in Windows to wait for any key. This is obviously specific to Windows, but the entire
program is anyway
*/
func Pause() {
	proc := exec.Command("cmd", "/C", "pause")
	proc.Stdin = os.Stdin
	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr
	err := proc.Start()
	if err != nil {
		log.Fatal(err)
	}
}
