package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func Test_checkIfValidFile(t *testing.T) {
	// Creating a temporal and empty CSV file
	tempFile, err := ioutil.TempFile("", "file*.csv")
	if err != nil {
		panic(err) // This should never happen
	}

	// Once all the tests are done. We delete the temporal file
	defer os.Remove(tempFile.Name())

	tests := []struct {
		name     string
		fileName string
		want     bool
		wantErr  bool
	}{
		// TODO: Add test cases.
		{"File does exist", tempFile.Name(), true, false},
		{"File does not exist", "nowhere/file.csv", false, true},
		{"File is not CSV", "file.txt", false, true},
	}

	// Iterating over our test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkIfValidFile(tt.fileName)

			//checking error
			if (err != nil) != tt.wantErr {
				t.Errorf("checkIfValidFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Checking the returning value
			if got != tt.want {
				t.Errorf("checkIfValidFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
