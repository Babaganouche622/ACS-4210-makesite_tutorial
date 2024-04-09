package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
)

func main() {
	// Define command-line flags
	fileFlag := flag.String("file", "", "Name of the input text file")
	dirFlag := flag.String("dir", ".", "Directory containing input text files")
	flag.Parse()

	// Initialize counters for text, Markdown, and HTML files
	var txtCount, mdCount, htmlCount int
	var totalSize float64
	var totalTime float64

	// Always walk through a directory, you can search for the given file in the directory
	if *dirFlag != "" {
		txtFiles, err := findTextFiles(*fileFlag, *dirFlag)
		if err != nil {
			log.Fatal(err)
		}
		// Generate HTML pages for each .txt or .md file found
		for _, txtFile := range txtFiles {
			currentSize, currentTime := generateHTMLFromText(txtFile)
			totalSize += currentSize
			totalTime += currentTime
			htmlCount++
			if strings.HasSuffix(txtFile, ".txt") {
				txtCount++
			} else if strings.HasSuffix(txtFile, ".md") {
				mdCount++
			}
		}
		// }
	}

	pluralFile := "file"
	if htmlCount > 1 {
		pluralFile += "s"
	}
	fmt.Printf("\n\n\x1b[32mDone! We just built %d HTML pages (%.1fkB total) %s for you in %s.\x1b[0m\n", htmlCount, totalSize, pluralFile, fmt.Sprintf("%.2fs", totalTime))
	fmt.Printf("\x1b[35mYou had %d .txt and %d .md %s.\x1b[0m\n", txtCount, mdCount, pluralFile)
}

// Function to find all .txt files in a directory
func findTextFiles(file string, dir string) ([]string, error) {
	var txtFiles []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == file {
			txtFiles = append(txtFiles, path)
		} else if file == "" && !info.IsDir() && (strings.HasSuffix(info.Name(), ".txt") || strings.HasSuffix(info.Name(), ".md")) {
			txtFiles = append(txtFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return txtFiles, nil
}

// Function to generate HTML from a text file
func generateHTMLFromText(filename string) (float64, float64) {
	startTime := time.Now()

	// Read the contents of the text file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("\x1b[34mError opening file: %s\x1b[0m\n", err)
	}
	defer file.Close()

	// Get the file stats
	fileStat, err := os.Stat(filename)
	if err != nil {
		log.Printf("\x1b[34mFile %s does not exist\x1b[0m\n", filename)
	}

	content := make([]byte, fileStat.Size())
	_, err = file.Read(content)
	if err != nil {
		log.Fatalf("\x1b[34mError reading file: %s\x1b[0m\n", err)
	}

	// Parse the HTML template
	tmpl, err := template.ParseFiles("template.tmpl")
	if err != nil {
		log.Fatalf("\x1b[34mError parsing template file: %s\x1b[0m\n", err)
	}

	// Get the base name of the input file
	baseName := filepath.Base(filename)

	// Construct the HTML file path relative to the project directory
	htmlFilename := filepath.Join("output", strings.TrimSuffix(baseName, filepath.Ext(baseName))+".html")

	// Ensure the output directory exists
	err = os.MkdirAll(filepath.Dir(htmlFilename), 0755)
	if err != nil {
		log.Fatalf("\x1b[34mError, couldn't make directory\x1b[0m\n", err)
	}

	f, err := os.Create(htmlFilename)
	if err != nil {
		log.Fatalf("\x1b[34mError creating HTML file: %s\x1b[0m\n", err)
	}
	defer f.Close()

	// Execute the template with the contents of the text file
	if strings.HasSuffix(filename, ".md") {
		// If the input file is Markdown, parse it before executing the template
		htmlContent := markdown.ToHTML(content, nil, nil)
		err = tmpl.Execute(f, template.HTML(htmlContent))
	} else {
		// Otherwise, directly execute the template
		err = tmpl.Execute(f, string(content))
	}
	if err != nil {
		log.Fatalf("\x1b[31mTemplate coudln't execute: \x1b[0m\n", err)
	}

	endTime := time.Now()

	elapsedTime := float64(endTime.Sub(startTime).Microseconds()) / 1000
	elapsedTimeString := fmt.Sprintf("%.2fs", elapsedTime)

	// Convert the file size to kilobytes
	var fileSizeKB float64 = float64(fileStat.Size()) / 1024

	// Print the generated HTML filename
	log.Println("\x1b[32mSucess! Generated HTML file\x1b[0m")
	fmt.Printf("File size of %s was: \x1b[33m%.1fKB\x1b[0m\n", filename, fileSizeKB)
	fmt.Printf("Time taken to convert %s: \x1b[33m%v\x1b[0m\n", filename, elapsedTimeString)
	fmt.Printf("\x1b[34m%s\x1b[0m converted to \x1b[34m%s \x1b[0m\n\n", filename, htmlFilename)
	return fileSizeKB, elapsedTime
}
