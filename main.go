package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Uploading File \n")

	//1. parse input, type multipart/form-data
	r.ParseMultipartForm(10 << 20)

	//2. retrieve file from posted form-data
	file, handler, err := r.FormFile("myFile")
	defer file.Close()
	if err != nil {
		fmt.Println("Error Retrieveing file from form-data")
		fmt.Println(err)
		return
	}

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	//3 Write temporary file on our server
	tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	tempFile.Write(fileBytes)

	//4 return wheter or not this has been succesful
	fmt.Fprintf(w, "Succesfully Uploaded File \n")
}

func setupRoutest() {
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":8080", nil)
}

func main() {
	fmt.Println("Go File Upload Tutorial")
	setupRoutest()
}
