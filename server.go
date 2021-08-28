package main

import (
//    "bytes"
    "fmt"
    "log"
//    "time"
    "net/http"
    "encoding/json"
)

var MyClient = &http.Client{}

/*
id de la playlist
*/
func getListPlaylist(id string) {
    req, _ := http.NewRequest("GET", "https://api.spotify.com/v1/playlists/"+id+"/tracks", nil)
    req.Header.Add("Accept", "application/json")
    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("Authorization", SpotifyAPI)

    res, err := MyClient.Do(req)
    if err != nil {
        fmt.Printf("error:", err)
        return
    }

    if res.StatusCode > 300 {
        fmt.Printf("Erreur !!\n")
        return
    }

    ree := &SpotifyPlaylist{}
    defer res.Body.Close()
    err = json.NewDecoder(res.Body).Decode(&ree)
    if err != nil {
        fmt.Printf("error:", err)
        return
    }

    fmt.Printf(ree.Href)
}

func formHandler(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        fmt.Fprintf(w, "ParseForm() err: %v", err)
        return
    }

    fmt.Fprintf(w, "POST request successful\n")
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
    getListPlaylist("6OGZZ8tI45MB1d3EUEqNKI")
    return

    fileServer := http.FileServer(http.Dir("./static"))
    http.Handle("/", fileServer)
    http.HandleFunc("/hello", index)
    http.HandleFunc("/back", formHandler)

    fmt.Printf("Starting the server…\n")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
