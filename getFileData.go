package main

import (
	"errors"
	"flag"
	"os"
)

func GetFileData() (inputFile, error) {
	// We need to validate that we're getting the correct number of arguments
	if len(os.Args) < 2 {
		return inputFile{}, errors.New("a filepath argument is required")
	}

	// Defining option flags. For this, we're using the Flag package from the standard library
	// We need to define three arguments: the flag's name, the default value, and a short description
	// (displayed whith the option --help)
	separator := flag.String("separator", "comma", "Column Separator")
	pretty := flag.Bool("pretty", false, "Generate Pretty JSON")

	//parse all the flags
	flag.Parse()

	// The only argument (that is not a flag option) is the file location (CSV file)
	fileLocation := flag.Arg(0)

	// Validating whether or not we received "comma" or "semicolon" from the parsed arguments.
	// If we dind't receive any of those. We should return an error
	if !(*separator == "comma" || *separator == "semicolon") {
		return inputFile{}, errors.New("only comma and semicolons are allowed")
	}

	// If we get to this endpoint, our programm arguments are validated
	// We return the corresponding struct instance with all the required data
	return inputFile{fileLocation, *separator, *pretty}, nil
}
