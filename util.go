package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func getFileData() (InputFile, error) {
	// validating correct num of args
	if len(os.Args) < 2 {
		return InputFile{}, errors.New("a filepath argument is required")
	}

	// Defining option flag. For this, we are using Flag package from GO standard library
	// We are defining 3 arguments: flag's name, flag's default value and a short description (displayed whith the option --help)
	separator := flag.String("separator", "comma", "Column separator")
	pretty := flag.Bool("pretty", false, "Generate Pretty JSON")

	// Parse all arguments from terminal
	flag.Parse()

	// fileLocation argument (that is not a flag)
	fileLocation := flag.Arg(0)

	// Validating whether or not we received "comma" or "semicolon" from the parsed arguments.
	// If we dind't receive any of those. We should return an error
	if !(*separator == "comma" || *separator == "semicolon") {
		return InputFile{}, errors.New("only comma or semicolon separators are allowed")
	}

	// If we get to this endpoint, our programm arguments are validated
	// We return the corresponding struct instance with all the required data
	return InputFile{fileLocation, *separator, *pretty}, nil
}

func checkIfValidFile(filename string) (bool, error) {
	// Checking if entered file is CSV by using the filepath package from the standard library
	if fileExt := filepath.Ext(filename); fileExt != ".csv" {
		return false, fmt.Errorf("file %s is not csv", filename)
	}

	// Checking if filepath entered belongs to an existing file. We use the Stat method from the os package (standard library)
	if _, err := os.Stat(filename); err != nil && os.IsNotExist(err) {
		return false, fmt.Errorf("file %s does not exist", filename)
	}

	// If we get to this point, it means this is a valid file
	return true, nil
}
