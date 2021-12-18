package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
)

func exitGracefully(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

func check(err error) {
	if err != nil {
		exitGracefully(err)
	}
}

func processLine(headers []string, line []string) (map[string]string, error) {
	// Validating if we're getting the same number of headers and columns. Otherwise, we return an error
	if len(headers) != len(line) {
		return nil, errors.New("Line doesn't match headers format. Skipping")
	}

	// Creating the map we're going to populate
	recordMap := make(map[string]string)

	// For each header we're going to set a new map key with the corresponding column value
	for i, name := range headers {
		recordMap[name] = line[i]
	}

	// Returning our generated map
	return recordMap, nil
}

func processCSVFile(fileData inputFile, writerChan chan<- map[string]string) {
	// opening our file for reading
	file, err := os.Open(fileData.Filepath)
	check(err) // Checking for errors, we shouldn't get any

	// Don't forget to close the file once everything is done
	defer file.Close()

	// Defining a "headers" and "line" slice
	var headers, line []string

	// initializing out csv reader
	reader := csv.NewReader(file)

	// The default character separator is comma (,)
	// so we need to change to semicolon if we get that option from the terminal
	if fileData.separator == "semicolon" {
		reader.Comma = ';'
	}
	// Reading the first line, where we will find our headers
	headers, err = reader.Read()
	check(err) // Again, error checking

	// Now we're going to iterate over each line from the CSV file
	for {
		// We read one row (line) from the CSV.
		// This line is a string slice, with each element representing a column
		line, err = reader.Read()

		// If we get to End of the File, we close the channel and break the for-loop
		if err == io.EOF {
			close(writerChan)
			break
		} else if err != nil {
			exitGracefully(err) // If this happens, we got an unexpected error
		}

		// Processiong a CSV line\
		record, err := processLine(headers, line)

		if err != nil { // If we get an error here, it means we got a wrong number of columns, so we skip this line
			fmt.Printf("Line: %sError: %s\n", line, err)
			continue
		}

		// Otherwise, we send the processed record to the writer channel
		writerChan <- record
	}
}
