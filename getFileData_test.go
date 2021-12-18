package main

import (
	"flag"
	"os"
	"reflect"
	"testing"
)

func Test_getFileData(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name    string    // name of the test
		want    inputFile // What inputFile instance we want our function to return.
		wantErr bool      //Whether or not we want an error
		osArgs  []string  // The command arguments used for this test
	}{
		// Here we're declaring each unit test input and output data as defined before
		{"Default parameters", inputFile{"file.csv", "comma", false}, false, []string{"cmd", "file.csv"}},
		{"No parameters", inputFile{}, true, []string{"cmd"}},
		{"Semicolon enabled", inputFile{"file.csv", "semicolon", false}, false, []string{"cmd", "--separator=semicolon", "file.csv"}},
		{"Pretty enabled", inputFile{"file.csv", "comma", true}, false, []string{"cmd", "--pretty", "file.csv"}},
		{"Pretty and semicolon enabled", inputFile{"file.csv", "semicolon", true}, false, []string{"cmd", "--pretty", "--separator=semicolon", "file.csv"}},
		{"Separator not identified", inputFile{}, true, []string{"cmd", "--separator=pipe", "file.csv"}},
	}

	//iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Saving the original os.Args reference
			actualOsArgs := os.Args
			// This defer function will run after the test is done
			defer func() {
				os.Args = actualOsArgs                                           // Restoring the original os.Args reference
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError) // Reseting the Flag command line. So that we can parse flags again
			}()

			os.Args = tt.osArgs             // Setting the specific command args for this test
			got, err := GetFileData()       // Runing the function we want to test
			if (err != nil) != tt.wantErr { // Asserting whether or not we get the corret error value
				t.Errorf("getFileData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) { // Asserting whether or not we get the corret wanted value
				t.Errorf("getFileData() = %v, want %v", got, tt.want)
			}
		})
	}
}
