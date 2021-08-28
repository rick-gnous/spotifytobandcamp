package main

import (
    "fmt"
    "log"
    "net/http"
)

func formHandler(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        fmt.Fprintf(w, "ParseForm() err: %v", err)
        return
    }

    fmt.Fprintf(w, "POST request successful")
    type_id := r.FormValue("type-id")
    id := r.FormValue("id")

    fmt.Fprintf(w, "Type ID = %s\n", type_id)
    fmt.Fprintf(w, "ID = %s\n", id)
}

func index(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/hello" {
        http.Error(w, "404", http.StatusNotFound)
        return
    }

    if r.Method != "GET" {
        http.Error(w, "Use GET pls", http.StatusNotFound)
        return
    }

    fmt.Fprintf(w, "Helo!")
}

func main() {
    fileServer := http.FileServer(http.Dir("./static"))
    http.Handle("/", fileServer)
    http.HandleFunc("/hello", index)
    http.HandleFunc("/back", formHandler)

    fmt.Printf("Starting the server…\n")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
