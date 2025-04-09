package main

import ( 
	"fmt"
	"flag"
	"os"
	"path/filepath"
	"github.com/sahilm/fuzzy"
)


func findFiles(dir, name string, fuzzySearch bool) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if fuzzySearch {
			matches := fuzzy.Find(name, []string{info.Name()})
			if len(matches) > 0 {
				fmt.Println(path)
			}
		} else {
			// Exact match
		if info.Name() == name {
			fmt.Println(path)
		}
	}
		return nil
	})
}

func main() {
	// define command-line flags
	dirPtr := flag.String("d", ".", "Directory to search in")
	namePtr := flag.String("n", "", "Name of the file to search for")
	fuzzyPtr := flag.Bool("f", false, "Enable fuzzy search")

	// Parse the flags
	flag.Parse()

	// Check if the name is provided
	if *namePtr == "" {
		fmt.Println("Please provide a file name using the -name flag.")
		return
	}


	// Call the findFiles function
	err := findFiles(*dirPtr, *namePtr, *fuzzyPtr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

