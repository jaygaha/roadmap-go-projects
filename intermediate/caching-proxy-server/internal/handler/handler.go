package handler

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/jaygaha/roadmap-go-projects/intermediate/caching-proxy-server/internal/cache"
	"github.com/jaygaha/roadmap-go-projects/intermediate/caching-proxy-server/internal/config"
)

// copyHeaders copies all header key-values from src to dst
func copyHeaders(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

// fetchAndServeFromOrigin fetches the request from the origin server and serve it to the client
func fetchAndServeFromOrigin(w http.ResponseWriter, r *http.Request, shouldCache bool) {
	// construct the full target url for the origin server using config.OriginServerURL
	targetURL := config.OriginServerURL.ResolveReference(r.URL)

	// create a new HTTP request to the origin server
	originReq, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
	if err != nil {
		log.Printf("Error creating origin request for %s: %v", targetURL.String(), err)
		http.Error(w, "Error creating request to origin server", http.StatusInternalServerError)
		return
	}

	// copy headers from the incoming client request to the origin request
	copyHeaders(originReq.Header, r.Header)
	// Remove hop-by-hop headers that should not be forwarded to the origin server
	originReq.Header.Del("Connection")
	originReq.Header.Del("Proxy-Connection")
	originReq.Header.Del("Keep-Alive")
	originReq.Header.Del("Transfer-Encoding") // Let the http client handle chunked encoding

	// execut the request to the origin server
	httpClient := &http.Client{
		Timeout: 30 * time.Second, // set a timeout for the request
	}

	originResp, err := httpClient.Do(originReq)
	if err != nil {
		log.Printf("Error fetching from origin server for %s: %v", targetURL.String(), err)
		http.Error(w, "Error fetching from origin server", http.StatusBadGateway)
		return
	}
	defer originResp.Body.Close()

	// read the response body from the origin server
	originRespBody, err := io.ReadAll(originResp.Body)
	if err != nil {
		log.Printf("Error reading response body from origin server for %s: %v", targetURL.String(), err)
		http.Error(w, "Error reading response from origin server", http.StatusInternalServerError)
		return
	}

	// cache the response if needed
	isCacheableResp := r.Method == http.MethodGet &&
		shouldCache &&
		originResp.StatusCode >= http.StatusOK &&
		originResp.StatusCode < http.StatusMultipleChoices // 300 status code

	if isCacheableResp {
		// create a new CacheEntry and add it to the cache
		cacheEntry := cache.CacheEntry{
			Headers:      originResp.Header,
			ResponseData: originRespBody,
			StatusCode:   originResp.StatusCode,
			CreatedAt:    time.Now(),
		}
		cache.Set(r.URL.String(), cacheEntry)
		w.Header().Set("X-Cache", "MISS") // It was a miss, now cached
		log.Println("Cache Action: MISS")
	} else if r.Method == http.MethodGet {
		w.Header().Set("X-Cache", "MISS") // A miss, but not cached like error from origin
		log.Println("Cache Action: MISS")
	} else {
		w.Header().Set("X-Cache", "SKIP") // Not a cacheable response, no caching
		log.Println("Cache Action: SKIP")
	}

	// finally serve the response to the client
	copyHeaders(w.Header(), originResp.Header)
	w.WriteHeader(originResp.StatusCode)
	w.Write(originRespBody)
}

// HandleClearCache handles the cache clear request
func HandleClearCache(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	cache.Clear()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Cache cleared."))
}

// HandleProxyRequest is the core HTTP handler for proxuing and caching requests
func HandleProxyRequest(w http.ResponseWriter, r *http.Request) {
	reqKey := r.URL.String() // use the full URL path and query as the cache key

	// Step 1: Handle non-cacheable requests
	// Forward non-cacheable requests to the origin server and return
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		fetchAndServeFromOrigin(w, r, false)
		return
	}

	// Step 2: Check if the request is cacheable
	cachedEntry, found := cache.Get(reqKey)
	if found {
		log.Printf("Cache Action: HIT for %s", reqKey)
		w.Header().Set("X-Cache", "HIT")
		copyHeaders(w.Header(), cachedEntry.Headers)

		// ensure Content-Type is correctly set from cache
		w.Header().Set("Content-Type", cachedEntry.Headers.Get("Content-Type"))
		w.WriteHeader(cachedEntry.StatusCode)
		w.Write(cachedEntry.ResponseData)
		return
	}

	// Step 3: cache MISS - fetch from origin and potentially cache
	log.Printf("Cache Action: MISS for %s", reqKey)
	fetchAndServeFromOrigin(w, r, true) // true means try to cache this response
}
