package main

import (
    "fmt"
    "log"
    "strings"
    "net/http"
    "encoding/json"
    "github.com/undertideco/bandcamp"
)

var MyClient = &http.Client{}

type BandcampAlbum struct {
    find bool
    url string
}

/* 
check artist and album 
items[x].track.album.name et items[x].track.album.artists[0].name
*/
func testBandcamp(album string, artist string) BandcampAlbum{
    bandcampClient := bandcamp.NewClient()

    results, err := bandcampClient.Search(album)
    if err != nil {
        log.Println(err)
        return BandcampAlbum{false, ""}
    }

    if (strings.Contains(results[0].Title, album) || strings.Contains(album, results[0].Title)) && strings.Compare(results[0].Artist, artist) == 0 {
        return BandcampAlbum{true, results[0].URL}
    } else {
        return BandcampAlbum{false, ""}
    }
}

/*
id de la playlist
*/
func getListPlaylist(id string) {
    //req, _ := http.NewRequest("GET", "https://api.spotify.com/v1/playlists/"+id+"/tracks", nil)
    req, _ := http.NewRequest("GET", "http://localhost:8001", nil)
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

    tmp := BandcampAlbum{}

    for i := 0; i < len(ree.Items); i++ {
        tmp = testBandcamp(ree.Items[i].Track.Album.Name, ree.Items[i].Track.Album.Artists[0].Name)
        if tmp.find {
            fmt.Printf("Find !! url: %s\n", tmp.url)
        }
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
