package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"
)
func main() {
	    // Read the contents of first-post.txt
			content, err := ioutil.ReadFile("first-post.txt")
			if err != nil {
					log.Fatal(err)
			}
	
			// Parse the HTML template
			tmpl, err := template.ParseFiles("template.tmpl")
			if err != nil {
					log.Fatal(err)
			}
	
			// Create a new file first-post.html
			f, err := os.Create("first-post.html")
			if err != nil {
					log.Fatal(err)
			}
			defer f.Close()
	
			// Execute the template with the contents of first-post.txt
			err = tmpl.Execute(f, string(content))
			if err != nil {
					log.Fatal(err)
			}

			// Execute the template with the contents of first-post.txt and print it to stdout
			err = tmpl.Execute(os.Stdout, string(content))
			if err != nil {
					log.Fatal(err)
			}
}
