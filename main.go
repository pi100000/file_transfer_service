package main

import (
	"fmt"
	"net/http"
)

func main() {
    http.HandleFunc("/", homeHandler)
    fmt.Println("Server listening on port 8080...")
    http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request)  {
	http.ServeFile(w,r, "templates/home.html")
}



