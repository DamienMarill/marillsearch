package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	picReAPI = `https://pic.re/image`
)

func FetchAnimeWallpaper(maxSize int) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	url := picReAPI
	if maxSize > 0 {
		url = fmt.Sprintf("%s?max_size=%d", picReAPI, maxSize)
	}

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36`)

	resp, err := client.Do(request)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch wallpaper from pic.re: %v", err)
	}
	defer resp.Body.Close()

	// pic.re redirects to the actual image URL
	return resp.Request.URL.String(), nil
}

// Legacy function for backward compatibility
func FetchBingWallpaper() (string, error) {
	return FetchAnimeWallpaper(1920)
}

// Parse max_size from request query parameter
func ParseMaxSize(r *http.Request) int {
	maxSizeStr := r.URL.Query().Get("max_size")
	if maxSizeStr == "" {
		return 1920 // Default to 1920 for HD
	}
	
	maxSize, err := strconv.Atoi(maxSizeStr)
	if err != nil || maxSize < 256 || maxSize > 4096 {
		return 1920 // Fallback to safe default
	}
	
	return maxSize
}
