package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"log"
	"bufio"
)

func main() {
   if len(os.Args) < 2 {
   		fmt.Println("Not a valid path.")
   		return
   	}

   	path := os.Args[1]

   	if isFile(path) {
   		if (!isPHPFile(path)) {
   		    fmt.Println("Not a php file.")
            return
   		}

   		checkDumpDieOccurences(path)
   		return
   	}

    if isDir(path) {
        // Get all files in directory, and subdirectory and subsubdirectory .. you got it
        // We assign the response of the recursive function to two different error variables, one for the directory being iterated and one for the file.
        // If no error occurs, check for dump occurrences; otherwise, display the error.
        dirError := filepath.Walk(path, func(filePath string, fileInfo os.FileInfo, fileError error) error {
            // file error is not null, stop the script and display the error
            if fileError != nil {
                fmt.Printf("Error accessing %s: %s\n", filePath, fileError)
                return fileError
            }

            if isPHPFile(filePath) {
                checkDumpDieOccurences(filePath)
            }

            // Not a php file, continue
            return nil
        })

        if dirError != nil {
            fmt.Printf("Cannot read the directory %s: %s\n", path, dirError)
            return
        }

        return
    }
}

func isFile(path string) (bool) {
    fileInfo, err := os.Stat(path)
    if err != nil {
        return false
    }
    if fileInfo.Mode().IsRegular() {
        return true
    }
    return false
}

func isDir(path string) (bool) {
    fileInfo, err := os.Stat(path)
    if err != nil {
        return false
    }
    if fileInfo.Mode().IsDir() {
        return true
    }
    return false
}

func isPHPFile(path string) bool {
    ext := filepath.Ext(path)
    ext = strings.ToLower(ext)
    return ext == ".php"
}

func checkDumpDieOccurences(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Cannot open the file %s : %v", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	lineNumber := 1
	for scanner.Scan() {
	    // Get the line then lowercase it to be case insensitive
		line := scanner.Text()
        lowerLine := strings.ToLower(line)

        // Check occurences
		if strings.Contains(lowerLine, "var_dump(") || strings.Contains(line, "dump(") {
			fmt.Printf("%s: dump/var_dump found on line %d \n", filePath, lineNumber)
		}

		if strings.Contains(lowerLine, "die(") || strings.Contains(line, "die;") {
			fmt.Printf("%s: die found on line %d \n", filePath, lineNumber)
		}

		lineNumber++
	}
}