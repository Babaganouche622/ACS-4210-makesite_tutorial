package main

import (
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
)

func main() {
	// Define command-line flags
	fileFlag := flag.String("file", "", "Name of the input text file")
	dirFlag := flag.String("dir", ".", "Directory containing input text files")
	flag.Parse()

	// If both flags are provided, prioritize `file` flag
	if *fileFlag != "" {
		generateHTMLFromText(*fileFlag)
	} else {
		// If only `dir` flag is provided, find all .txt files in the directory
		if *dirFlag != "" {
			txtFiles, err := findTextFiles(*dirFlag)
			if err != nil {
				log.Fatal(err)
			}
			// Generate HTML pages for each .txt file found
			for _, txtFile := range txtFiles {
				generateHTMLFromText(txtFile)
			}
		}
	}
}

// Function to find all .txt files in a directory
func findTextFiles(dir string) ([]string, error) {
	var txtFiles []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".txt") {
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
func generateHTMLFromText(filename string) {
	// Read the contents of the text file
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	// Parse the HTML template
	tmpl, err := template.ParseFiles("template.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	// Create a new HTML file based on the text file name
	htmlFilename := strings.TrimSuffix(filename, filepath.Ext(filename)) + ".html"
	f, err := os.Create(htmlFilename)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	// Print the generated HTML filename
	log.Println("Generated HTML file:", htmlFilename)
}
