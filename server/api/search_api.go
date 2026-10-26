package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"grafikart/grafisearch/search"
	"grafikart/grafisearch/server"
	"grafikart/grafisearch/utils"
)

func SearchWithParser(fn func(string) ([]search.SearchResult, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")

		if q == "" {
			serveError(w, errors.New("missing query parameter"))
			return
		}
		results, err := fn(server.ReplaceFilterBangs(q))
		if err != nil {
			serveError(w, err)
			return
		}

		w.Header().Set("Cache-Control", "public, max-age=3600")
		w.Header().Set("Content-Type", "application/json")

		data, err := json.Marshal(results)
		if err != nil {
			serveError(w, err)
			return
		}
		w.Write(data)
	}
}

func ImageSearchWithParser(fn func(string) ([]search.ImageResult, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		q := r.URL.Query().Get("q")
		userAgent := r.Header.Get("User-Agent")
		clientIP := r.Header.Get("X-Forwarded-For")
		if clientIP == "" {
			clientIP = r.RemoteAddr
		}

		utils.LogInfo("IMAGE_SEARCH_START", map[string]interface{}{
			"query":      q,
			"client_ip":  clientIP,
			"user_agent": userAgent,
		})

		if q == "" {
			utils.LogError("IMAGE_SEARCH_ERROR", errors.New("missing query parameter"), map[string]interface{}{
				"query": q,
			})
			serveError(w, errors.New("missing query parameter"))
			return
		}

		processedQuery := server.ReplaceFilterBangs(q)
		if processedQuery != q {
			utils.LogInfo("IMAGE_SEARCH_QUERY_PROCESSED", map[string]interface{}{
				"original":  q,
				"processed": processedQuery,
			})
		}

		utils.LogInfo("IMAGE_SEARCH_CALLING_PROVIDER", map[string]interface{}{
			"query": processedQuery,
		})

		results, err := fn(processedQuery)
		duration := time.Since(start)

		if err != nil {
			utils.LogError("IMAGE_SEARCH_PROVIDER_ERROR", err, map[string]interface{}{
				"query":    q,
				"duration": duration,
			})
			serveError(w, err)
			return
		}

		utils.LogInfo("IMAGE_SEARCH_SUCCESS", map[string]interface{}{
			"query":    q,
			"results":  len(results),
			"duration": duration,
		})

		w.Header().Set("Cache-Control", "public, max-age=3600")
		w.Header().Set("Content-Type", "application/json")

		data, err := json.Marshal(results)
		if err != nil {
			utils.LogError("IMAGE_SEARCH_JSON_ERROR", err, map[string]interface{}{
				"query": q,
			})
			serveError(w, err)
			return
		}
		
		utils.LogInfo("IMAGE_SEARCH_RESPONSE", map[string]interface{}{
			"query":         q,
			"response_size": len(data),
		})
		w.Write(data)
	}
}
