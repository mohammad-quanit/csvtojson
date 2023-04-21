package main

import (
	"flag"
	"os"
	"reflect"
	"testing"
)

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
		{"No parameters", InputFile{}, true, []string{"cmd"}},
		{"Semicolon enabled", InputFile{"test.csv", "semicolon", false}, false, []string{"cmd", "--separator=semicolon", "test.csv"}},
		{"Pretty enabled", InputFile{"test.csv", "comma", true}, false, []string{"cmd", "--pretty", "test.csv"}},
		{"Pretty and semicolon enabled", InputFile{"test.csv", "semicolon", true}, false, []string{"cmd", "--pretty", "--separator=semicolon", "test.csv"}},
		{"Separator not identified", InputFile{}, true, []string{"cmd", "--separator=pipe", "test.csv"}},
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
