package main

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download/", downloadHandler)

    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    fmt.Println("Server listening on port 8080...")
    http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request)  {
	http.ServeFile(w,r, "templates/home.html")
}

func uploadHandler(w http.ResponseWriter, r *http.Request)  {
    err := r.ParseMultipartForm(10 << 30)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	file, handler, err := r.FormFile("file")
	if err != nil {
        http.Error(w, "Failed to retrieve file from form", http.StatusBadRequest)
        return
    }
	defer file.Close()

    f, err := os.Create("static/" + handler.Filename)
    if err != nil {
        http.Error(w, "Failed to create file on server", http.StatusInternalServerError)
        return
    }
    defer f.Close()

	 _, err = io.Copy(f, file)
	 if err != nil {
		 http.Error(w, "Failed to copy file data", http.StatusInternalServerError)
		 return
	 }
 
	 http.Redirect(w, r, "/", http.StatusSeeOther)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
    fileName := r.URL.Path[len("/download/"):]

    filePath := filepath.Join("static", fileName)
    _, err := os.Stat(filePath)
    if err != nil {
        http.Error(w, "File not found", http.StatusNotFound)
        return
    }

    file, err := os.Open(filePath)
    if err != nil {
        http.Error(w, "Failed to open file", http.StatusInternalServerError)
        return
    }
    defer file.Close()

    contentType := mime.TypeByExtension(filepath.Ext(filePath))
    if contentType == "" {
        contentType = "application/octet-stream" 
    }

    w.Header().Set("Content-Type", contentType)

    w.Header().Set("Content-Disposition", "attachment; filename="+fileName)

    _, err = io.Copy(w, file)
    if err != nil {
        http.Error(w, "Failed to send file", http.StatusInternalServerError)
        return
    }
}
