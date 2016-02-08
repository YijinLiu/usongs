package main

import(
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func WriteFile(fileName string, content []byte, mode os.FileMode) error {
	folder := filepath.Dir(fileName)
	// Create the folder if needed.
	if err := os.MkdirAll(folder, os.FileMode(mode|0711)); err != nil {
		return err
	}
	if err := ioutil.WriteFile(fileName, content, mode); err != nil {
		return err
	}
	return nil
}

func CopyFile(src, dst string) error {
	fi, err := os.Stat(src)
	if err != nil {
		return err
	}
	log.Printf("Copying '%s' to '%s' ...", src, dst)
	if content, err := ioutil.ReadFile(src); err != nil {
		return err
	} else if err := WriteFile(dst, content, fi.Mode()); err != nil {
		return err
	}
	return nil
}