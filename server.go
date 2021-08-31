package main

import (
    "fmt"
    "log"
    "time"
    "strings"
    "net/http"
    "encoding/json"
    "github.com/undertideco/bandcamp"
)

type BandcampAlbum struct {
    find bool
    url string
}

type RespBandcamp struct {
    Done int `json:"done"`
    Todo int `json:"todo"`
    Url []string `json:"url"`
}

func (rp *RespBandcamp) Add(str string) []string {
    rp.Url = append(rp.Url, str)
    return rp.Url
}

var MyClient = &http.Client{}
var MyResp = &RespBandcamp{}

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
    /*
    req, _ := http.NewRequest("GET",
                              "https://api.spotify.com/v1/playlists/"+id+"/tracks",
                              nil)
                              */
    req, _ := http.NewRequest("GET",
                              "http://localhost:8001",
                              nil)
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
    MyResp.Todo = len(ree.Items)
    MyResp.Done = 0

    for i := 0; i < len(ree.Items); i++ {
        tmp = testBandcamp(ree.Items[i].Track.Album.Name,
                           ree.Items[i].Track.Album.Artists[0].Name)
        if tmp.find {
            //fmt.Printf("Find !! url: %s\n", tmp.url)
            //MyResp.url = append(MyResp.url, tmp.url)
            //MyResp.url = append (MyResp.url, tmp.url)
            MyResp.Add(tmp.url)
            //fmt.Printf("tmp %s \n", MyResp.url[0])
            //fmt.Printf("len=%d cap=%d %v\n", len(MyResp.url), cap(MyResp.url), MyResp.url)
        }
        if i % 25 == 0 {
            time.Sleep(5 * time.Second)
        }
        MyResp.Done++
    }
    fmt.Printf("\nFinish\n")
}

func formHandler(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        fmt.Fprintf(w, "ParseForm() err: %v", err)
        return
    }

    //fmt.Fprintf(w, "POST request successful\n")
    //type_id := r.FormValue("type-id")
    //id := r.FormValue("id")

    //fmt.Fprintf(w, "Type ID = %s\n", type_id)
    //fmt.Fprintf(w, "ID = %s\n", id)

    w.Header().Set("Location", "/feudecamp.html")
    w.WriteHeader(http.StatusSeeOther)
    //go getListPlaylist("6OGZZ8tI45MB1d3EUEqNKI")
    go getListPlaylist(r.FormValue("id"))
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

func test(w http.ResponseWriter, r *http.Request) {
    getListPlaylist("6OGZZ8tI45MB1d3EUEqNKI")
}

func getNew(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(MyResp)
    MyResp.Url = nil
}

func main() {
    fileServer := http.FileServer(http.Dir("./static"))
    http.Handle("/", fileServer)
    http.HandleFunc("/hello", index)
    http.HandleFunc("/back", formHandler)
    http.HandleFunc("/test", test)
    http.HandleFunc("/refresh", getNew)

    fmt.Printf("Starting the server…\n")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
