package search

import (
	"encoding/json"
	"fmt"
	"grafikart/grafisearch/utils"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Pinterest internal API response structure
type PinterestAPIResponse struct {
	ResourceResponse struct {
		Data []PinterestPin `json:"data"`
	} `json:"resource_response"`
}

type PinterestPin struct {
	Images map[string]struct {
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"images"`
	RichSummary struct {
		DisplayName string `json:"display_name"`
		URL         string `json:"url"`
	} `json:"rich_summary"`
	GridTitle     string `json:"grid_title"`
	Description   string `json:"description"`
	OriginalLink  string `json:"original_link"`
}

func GetPinterestImageResults(q string) ([]ImageResult, error) {
	start := time.Now()
	utils.LogInfo("PINTEREST_IMAGES_START", map[string]interface{}{
		"query": q,
	})

	// Use Pinterest's internal search API
	apiURL := fmt.Sprintf("https://www.pinterest.com/resource/BaseSearchResource/get/?source_url=/search/pins/&data={\"options\":{\"query\":\"%s\",\"scope\":\"pins\"},\"context\":{}}", url.QueryEscape(q))
	
	utils.LogInfo("PINTEREST_IMAGES_API_CALL", map[string]interface{}{
		"query": q,
		"url":   apiURL,
	})

	// Create request with proper headers to mimic browser
	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		duration := time.Since(start)
		utils.LogError("PINTEREST_IMAGES_REQUEST_ERROR", err, map[string]interface{}{
			"query":    q,
			"duration": duration,
		})
		return nil, err
	}

	// Add headers to mimic browser request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", "https://www.pinterest.com/")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := client.Do(req)
	if err != nil {
		duration := time.Since(start)
		utils.LogError("PINTEREST_IMAGES_FETCH_ERROR", err, map[string]interface{}{
			"query":    q,
			"duration": duration,
		})
		return nil, err
	}
	defer resp.Body.Close()

	utils.LogInfo("PINTEREST_IMAGES_HTTP_RESPONSE", map[string]interface{}{
		"query":       q,
		"status_code": resp.StatusCode,
		"status":      resp.Status,
	})

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		duration := time.Since(start)
		utils.LogError("PINTEREST_IMAGES_READ_ERROR", err, map[string]interface{}{
			"query":    q,
			"duration": duration,
		})
		return nil, err
	}

	// Parse JSON response
	var apiResponse PinterestAPIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		duration := time.Since(start)
		utils.LogError("PINTEREST_IMAGES_JSON_ERROR", err, map[string]interface{}{
			"query":    q,
			"duration": duration,
			"body":     string(body)[:200], // Log first 200 chars for debugging
		})
		return nil, err
	}

	results := []ImageResult{}
	
	utils.LogInfo("PINTEREST_IMAGES_PARSING_DATA", map[string]interface{}{
		"query": q,
		"pins":  len(apiResponse.ResourceResponse.Data),
	})

	// Process Pinterest pins
	for i, pin := range apiResponse.ResourceResponse.Data {
		if i >= 20 { // Limit results
			break
		}

		result := convertPinterestPin(pin, q)
		if result != nil {
			results = append(results, *result)
			
			utils.LogInfo("PINTEREST_IMAGES_RESULT_FOUND", map[string]interface{}{
				"query":     q,
				"index":     len(results),
				"title":     result.Title,
				"url":       result.URL,
				"source":    result.Source,
			})
		}
	}

	duration := time.Since(start)
	utils.LogInfo("PINTEREST_IMAGES_COMPLETED", map[string]interface{}{
		"query":    q,
		"results":  len(results),
		"duration": duration,
	})

	return results, nil
}

func convertPinterestPin(pin PinterestPin, query string) *ImageResult {
	var imageURL, thumbnailURL string
	var width, height int
	
	// Get the best image quality available
	if origImg, exists := pin.Images["orig"]; exists {
		imageURL = origImg.URL
		width = origImg.Width
		height = origImg.Height
	} else if img474, exists := pin.Images["474x"]; exists {
		imageURL = img474.URL
		width = img474.Width
		height = img474.Height
	}
	
	// Get thumbnail
	if thumb236, exists := pin.Images["236x"]; exists {
		thumbnailURL = thumb236.URL
	} else if thumb474, exists := pin.Images["474x"]; exists {
		thumbnailURL = thumb474.URL
	} else {
		thumbnailURL = imageURL // Fallback
	}
	
	// Get title from various sources
	title := pin.GridTitle
	if title == "" {
		title = pin.Description
	}
	if title == "" {
		title = pin.RichSummary.DisplayName
	}
	if title == "" {
		title = fmt.Sprintf("Pinterest image for: %s", query)
	}
	
	// Get source URL
	sourceURL := pin.OriginalLink
	if sourceURL == "" {
		sourceURL = pin.RichSummary.URL
	}
	if sourceURL == "" {
		sourceURL = "https://www.pinterest.com/"
	}
	
	// Parse domain from source URL
	domain := "pinterest.com"
	if sourceURL != "" && sourceURL != "https://www.pinterest.com/" {
		if u, err := url.Parse(sourceURL); err == nil {
			domain = u.Host
		}
	}
	
	if imageURL != "" {
		return &ImageResult{
			URL:       imageURL,
			Title:     title,
			Thumbnail: thumbnailURL,
			Source:    sourceURL,
			Domain:    domain,
			Width:     width,
			Height:    height,
		}
	}
	
	return nil
}