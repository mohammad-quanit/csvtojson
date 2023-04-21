package main

import (
	"flag"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func Test_checkIfValidFile(t *testing.T) {
	// creating a temporary csv file
	tempFile, err := ioutil.TempFile("", "test*.csv")
	if err != nil {
		panic(err)
	}

	// Once all the tests are done. We delete the temporal file
	defer os.Remove(tempFile.Name())

	tests := []struct {
		name     string
		filename string
		want     bool
		wantErr  bool
	}{
		{"File does exist", tempFile.Name(), true, false},
		{"File does not exist", "folder/" + tempFile.Name(), false, true},
		{"File is not csv", "test.txt", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkIfValidFile(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkIfValidFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkIfValidFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetFileData(t *testing.T) {
	// Defining our test slice. Table driven approach
	tests := []struct {
		name    string
		want    InputFile
		wantErr bool     // whether or not we want an error.
		osArgs  []string // The command arguments used for this test
	}{
		// Here we're declaring each unit test input and output data as defined before
		{"Default parameters", InputFile{"test.csv", "comma", false}, false, []string{"cmd", "test.csv"}},
		{"Semicolon enabled", InputFile{"test.csv", "semicolon", false}, false, []string{"cmd", "--separator=semicolon", "test.csv"}},
		{"Pretty enabled", InputFile{"test.csv", "comma", true}, false, []string{"cmd", "--pretty", "test.csv"}},
		{"Pretty and semicolon enabled", InputFile{"test.csv", "semicolon", true}, false, []string{"cmd", "--pretty", "--separator=semicolon", "test.csv"}},
		{"Separator not identified", InputFile{}, true, []string{"cmd", "--separator=pipe", "test.csv"}},
		{"No parameters", InputFile{}, true, []string{"cmd"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Saving the original os.Args reference
			actualOsArgs := os.Args

			// This defer function run after test is done
			defer func() {
				os.Args = actualOsArgs                                           // Restoring the actual OS args ref
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError) // Reseting the Flag command line. So that we can parse flags again
			}()

			os.Args = tt.osArgs       // Setting the specific command args for this test
			got, err := getFileData() // Runing the function we want to test

			// Asserting whether or not we get the corret error value
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Asserting whether or not we get the corret wanted value
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFileData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_processCsvFile(t *testing.T) {
	// Defining the maps we're expecting to get from our function
	wantMapSlice := []map[string]string{
		{"COL1": "1", "COL2": "2", "COL3": "3"},
		{"COL1": "4", "COL2": "5", "COL3": "6"},
	}

	// type args struct {
	// 	fileData      InputFile
	// 	writerChannel chan<- map[string]string
	// }
	tests := []struct {
		name      string // The name of the test
		csvString string // The content of our tested CSV file
		seperator string // The separator used for each test case
	}{
		{"Comma separator", "COL1,COL2,COL3\n1,2,3\n4,5,6\n", "comma"},
		{"Semicolon separator", "COL1;COL2;COL3\n1;2;3\n4;5;6\n", "semicolon"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Creating a CSV temp file for testing
			tmpFile, err := ioutil.TempFile("", "test*.csv")
			exitGracefully(err)

			defer os.Remove(tmpFile.Name()) // Removing the CSV test file before living

			// Defining the inputFile struct that we're going to use as one parameter of our function
			testFileData := InputFile{
				filePath:  tmpFile.Name(),
				pretty:    false,
				seperator: tt.seperator,
			}

			// Defining the writerChanel
			writerChannel := make(chan map[string]string)

			// calling the function as goroutine
			go processCsvFile(testFileData, writerChannel)

			// iterating over the slice containing the expected map values
			for _, wantMap := range wantMapSlice {
				record := <-writerChannel // Waiting for the record that we want to compare
				if !reflect.DeepEqual(record, wantMap) {
					t.Errorf("processCsvFile() = %v, want %v", record, wantMap)
				}
			}
		})
	}
}
