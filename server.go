package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type contact struct {
	Email         string
	FirstName     string
	Last4SSN      string
	LastName      string
	Message       string
	PhoneNumber   string
	TimeContacted string
}

func uploadFiles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

	// Access the other form fields.
	var c contact = contact{
		Email:         r.FormValue("email"),
		FirstName:     r.FormValue("firstName"),
		Last4SSN:      r.FormValue("last4SSN"),
		LastName:      r.FormValue("lastName"),
		Message:       r.FormValue("message"),
		PhoneNumber:   r.FormValue("phoneNumber"),
		TimeContacted: time.Now().String(),
	}

	contactInfoBytes, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(contactInfoBytes))
	fmt.Fprintf(w, "Successfully Uploaded Contact Info: %s\n", string(contactInfoBytes))

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Println("Error Parsing Form")
		fmt.Println(err)
		return
	}

	// Access the map of uploaded files by key.
	files := r.MultipartForm.File["files"]

	// Iterate through each uploaded file.
	for _, fileHeader := range files {
		// Open the file.
		file, err := fileHeader.Open()
		if err != nil {
			fmt.Println("Error Retrieving the File")
			fmt.Println(err)
			return
		}
		defer file.Close()

		fmt.Printf("Uploaded File: %+v\n", fileHeader.Filename)
		fmt.Printf("File Size: %+v\n", fileHeader.Size)
		fmt.Printf("MIME Header: %+v\n", fileHeader.Header)

		// Create a temporary file within our temp-images directory that follows
		// a particular naming pattern.
		tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer tempFile.Close()

		// Copy the contents of the uploaded file to the temporary file.
		_, err = io.Copy(tempFile, file)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Fprintf(w, "Successfully Uploaded File: %s\n", fileHeader.Filename)
	}

	fmt.Fprintf(w, "All files successfully uploaded\n")
}

func main() {
	port := 8080
	url := fmt.Sprintf(":%d", port)
	fmt.Println("listening and serving on", url)

	http.Handle("/", http.FileServer(http.Dir("")))
	http.HandleFunc("/upload", uploadFiles)
	http.ListenAndServe(url, nil)
}
