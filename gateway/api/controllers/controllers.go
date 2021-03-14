package controllers

import (
	"net/http"
	"strings"
)

var allowedHeaders = []string{"X-", "Content-Type", "User-Agent", "Content-Length"}

func stripHeaders(headers http.Header) {
	for k := range headers {
		found := false
		for _, prefix := range allowedHeaders {
			if strings.Contains(k, prefix) {
				found = true
				break
			}
		}

		if !found {
			headers.Del(k)
		}
	}
}

func copyHeaders(destination http.Header, source *http.Header) {
	for k, v := range *source {
		vClone := make([]string, len(v))
		copy(vClone, v)
		(destination[k]) = v
	}
}
