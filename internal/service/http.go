package service

import (
	"main/pkg"
	"net/http"
)

func newHttpClient(conf *pkg.HttpClientParams) http.Client {
	return http.Client{
		Timeout: conf.Timeout,
	}
}
