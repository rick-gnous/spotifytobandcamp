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
    "github.com/gofiber/fiber/v2"
    "github.com/undertideco/bandcamp"
)

var MyClient = &http.Client{}
var MyResp = &RespBandcamp{}
var SpotifyAPI = &TokenUser{}

func loginSpotify(c *fiber.Ctx) error {
    c.Set("Location", "https://accounts.spotify.com/authorize?client_id="+ClientID+"&response_type=token&redirect_uri="+RedirectURI)
    return c.SendStatus(303)
}

/* 
check artist and album 
items[x].track.album.name et items[x].track.album.artists[0].name
*/
func searchAlbumBandcamp(album string, artist string) (bool, string) {
    bandcampClient := bandcamp.NewClient()

    results, err := bandcampClient.Search(album)
    if err != nil {
        log.Println(err)
        return false, ""
    }

    if (strings.Contains(results[0].Title, album) || strings.Contains(album, results[0].Title)) && strings.Compare(results[0].Artist, artist) == 0 {
        return true, results[0].URL
    } else {
        return false, ""
    }
}

func searchArtistBandcamp(artist string) (bool, string) {
    bandcampClient := bandcamp.NewClient()

    results, err := bandcampClient.Search(artist)
    if err != nil {
        log.Println(err)
        return false, ""
    }

    if strings.Compare(results[0].Artist, artist) == 0 {
        return true, strings.Split(results[0].URL, "/album/")[0]
    } else {
        return false, ""
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
    var find bool
    var tmp string
    MyResp.Todo = len(playlist.Items)
    MyResp.Done = 0

    for i := 0; i < len(playlist.Items); i++ {
        find, tmp = searchAlbumBandcamp(playlist.Items[i].Track.Album.Name,
            playlist.Items[i].Track.Album.Artists[0].Name)
        if find {
            MyResp.AddAlbum(newUrlBandcamp(
                playlist.Items[i].Track.Album.Artists[0].Name,
                playlist.Items[i].Track.Album.Name,
                playlist.Items[i].Track.Album.ExternalUrls.Spotify,
                tmp))
        } else {
            find, tmp = searchArtistBandcamp(playlist.Items[i].Track.Album.Artists[0].Name)
            if find {
                MyResp.AddArtist(newUrlBandcamp(
                    playlist.Items[i].Track.Album.Artists[0].Name,
                    playlist.Items[i].Track.Album.Name,
                    playlist.Items[i].Track.Album.ExternalUrls.Spotify,
                    tmp))
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

func formHandler (c *fiber.Ctx) error {
    c.Set("Location", "/feudecamp.html")
    go getListPlaylist(c.FormValue("id"))
    return c.SendStatus(303)
}

func hello(c *fiber.Ctx) error {
    return c.SendString("Helo!")
}

func getNew(c *fiber.Ctx) error {
    c.JSON(MyResp)
    MyResp.Albums = nil
    MyResp.Artists = nil
    MyResp.Notfound = nil
    return c.SendStatus(201)
}

func mytoken(c *fiber.Ctx) error {
    return c.BodyParser(&SpotifyAPI)
}

func main() {
    app := fiber.New()

    app.Static("/", "./static")

    app.Get("/hello", hello)
    app.Get("/refresh", getNew)
    app.Get("/spotify", loginSpotify)

    app.Post("/back", formHandler)
    app.Post("/mytoken", mytoken)
    fmt.Printf("Starting the server…\n")
    app.Listen(":8080")
}
