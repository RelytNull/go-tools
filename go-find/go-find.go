package main

import ( 
	"fmt"
	"flag"
	"os"
	"path/filepath"
)

func findFiles(dir, name string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if info.Name() == name {
			fmt.Println(path)
		}
		return nil
	})
}

func main() {
	// define command-line flags
	dirPtr := flag.String("dir", ".", "Directory to search in")
	namePtr := flag.String("name", ".", "Name of the file to search for")

	// Parse the flags
	flag.Parse()

	// Check if the name is provided
	if *namePtr == "" {
		fmt.Println("Please provide a file name using the -name flag.")
		return
	}


	// Call the findFiles function
	err := findFiles(*dirPtr, *namePtr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

