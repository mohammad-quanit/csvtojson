package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func checkIfValidFile(filename string) (bool, error) {
	// Checking if entered file is CSV by using the filepath package from the standard library
	if fileExt := filepath.Ext(filename); fileExt != ".csv" {
		return false, fmt.Errorf("file %s is not a csv file", filename)
	}

	// Checking if filepath entered belongs to an existing file. We use the Stat method from the os package (standard library)
	if _, err := os.Stat(filename); err != nil && os.IsNotExist(err) {
		return false, fmt.Errorf("file %s doesn't exist", filename)
	}

	// If we get to this point, it means this is a valid file
	return true, nil

}
