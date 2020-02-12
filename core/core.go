package core

import (
	"context"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/user"
)

type Downloader struct {
	Root        string
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
	usr, err := user.Current()
	if err != nil {
		log.Printf("Failed getting user home directory. Is USER set?\n")
	}
	root := usr.HomeDir + string(os.PathSeparator) + ".gud"

	return &Downloader{
		BaseURL: u,
		Context: ctx,
		Cancel:  cancel,
		Client:  client,
		Root:    root,
	}
}

func (d *Downloader) SetBaseURL(url *url.URL) {
	d.BaseURL = url
}

func (d *Downloader) SetCSRF(csrf string) {
	d.CSRF = csrf
}
