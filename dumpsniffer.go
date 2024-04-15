package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"log"
	"bufio"
	"time"
	"runtime"
	"sync"
)

func main() {
   if len(os.Args) < 2 {
   		fmt.Println("Not a valid path.")
   		return
   	}

   startTime := time.Now()
   var memStats runtime.MemStats
   filesInspected := 0
   path := os.Args[1]

   if isFile(path) {
       if (!isPHPFile(path)) {
   		    fmt.Println("Not a php file.")
            return
   		}

        occurences := checkDumpDieOccurences(path)
        runtime.ReadMemStats(&memStats)
        displayTimeAndFiles(startTime, 1, occurences, memStats.Alloc)
        return
   }

    if isDir(path) {
        // Get all files in directory, and subdirectory and subsubdirectory .. you got it
        // We assign the response of the recursive function to two different error variables, one for the directory being iterated and one for the file.
        // If no error occurs, check for dump occurrences; otherwise, display the error.
        occurences := 0

        // WaitGroup and Mutex are used to synchronize the goroutines
        var wg sync.WaitGroup

        dirError := filepath.Walk(path, func(filePath string, fileInfo os.FileInfo, fileError error) error {
            // file error is not null, stop the script and display the error
            if fileError != nil {
                fmt.Printf("Error accessing %s: %s\n", filePath, fileError)
                return fileError
            }

            if isPHPFile(filePath) {
                filesInspected++
                wg.Add(1)
                go func(filePath string) {
                    defer wg.Done()
                    occurences += checkDumpDieOccurences(filePath)
                }(filePath)
            }

            return nil
        })

        if dirError != nil {
            fmt.Printf("Cannot read the directory %s: %s\n", path, dirError)
            return
        }

        wg.Wait()
        runtime.ReadMemStats(&memStats)
        displayTimeAndFiles(startTime, filesInspected, occurences, memStats.Alloc)
        return
    }
}

func isFile(path string) bool {
    fileInfo, err := os.Stat(path)
    if err != nil {
        return false
    }

    return fileInfo.Mode().IsRegular()
}

func isDir(path string) bool {
    fileInfo, err := os.Stat(path)
    if err != nil {
        return false
    }

    return fileInfo.Mode().IsDir()
}

func isPHPFile(path string) bool {
    ext := filepath.Ext(path)
    ext = strings.ToLower(ext)
    return ext == ".php"
}

func checkDumpDieOccurences(filePath string) int {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Cannot open the file %s : %v", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
    occurencesFound := 0

	lineNumber := 1
	for scanner.Scan() {
	    // Get the line then lowercase it to be case insensitive
		line := scanner.Text()
        lowerLine := strings.ToLower(line)

        // Check occurences
		if strings.Contains(lowerLine, "dump(") {
			fmt.Printf("%s: dump/var_dump found on line %d \n", filePath, lineNumber)
			occurencesFound++
		}

		if strings.Contains(lowerLine, "die(") || strings.Contains(lowerLine, "die;") {
			fmt.Printf("%s: die found on line %d \n", filePath, lineNumber)
			occurencesFound++
		}

		lineNumber++
	}

	return occurencesFound
}

func displayTimeAndFiles(startTime time.Time, filesInspected int, occurences int, memoryAllocated uint64) {
    memoryInKo := float64(memoryAllocated) / float64(1024)
    elapsedTime := time.Since(startTime)
    fmt.Printf("Elapsed time: %s\n", elapsedTime)
    fmt.Printf("Number of files inspected: %d\n", filesInspected)
    fmt.Printf("%d occurences found\n", occurences)
    fmt.Printf("%.2fko memory allocated\n", memoryInKo)
}