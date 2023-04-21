package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
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

func processCsvFile(fileData InputFile, writerChannel chan<- map[string]string) {
	// Opening file for reading
	file, err := os.Open(fileData.filePath)
	exitGracefully(err)

	// Closing the file when processing is done
	defer file.Close()

	//Defining a headers and line slice
	var headers, line []string

	// Initialiing our csv reader
	reader := csv.NewReader(file)

	// The default character separator is comma (,) so we need to change to semicolon if we get that option from the terminal
	if fileData.seperator == "semicolon" {
		reader.Comma = ';'
	}

	// Reading the first line, where we will find our headers
	headers, err = reader.Read()
	exitGracefully(err)

	// Now we're going to iterate over each line from the CSV file
	for {
		// We read one row (line) from the CSV.
		// This line is a string slice, with each element representing a column
		line, err = reader.Read()

		// If we get to End of the File, we close the channel and break the for-loop
		if err == io.EOF {
			close(writerChannel)
			break
		} else if err != nil {
			exitGracefully(err)
		}

		// Processiong a CSV line
		record, err := processLine(headers, line)

		if err != nil {
			fmt.Printf("Line: %s, Error: %s \n", line, err)
			continue
		}
		// Otherwise, we send the processed record to the writer channel
		writerChannel <- record
	}
}

func processLine(headers []string, dataList []string) (map[string]string, error) {
	// Validating if we're getting the same number of headers and columns. Otherwise, we return an error
	if len(dataList) != len(headers) {
		return nil, errors.New("line doesn't match headers format. Skipping")
	}
	// Creating the map we're going to populate
	recordMap := make(map[string]string)

	// For each header we're going to set a new map key with the corresponding column value
	for i, record := range headers {
		recordMap[record] = dataList[i]
	}

	// Returning our generated map
	return recordMap, nil
}

// func writeJsonFile(csvPath string, writerChannel <-chan map[string]string, done chan<- bool, pretty bool) {
// 	writeString := createStringWriter(csvPath) // Instantiating a JSON writer function
// 	jsonFunc, breakLine := getJSONFunc(pretty)                        //Instantiating the JSON parse function and the breakline character

// 	// log for informing
// 	fmt.Println("Writing JSON file...")

// 	// Writing the first character of our JSON file. We always start with a "[" since we always generate array of record
// 	writeString("["+breakLine, false)

// 	first := true

// 	for {
// 		// waiting from pushed records
// 	}
// }

func exitGracefully(err error) {
	// Exiting gracefullly on error
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
