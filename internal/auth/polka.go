package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	apiToken := headers.Get("Authorization")
	if apiToken == "" {
		return "", fmt.Errorf("No 'Authorization' header")
	}
	token_string, ok := strings.CutPrefix(apiToken, "ApiKey ")
	if ok {
		return token_string, nil
	} else {
		return "", fmt.Errorf("Error parsing the token")
	}
}
