package main

// based idea from https://dasarpemrogramangolang.novalagung.com/B-form-upload-file.html thx to @novalagung for tutorial on his repo

// got referenced from
// 1. https://dasarpemrogramangolang.novalagung.com/B-form-upload-file.html
// 2. https://stackoverflow.com/questions/768431/how-do-i-make-a-redirect-in-php
// 3. https://stackoverflow.com/a/59817862

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	// this need to serve that uploads since every file now put on uploads
	// without this file cannot be render by golang
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	// for rendering the html
	http.HandleFunc("/", routerIndexGet)
	http.HandleFunc("/upload", routerUploadPost)

	fmt.Println("Server Upload Test running...")
	fmt.Println("Visit http://localhost:9999 to test the code")
	http.ListenAndServe(":9999", nil)
}

func routerIndexGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Serve the HTML template
	tmpl, err := template.ParseFiles("view.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func routerUploadPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// for handling file
	r.ParseMultipartForm(8 << 20)            // parse the uploaded file size i limited on 8MB to save some memory
	file, handler, err := r.FormFile("file") // here came from the name of the form "name="file""
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Check if the file is a PNG
	if !strings.HasSuffix(handler.Filename, ".png") {
		http.Error(w, "Only PNG files are allowed", http.StatusBadRequest)
		return
	}

	// Construct the file path in the upload directory and rename it to "upload.png"
	filePath := filepath.Join("uploads", "upload.png")

	// Create the new file
	newFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer newFile.Close()

	// Copy the uploaded file to the new file
	_, err = io.Copy(newFile, file)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Redirect to the index page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
