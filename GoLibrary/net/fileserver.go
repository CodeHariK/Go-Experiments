package net

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const uploadDirectory = "./uploads/"

func FileServer() {
	// Ensure the uploads directory exists
	if err := os.MkdirAll(uploadDirectory, 0o755); err != nil {
		fmt.Println("Error creating upload directory:", err)
		return
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download", downloadHandler)

	port := 8080
	fmt.Printf("Server listening on :%d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Display a simple HTML form
	tmpl, err := template.New("index").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>File Uploader/Downloader</title>
		</head>
		<body>
			<h1>File Uploader/Downloader</h1>
			<form action="/upload" method="post" enctype="multipart/form-data">
				<label for="file">Choose a file:</label>
				<input type="file" name="file" id="file" required>
				<button type="submit">Upload</button>
			</form>
			<hr>
			<h2>Download Files</h2>
			<ul>
				{{range .}}
					<li><a href="/download?file={{.}}">{{.}}</a></li>
				{{end}}
			</ul>
		</body>
		</html>
	`)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// List files in the uploads directory
	files, err := listFiles(uploadDirectory)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render the HTML template with the list of files
	if err := tmpl.Execute(w, files); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data, including the uploaded file
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a unique filename based on the current timestamp
	filename := time.Now().Format("20060102150405") + "_" + handler.Filename
	filePath := filepath.Join(uploadDirectory, filename)

	// Create the file on the server
	dst, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the server
	if _, err := io.Copy(dst, file); err != nil {
		fmt.Println("Error copying file:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Redirect to the index page after successful upload
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// Get the filename from the query parameter
	filename := r.URL.Query().Get("file")
	if filename == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Open the file on the server
	filePath := filepath.Join(uploadDirectory, filename)
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Set the appropriate content type and headers for file download
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", getFileSize(filePath)))

	// Copy the file to the response writer
	if _, err := io.Copy(w, file); err != nil {
		fmt.Println("Error copying file to response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func listFiles(directory string) ([]string, error) {
	var files []string

	// Walk through the files in the directory
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Exclude directories from the list
		if !info.IsDir() {
			files = append(files, info.Name())
		}
		return nil
	})

	return files, err
}

func getFileSize(filePath string) int64 {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0
	}
	return fileInfo.Size()
}
