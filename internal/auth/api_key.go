package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (apiKey string, err error) {
	apiKeyHeader := headers.Get("Authorization")
	if len(apiKeyHeader) == 0 {
		err = errors.New("no apiKey header found")
		return
	}

	var found bool
	apiKey, found = strings.CutPrefix(apiKeyHeader, "ApiKey ")
	if !found {
		err = errors.New("ApiKey prefix not found in header")
		return
	}

	return apiKey, nil
}
