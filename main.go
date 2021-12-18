package main

import (
	"fmt"
)

type inputFile struct {
	Filepath  string
	separator string
	Pretty    bool
}

func main() {
	fmt.Println("CSV to JSON Cli Tool")
	fileData, err := GetFileData()
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	fmt.Println(fileData)
}
