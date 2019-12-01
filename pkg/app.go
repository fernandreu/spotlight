package app

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type ImageFile struct {
	folder string
	name   string
	format string
	width  int
	height int
}

func (i *ImageFile) FullPath() string {
	return filepath.Join(i.folder, i.name)
}

func (i *ImageFile) Hash() (string, int, error) {
	b, err := ioutil.ReadFile(i.FullPath())
	if err != nil {
		return "", 0, err
	}

	return fmt.Sprintf("%x", md5.Sum(b)), len(b), nil
}

func (i *ImageFile) SaveAs(path string) error {
	input, err := ioutil.ReadFile(i.FullPath())
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, input, 0644)
	if err != nil {
		return err
	}

	return nil
}

func parseImageFile(folder string, name string) *ImageFile {
	fullPath := filepath.Join(folder, name)
	file, err := os.Open(fullPath)
	if err != nil {
		return nil
	}

	defer tryClose(file)

	img, format, err := image.DecodeConfig(file) // format will be "jpeg", "png", etc.
	if err != nil {
		return nil
	}

	if img.Width < 1024 || img.Height < 768 || img.Width <= img.Height {
		return nil
	}

	result := ImageFile{
		folder: folder,
		name:   name,
		format: format,
		width:  img.Width,
		height: img.Height,
	}

	return &result
}

func listPictures(folder string) []ImageFile {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		log.Fatal(err)
	}

	var result []ImageFile
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		file := parseImageFile(folder, f.Name())
		if file != nil {
			result = append(result, *file)
		}
	}

	return result
}

// Reads a file with "name : hash" entries. Returns both a name->hash and a hash->name map
func readCheckSumFile(path string) (map[string]string, map[string]string) {
	filesByHash := make(map[string]string)
	filesByName := make(map[string]string)
	if file, err := os.Open(path); err == nil {
		defer tryClose(file)
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
	defer tryClose(file)

	w := bufio.NewWriter(file)
	defer tryFlush(w)
	for key, value := range filesByName {
		_, _ = fmt.Fprintf(w, "%s : %s\n", key, value)
	}
}

func ProcessFiles(destFolder string) {
	const SubPath = "\\AppData\\Local\\Packages\\Microsoft.Windows.ContentDeliveryManager_cw5n1h2txyewy\\LocalState\\Assets\\"
	originFolder := filepath.Join(os.Getenv("USERPROFILE"), SubPath)

	fmt.Printf("Copying all available Windows Spotlight pictures\n\n")
	fmt.Printf("Origin folder:\n%s\n\n", originFolder)
	fmt.Printf("Destination folder:\n%s\n\n", destFolder)

	// Read the Hash checksum of the existing files from a "CheckSums.txt" file in the same folder,
	// if it exists. File is simply stored as "filename : hash" on each line
	checkSumFile := filepath.Join(destFolder, "CheckSums.txt")
	filesByName, filesByHash := readCheckSumFile(checkSumFile)

	// Scan all files in the destination folder in case there is a mismatch between the stored
	// hashes and the actual files in the folder
	for _, file := range listPictures(destFolder) {
		_, ok := filesByName[file.name]
		if ok {
			continue
		}

		hash, _, err := file.Hash()
		if err != nil {
			continue
		}

		filesByName[file.name] = hash
		filesByHash[hash] = file.name
	}

	// Scan the origin folder now. Make sure the pictures are at least 1024x768, as there might be thumbnails or
	// different versions of the same picture too
	count := 0
	for _, file := range listPictures(originFolder) {

		hash, size, err := file.Hash()
		if err != nil {
			continue
		}

		if _, ok := filesByHash[hash]; !ok {
			destName := file.name + "." + file.format
			err = file.SaveAs(filepath.Join(destFolder, destName))
			if err != nil {
				log.Print(err)
				continue
			}

			fmt.Printf("File copied: %s (%dx%d, %d kb)\n", destName, file.width, file.height, size/1024)
			count += 1

			filesByHash[hash] = destName
			filesByName[destName] = hash
		}
	}

	fmt.Printf("Finished copying; %d new pictures found in total\n", count)

	// Write the full list of files and hashes back to the text file
	writeCheckSumFile(checkSumFile, filesByName)
}

func tryClose(f *os.File) {
	_ = f.Close()
}

func tryFlush(w *bufio.Writer) {
	_ = w.Flush()
}
