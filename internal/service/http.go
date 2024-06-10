package service

import (
	"main/pkg"
	"net/http"
	"time"
)

func newHttpClient(conf *pkg.HttpClientParams) http.Client {
	return http.Client{
		Timeout: time.Duration(conf.Timeout) * time.Second,
	}
}
