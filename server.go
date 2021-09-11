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
var Session = session.New()
var Queue = make(map[string]*RespBandcamp)

func checkToken(sess *session.Session) {
    expire := sess.Get("expire")
    if expire != nil {
        if time.Now().Unix() > sess.Get("creation").(int64) + expire.(int64) {
            sess.Destroy()
        }
    }
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

func getAllTracksPlaylist(token, tokentype, id string, offset int) SpotifyPlaylist {
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
        log.Print(err.Error())
        return ret
    }

    ret = *playlist
    if ret.Total > offset {
        r := getAllTracksPlaylist(token, tokentype, id, offset + 100)
        ret.Items = append(ret.Items, r.Items...)
    }

    return ret
}

/*
id de la playlist
*/
func getListPlaylist(id, token, tokentype string) {
    playlist := getAllTracksPlaylist(token, tokentype, id, 0)

    var find bool
    var tmp, artist, album, urlSpotify string
    var MyResp = &RespBandcamp{}
    MyResp.Todo = len(playlist.Items)
    MyResp.Done = 0
    Queue[token] = MyResp

    for i := 0; i < len(playlist.Items); i++ {
        album = playlist.Items[i].Track.Album.Name
        artist = playlist.Items[i].Track.Album.Artists[0].Name
        urlSpotify = playlist.Items[i].Track.Album.ExternalUrls.Spotify

        if (MyResp.ContainsAlbum(album, artist)) {
            MyResp.Todo--
            continue
        }

        find, tmp = searchAlbumBandcamp(album, artist)
        if find {
            MyResp.AddAlbum(newUrlBandcamp(artist, album, urlSpotify, tmp))
        } else {
            find, tmp = searchArtistBandcamp(artist)
            if find {
                MyResp.AddArtist(newUrlBandcamp(artist, album, urlSpotify, tmp))
            } else {
                MyResp.AddNotfound(newUrlWoBandcamp(artist, album, urlSpotify))
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

    checkToken(sess)

    if sess.Get("token") == nil {
        panic("Vous n’êtes pas connecté à Spotify.")
    }

    token := sess.Get("token").(string)
    tokentype := sess.Get("tokentype").(string)
    id := c.FormValue("id")

    e := testSpotifyPlaylist(token, tokentype, id)
    if e != nil {
        panic(e.Error())
    }

    Queue[token] = nil
    c.Set("Location", "/feudecamp")
    go getListPlaylist(id, token, tokentype)
    return c.SendStatus(303)
}

func getNew(c *fiber.Ctx) error {
    sess, _ := Session.Get(c)

    c.JSON(Queue[sess.Get("token").(string)])
    return c.SendStatus(201)
}

func mytoken(c *fiber.Ctx) error {
    tmp := newTokenUser()
    err := c.BodyParser(&tmp)
    if err != nil {
        log.Panic(err.Error())
    } else {
        sess, err := Session.Get(c)
        if err != nil {
            log.Panic(err.Error())
        }

        expire, _ := strconv.Atoi(tmp.ExpiresIn)

        sess.Set("token", tmp.Token)
        sess.Set("expire", int64(expire))
        sess.Set("tokentype", tmp.TokenType)
        sess.Set("creation", time.Now().Unix())
        err = sess.Save()

        if err != nil {
            log.Panic(err.Error())
        }
    }

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

    checkToken(sess)

    return c.Render("index", fiber.Map{"connected": sess.Get("token") != nil,
        "url": SpotifyURL})
}

func fdc(c *fiber.Ctx) error {
    sess, _ := Session.Get(c)

    if sess.Get("token") == nil {
        panic("Vous n’êtes pas connecté.")
    } else if _, err := Queue[sess.Get("token").(string)]; !err {
        panic("Vous n’avez pas lancé de playlist.")
    }
    return c.Render("feudecamp", fiber.Map{})
}

func main() {
    //app := fiber.New(fiber.Config(Views: html, ViewsLayout: "layouts/main"))
    app := fiber.New(fiber.Config{Views: html.New("./views", ".html"),})
    app.Use(recover.New())
    app.Static("/", "./static")

    app.Get("/", index)
    app.Post("/", mytoken)
    app.Get("/feudecamp", fdc)
    app.Post("/feudecamp", getNew)
    app.Post("/back", formHandler)
    app.Get("/callback", spotifyCallback)

    log.Fatal(app.Listen(":8080"))
}
