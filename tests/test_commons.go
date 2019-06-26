package tests

import (
	"net/http"
	"net/url"
	"strings"
)

func NewPostRequestWithHeaders(url string, data url.Values, headers map[string]string) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	client := http.DefaultClient
	return client.Do(req)
}
