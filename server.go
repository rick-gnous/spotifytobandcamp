package main

import (
//    "bytes"
    "fmt"
    "log"
    "time"
    "net/http"
    "encoding/json"
)

var MyClient = &http.Client{}


type YoutubeAPI struct {
    Href  string `json:"href"`
    Items []struct {
        AddedAt time.Time `json:"added_at"`
        AddedBy struct {
            ExternalUrls struct {
                Spotify string `json:"spotify"`

            } `json:"external_urls"`
            Href string `json:"href"`
            ID   string `json:"id"`
            Type string `json:"type"`
            URI  string `json:"uri"`

        } `json:"added_by"`
        IsLocal      bool        `json:"is_local"`
        PrimaryColor interface{} `json:"primary_color"`
        Track        struct {
            Album struct {
                AlbumType string `json:"album_type"`
                Artists   []struct {
                    ExternalUrls struct {
                        Spotify string `json:"spotify"`

                    } `json:"external_urls"`
                    Href string `json:"href"`
                    ID   string `json:"id"`
                    Name string `json:"name"`
                    Type string `json:"type"`
                    URI  string `json:"uri"`

                } `json:"artists"`
                AvailableMarkets []string `json:"available_markets"`
                ExternalUrls     struct {
                    Spotify string `json:"spotify"`

                } `json:"external_urls"`
                Href   string `json:"href"`
                ID     string `json:"id"`
                Images []struct {
                    Height int    `json:"height"`
                    URL    string `json:"url"`
                    Width  int    `json:"width"`

                } `json:"images"`
                Name                 string `json:"name"`
                ReleaseDate          string `json:"release_date"`
                ReleaseDatePrecision string `json:"release_date_precision"`
                TotalTracks          int    `json:"total_tracks"`
                Type                 string `json:"type"`
                URI                  string `json:"uri"`

            } `json:"album"`
            Artists []struct {
                ExternalUrls struct {
                    Spotify string `json:"spotify"`

                } `json:"external_urls"`
                Href string `json:"href"`
                ID   string `json:"id"`
                Name string `json:"name"`
                Type string `json:"type"`
                URI  string `json:"uri"`

            } `json:"artists"`
            AvailableMarkets []string `json:"available_markets"`
            DiscNumber       int      `json:"disc_number"`
            DurationMs       int      `json:"duration_ms"`
            Episode          bool     `json:"episode"`
            Explicit         bool     `json:"explicit"`
            ExternalIds      struct {
                Isrc string `json:"isrc"`

            } `json:"external_ids"`
            ExternalUrls struct {
                Spotify string `json:"spotify"`

            } `json:"external_urls"`
            Href        string `json:"href"`
            ID          string `json:"id"`
            IsLocal     bool   `json:"is_local"`
            Name        string `json:"name"`
            Popularity  int    `json:"popularity"`
            PreviewURL  string `json:"preview_url"`
            Track       bool   `json:"track"`
            TrackNumber int    `json:"track_number"`
            Type        string `json:"type"`
            URI         string `json:"uri"`

        } `json:"track"`
        VideoThumbnail struct {
            URL interface{} `json:"url"`

        } `json:"video_thumbnail"`

    } `json:"items"`
    Limit    int         `json:"limit"`
    Next     string      `json:"next"`
    Offset   int         `json:"offset"`
    Previous interface{} `json:"previous"`
    Total    int         `json:"total"`

}

type Test struct {
    Tee string
}

type ErrorSpotify struct {
    href string
}

type Error struct {
    status, message string
}

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
    fmt.Printf("hm", res.Status)
    /*
    buf := new(bytes.Buffer)
    buf.ReadFrom(res.Body)

    fmt.Printf("coucou %+v", buf.String())
    */

    //ree := &ErrorSpotify{}
    var ree YoutubeAPI
    //defer res.Body.Close()
    //json.Unmarshal(res.Body, ree)
    err = json.NewDecoder(res.Body).Decode(&ree)
    if err != nil {
        fmt.Printf("error:", err)
        return
    }
    //json.NewDecoder(res.Body).Decode(interface{}(ree))
    //decoder := json.NewDecoder(res.Body).Decode(interface{}(ree))

    /*
    for token, _ :=  decoder.Token(); token != nil; token, _ =  decoder.Token() {
        fmt.Printf("value:", token)
    }
    */
    //decoder.Decode(interface{}(ree))
    fmt.Printf(ree.Href)
    fmt.Printf("coucou\n")
    return

    token := map[string]string{}
    fmt.Printf("coucou\n")

    for _, v := range token {
        fmt.Printf("coucou")
        fmt.Printf("\n", v)
    }
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
