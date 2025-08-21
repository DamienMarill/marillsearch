package api

import (
	"encoding/json"
	"grafikart/grafisearch/utils"
	"net/http"
)

type wallpaperResponse struct {
	Wallpaper string `json:"wallpaper"`
}

func WallpaperHandler(w http.ResponseWriter, r *http.Request) {
	// Parse max_size parameter from query
	maxSize := utils.ParseMaxSize(r)

	wallpaper, err := utils.FetchAnimeWallpaper(maxSize)
	if err != nil {
		serveError(w, err)
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=300") // 5 minutes cache for anime variety
	j, _ := json.Marshal(wallpaperResponse{Wallpaper: wallpaper})
	w.Write(j)
}
