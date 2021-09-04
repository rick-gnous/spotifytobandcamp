package main

import (
    "fmt"
    "log"
    "time"
    "errors"
    "strings"
    "strconv"
    "net/http"
    "encoding/json"
    "github.com/undertideco/bandcamp"
)

var MyClient = &http.Client{}
var MyResp = &RespBandcamp{}
var SpotifyAPI = &TokenUser{}

func loginSpotify(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        fmt.Fprintf(w, "ParseForm() err: %v", err)
        return
    }

    w.Header().Set("Location", "https://accounts.spotify.com/authorize?client_id="+ClientID+"&response_type=token&redirect_uri="+RedirectURI)
    w.WriteHeader(http.StatusSeeOther)
}

/* 
check artist and album 
items[x].track.album.name et items[x].track.album.artists[0].name
*/
func searchAlbumBandcamp(album string, artist string) BandcampAlbum{
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

func searchArtistBandcamp(artist string) BandcampAlbum {
    bandcampClient := bandcamp.NewClient()

    results, err := bandcampClient.Search(artist)
    if err != nil {
        log.Println(err)
        return BandcampAlbum{false, ""}
    }

    if strings.Compare(results[0].Artist, artist) == 0 {
        return BandcampAlbum{true, strings.Split(results[0].URL, "/album/")[0]}
    } else {
        return BandcampAlbum{false, ""}
    }
}

func getAllTracksPlaylist(id string, offset int) (SpotifyPlaylist, error) {
    ret := SpotifyPlaylist{}
    req, e := http.NewRequest("GET",
        "https://api.spotify.com/v1/playlists/"+id+"/tracks?offset="+strconv.FormatInt(int64(offset), 10),
        nil)
    if e != nil {
        fmt.Printf("%+v", e)
    }
    req.Header.Add("Accept", "application/json")
    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("Authorization", SpotifyAPI.TokenType + " " + SpotifyAPI.Token)
    res, err := MyClient.Do(req)

    if err != nil {
        return ret, err
    }

    if res.StatusCode > 300 {
        fmt.Printf("code %d\n", res.StatusCode)
        return ret, errors.New("Erreur token ou playlist inexistante.")
    }

    playlist := &SpotifyPlaylist{}
    defer res.Body.Close()
    err = json.NewDecoder(res.Body).Decode(&playlist)
    if err != nil {
        return ret, err
        fmt.Printf("error:", err)
    }

    ret = *playlist
    fmt.Printf("\n%d cc %d\n", ret.Total, offset)
    if ret.Total > offset {
        r, e := getAllTracksPlaylist(id, offset + 100)
        if e != nil {
            return ret, e
        }

        ret.Items = append(ret.Items, r.Items...)
    }

    return ret, nil
}

/*
id de la playlist
*/
func getListPlaylist(id string) {
    playlist, err := getAllTracksPlaylist(id, 0)
    if err != nil {
        fmt.Printf("Erreru!!\n")
        fmt.Printf("%+v", err)
        return
    }
    tmp := BandcampAlbum{}
    MyResp.Todo = len(playlist.Items)
    MyResp.Done = 0

    for i := 0; i < len(playlist.Items); i++ {
        tmp = searchAlbumBandcamp(playlist.Items[i].Track.Album.Name,
            playlist.Items[i].Track.Album.Artists[0].Name)
        if tmp.find {
            MyResp.AddAlbum(newUrlBandcamp(
                playlist.Items[i].Track.Album.Artists[0].Name,
                playlist.Items[i].Track.Album.Name,
                playlist.Items[i].Track.Album.ExternalUrls.Spotify,
                tmp.url))
        } else {
            tmp = searchArtistBandcamp(playlist.Items[i].Track.Album.Artists[0].Name)
            if tmp.find {
                MyResp.AddArtist(newUrlBandcamp(
                    playlist.Items[i].Track.Album.Artists[0].Name,
                    playlist.Items[i].Track.Album.Name,
                    playlist.Items[i].Track.Album.ExternalUrls.Spotify,
                    tmp.url))
            } else {
                MyResp.AddNotfound(newUrlWoBandcamp(
                    playlist.Items[i].Track.Album.Artists[0].Name,
                    playlist.Items[i].Track.Album.Name,
                    playlist.Items[i].Track.Album.ExternalUrls.Spotify))
            }
        }

        MyResp.Done++

        if i % 10 == 0 {
            time.Sleep(5 * time.Second)
        }
    }
    fmt.Printf("\nFinish\n")
}

func formHandler(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        fmt.Fprintf(w, "ParseForm() err: %v", err)
        return
    }

    w.Header().Set("Location", "/feudecamp.html")
    w.WriteHeader(http.StatusSeeOther)
    go getListPlaylist(r.FormValue("id"))
}

func hello(w http.ResponseWriter, r *http.Request) {
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

func getNew(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(MyResp)
    MyResp.Albums = nil
    MyResp.Artists = nil
    MyResp.Notfound = nil
}

func mytoken(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    err := json.NewDecoder(r.Body).Decode(&SpotifyAPI)
    if err != nil {
        fmt.Printf("error:", err)
    }
}

func main() {
    fileServer := http.FileServer(http.Dir("./static"))
    http.Handle("/", fileServer)
    http.HandleFunc("/hello", hello)
    http.HandleFunc("/back", formHandler)
    http.HandleFunc("/refresh", getNew)
    http.HandleFunc("/spotify", loginSpotify)
    http.HandleFunc("/mytoken", mytoken)

    fmt.Printf("Starting the server…\n")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
