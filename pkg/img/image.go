package img

import (
	"crypto/md5"
	"fmt"
	"github.com/fernandreu/spotlight/pkg/common"
	"image"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type ImageFile struct {
	Folder string
	Name   string
	Format string
	Width  int
	Height int
}

func (i *ImageFile) clone() ImageFile {
	result := ImageFile{
		Folder: i.Folder,
		Name:   i.Name,
		Format: i.Format,
		Width:  i.Width,
		Height: i.Height,
	}
	return result
}

func (i *ImageFile) FullPath() string {
	return filepath.Join(i.Folder, i.Name)
}

func (i *ImageFile) Hash() (string, int, error) {
	b, err := ioutil.ReadFile(i.FullPath())
	if err != nil {
		return "", 0, err
	}

	return fmt.Sprintf("%x", md5.Sum(b)), len(b), nil
}

func (i *ImageFile) SaveAs(folder string, name string) (*ImageFile, error) {
	input, err := ioutil.ReadFile(i.FullPath())
	if err != nil {
		return nil, err
	}

	path := filepath.Join(folder, name)
	err = ioutil.WriteFile(path, input, 0644)
	if err != nil {
		return nil, err
	}

	result := i.clone()
	result.Folder = folder
	result.Name = name
	return &result, nil
}

func (i *ImageFile) Rename(name string) error {
	err := os.Rename(i.FullPath(), filepath.Join(i.Folder, name))
	if err != nil {
		return err
	}

	i.Name = name
	return nil
}

func Parse(folder string, name string) *ImageFile {
	fullPath := filepath.Join(folder, name)
	file, err := os.Open(fullPath)
	if err != nil {
		return nil
	}

	defer common.TryCloseFile(file)

	img, format, err := image.DecodeConfig(file) // Format will be "jpeg", "png", etc.
	if err != nil {
		return nil
	}

	if img.Width < 1024 || img.Height < 768 || img.Width <= img.Height {
		return nil
	}

	result := ImageFile{
		Folder: folder,
		Name:   name,
		Format: format,
		Width:  img.Width,
		Height: img.Height,
	}

	return &result
}

func FindAll(folder string) []ImageFile {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		log.Fatal(err)
	}

	var result []ImageFile
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		file := Parse(folder, f.Name())
		if file != nil {
			result = append(result, *file)
		}
	}

	return result
}
