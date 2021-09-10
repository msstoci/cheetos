package api

import (
	"net/http"
)

func init() {
	SetDefaultClient(&DefaultClient{
		Client:  &http.Client{},
		Timeout: 1,
	})
}

var defaultClient Client

func SetDefaultClient(c Client) {
	defaultClient = c
}

func GetDefaultClient() Client {
	return defaultClient
}
