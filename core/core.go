package core

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type Downloader struct {
	BaseURL     *url.URL
	Client      *http.Client
	Context     context.Context
	Cancel      context.CancelFunc
	CSRF        string
	AccessToken string
	ClientID    string
}

func New() *Downloader {
	ctx, cancel := context.WithCancel(context.Background())
	apiUrl := "https://www.udemy.com"
	u, _ := url.ParseRequestURI(apiUrl)
	jar, _ := cookiejar.New(&cookiejar.Options{})
	client := &http.Client{
		Jar: jar,
	}
	return &Downloader{
		BaseURL: u,
		Context: ctx,
		Cancel:  cancel,
		Client:  client,
	}
}

func (d *Downloader) SetBaseURL(url *url.URL) {
	d.BaseURL = url
}

func (d *Downloader) SetCSRF(csrf string) {
	d.CSRF = csrf
}
