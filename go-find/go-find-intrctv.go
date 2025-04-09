package main

import ( 
	"fmt"
	"flag"
	"os"
	"path/filepath"
	"sync"
	"github.com/sahilm/fuzzy"
	"github.com/chzyer/readline"
)


func findFiles(dir, name string, fuzzySearch bool, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if fuzzySearch {
			matches := fuzzy.Find(name, []string{info.Name()})
			if len(matches) > 0 {
				results <- path
			}
		} else {
			// Exact match
		if info.Name() == name {
			results <- path
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error %v\n:", err)
	}
}

func main() {
	// define command-line flags
	dirPtr := flag.String("d", ".", "Directory to search in")
	//namePtr := flag.String("n", "", "Name of the file to search for")
	fuzzyPtr := flag.Bool("f", false, "Enable fuzzy search")
	// Parse the flags
	flag.Parse()
	// Create a channel to receive reults
	results := make(chan string)
	var wg sync.WaitGroup
	// Start a goroutine to search for files
	go func() {
		for {
			// Use readline for better input handling
			rl, err := readline.NewEx(&readline.Config{
				Prompt: "Enter file name to search: ",
			//	History: readline.Newlist(),
			})
			if err != nil {
				fmt.Println("Error creating readline instance:", err)
				return
			}
			
			// Read user input
			name, err := rl.Readline()
			if err != nil {
				fmt.Println("Error reading line:", err)
				return
			}


			// Clear previous results
			fmt.Println("Results:")

			// Start searching
			wg.Add(1)
			go findFiles(*dirPtr, name, *fuzzyPtr, results, &wg)

			// Wait for the search to complete
			go func() {
				wg.Wait()
				close(results)
			}()
		}
	}()

	// Display results as they come in
	for path := range results {
		fmt.Println(path)
	}
}

