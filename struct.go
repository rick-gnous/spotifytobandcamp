package main

import "time"

type SpotifyPlaylist struct {
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
