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
    "github.com/gofiber/template/html"
    "github.com/gofiber/fiber/v2/middleware/recover"
    "github.com/gofiber/fiber/v2/middleware/session"
)

var MyClient = &http.Client{}
var MyResp = &RespBandcamp{}
var SpotifyAPI = newTokenUser()
//var SpotifyAPI = &TokenUser{}
var Session = session.New()
var Errors string

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

func testSpotifyPlaylist(token, tokentype, id string) error {
    req, e := http.NewRequest("GET",
        "https://api.spotify.com/v1/playlists/"+id+"/tracks", nil)
    if e != nil {
        return e
    }
    req.Header.Add("Accept", "application/json")
    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("Authorization", tokentype + " " + token)
    res, err := MyClient.Do(req)

    if err != nil {
        return err
    }

    switch res.StatusCode {
        case 400: return errors.New("La requete s’est mal exécutée.")
        case 401: return errors.New("La Playlist semble être privée.")
        case 403: return errors.New("Accès refusé (token peut-être périmé).")
        case 404: return errors.New("Playlist inexistante.")
    }

    return nil
}

func getAllTracksPlaylist(token, tokentype, id string, offset int) (SpotifyPlaylist, error) {
    ret := SpotifyPlaylist{}
    req, _ := http.NewRequest("GET",
        "https://api.spotify.com/v1/playlists/"+id+"/tracks?offset="+strconv.FormatInt(int64(offset), 10),
        nil)

    req.Header.Add("Accept", "application/json")
    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("Authorization", tokentype + " " + token)
    res, _ := MyClient.Do(req)

    playlist := &SpotifyPlaylist{}
    defer res.Body.Close()
    err := json.NewDecoder(res.Body).Decode(&playlist)
    if err != nil {
        return ret, err
        fmt.Printf("error:", err)
    }

    ret = *playlist
    if ret.Total > offset {
        r, e := getAllTracksPlaylist(token, tokentype, id, offset + 100)
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
func getListPlaylist(id, token, tokentype string) {

    playlist, err := getAllTracksPlaylist(token, tokentype, id, 0)
    if err != nil {
        //fmt.Printf("error:", err.Error())
        log.Panic("error")
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
    sess, _ := Session.Get(c)

    if sess.Get("token") == nil {
        log.Panic("Vous n’êtes pas connecté à Spotify.")
    }

    token := sess.Get("token").(string)
    tokentype := sess.Get("tokentype").(string)
    id := c.FormValue("id")

    e := testSpotifyPlaylist(token, tokentype, id)
    if e != nil {
        log.Panic(e.Error())
    }


    c.Set("Location", "/feudecamp.html")
    go getListPlaylist(id, token, tokentype)
    return c.SendStatus(303)
}

func getNew(c *fiber.Ctx) error {
    /* read session */
    c.JSON(MyResp)
    MyResp.Albums = nil
    MyResp.Artists = nil
    MyResp.Notfound = nil
    return c.SendStatus(201)
}

func mytoken(c *fiber.Ctx) error {
    err := c.BodyParser(&SpotifyAPI)
    if err != nil {
        log.Panic(err.Error())
    } else {
        sess, err := Session.Get(c)
        if err != nil {
            log.Panic(err.Error())
        }

        sess.Set("token", SpotifyAPI.Token)
        sess.Set("expire", SpotifyAPI.ExpiresIn)
        sess.Set("tokentype", SpotifyAPI.TokenType)
        sess.Set("creation", time.Now().GoString())
        err = sess.Save()

        if err != nil {
            log.Panic(err.Error())
        }
    }

    //sess.Set("creation", time.Now())


    return c.SendStatus(201)
}

func spotifyCallback(c *fiber.Ctx) error {
    if c.Query("error") != "" {
        return c.Render("index", fiber.Map{"error": "Erreur lors de la connexion.",})
    } else {
        return c.Render("spotify-token", fiber.Map{})
    }
}

func index(c *fiber.Ctx) error {
    sess, _ := Session.Get(c)

    tmp := false
    if sess.Get("token") != nil {
        tmp = true
    }

    if Errors == "" {
        //return c.Render("index", fiber.Map{"connected": !SpotifyAPI.CheckEmpty(),
        return c.Render("index", fiber.Map{"connected": tmp,
        "url": SpotifyURL})
    } else {
        e := Errors
        Errors = ""
        return c.Render("index", fiber.Map{"connected": tmp,
        "error": e,
        "url": SpotifyURL})
    }
}

func main() {
    //app := fiber.New(fiber.Config(Views: html, ViewsLayout: "layouts/main"))
    app := fiber.New(fiber.Config{Views: html.New("./views", ".html"),})
    app.Use(recover.New())
    app.Static("/", "./static")

    app.Get("/", index)
    app.Post("/", mytoken)
    app.Post("/feudecamp", getNew)
    app.Post("/back", formHandler)
    app.Get("/callback", spotifyCallback)

    log.Fatal(app.Listen(":8080"))
}
