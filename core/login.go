package core

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/Sab94/go-udemy-dl/repo"
)

func (dl *Downloader) GetLogin() error {
	dl.BaseURL.Path = "/join/login-popup/"
	urlStr := dl.BaseURL.String()
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return err
	}
	dl.SetHeaders(req)
	resp, err := dl.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the csrf
	doc.Find("input").Each(func(i int, s *goquery.Selection) {
		name, _ := s.Attr("name")
		if name == "csrfmiddlewaretoken" {
			dl.CSRF, _ = s.Attr("value")
		}
	})

	return nil
}

func (dl *Downloader) DoLogin(email, password string) error {
	dl.BaseURL.Path = "/join/login-popup/"
	urlStr := dl.BaseURL.String()
	values := url.Values{}
	values.Set("csrfmiddlewaretoken", dl.CSRF)
	values.Set("email", email)
	values.Set("password", password)
	values.Set("locale", "en_US")
	reqPOST, err := http.NewRequest("POST", urlStr, strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}
	dl.SetHeaders(reqPOST)
	respost, err := dl.Client.Do(reqPOST)
	if err != nil {
		return err
	}
	defer respost.Body.Close()
	for _, v := range dl.Client.Jar.Cookies(dl.BaseURL) {
		if v.Name == "access_token" {
			dl.AccessToken = v.Value
		} else if v.Name == "client_id" {
			dl.ClientID = v.Value
		}
	}
	if dl.AccessToken == "" || dl.ClientID == "" {
		return errors.New("Please check credentials")
	}

	err = repo.Init(dl.Root, email, dl.ClientID, dl.AccessToken, dl.CSRF, dl.BaseURL.String(), dl.Business, dl.Client.Jar.Cookies(dl.BaseURL))
	if err != nil {
		return err
	}
	return nil
}
