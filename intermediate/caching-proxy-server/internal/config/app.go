package config

import (
	"flag"
	"log"
	"net/url"
	"time"
)

var (
	// ServerPort The port the caching server will listen to
	ServerPort string
	// OriginServerURL is the parsed URL of the upstream server we are caching for
	OriginServerURL *url.URL
	// CacheTTL is the time to live for the cache
	CacheTTL time.Duration
	// ClearCache, if true, will clear the existing caches
	ClearCache bool
)

// ParseFlags parses command line flags and validates them
func ParseFlags() {
	portFlg := flag.String("port", "8800", "The port the caching server will listen to")
	originServerFlg := flag.String("origin", "https://dummyjson.com", "The upstream server we are caching for")
	cacheTTLFlg := flag.Int("ttl", 300, "Cache duration in seconds")
	clearCacheFlg := flag.Bool("clear-cache", false, "Clear the cache and exit")

	flag.Parse()

	ServerPort = *portFlg
	ClearCache = *clearCacheFlg

	// Only validate origin server if clear-cache flag is not set
	if ClearCache {
		return
	}

	parsedOriginServerURL, err := url.Parse(*originServerFlg)
	if err != nil || parsedOriginServerURL.Scheme == "" || parsedOriginServerURL.Host == "" {
		log.Fatalf("Invalid origin server URL provided ('%s'): %v. Must be a complete URL (e.g., http://example.com).", *originServerFlg, err)
	}

	OriginServerURL = parsedOriginServerURL

	if *cacheTTLFlg <= 0 {
		log.Fatalf("Invalid cache TTL provided ('%d'). Must be a positive integer.", *cacheTTLFlg)
	}

	CacheTTL = time.Duration(*cacheTTLFlg) * time.Second
}
