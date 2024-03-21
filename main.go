package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/upload", uploadHandler)

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

