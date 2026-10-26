package search

type Link struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type SearchResult struct {
	URL      string `json:"url"`
	Title    string `json:"title"`
	Desc     string `json:"desc"`
	Domain   string `json:"domain"`
	SiteName string `json:"siteName,omitempty"`
	Author   string `json:"author,omitempty"`
	Related  []Link `json:"related,omitempty"`
}

type ImageResult struct {
	URL       string `json:"url"`        // Direct image URL
	Title     string `json:"title"`      // Image title/alt text
	Thumbnail string `json:"thumbnail"`  // Thumbnail URL
	Source    string `json:"source"`     // Source page URL
	Domain    string `json:"domain"`     // Source domain
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Size      string `json:"size,omitempty"`
}
