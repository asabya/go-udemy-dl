package core

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (dl *Downloader) GetLogin() {
	dl.BaseURL.Path = "/join/login-popup/"
	urlStr := dl.BaseURL.String()
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	req.Header.Set("User-Agent", "StackOverflow")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Authorization", "Basic YWQxMmVjYTljYmUxN2FmYWM2MjU5ZmU1ZDk4NDcxYTY6YTdjNjMwNjQ2MzA4ODI0YjIzMDFmZGI2MGVjZmQ4YTA5NDdlODJkNQ==")
	req.Header.Set("Host", "www.udemy.com")
	req.Header.Set("Referer", "https://www.udemy.com/join/login-popup")
	req.Header.Set("Origin", "https://www.udemy.com")
	req.Header.Set("Accept", "application/json")
	resp, err := dl.Client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer resp.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find("input").Each(func(i int, s *goquery.Selection) {
		name, _ := s.Attr("name")
		if name == "csrfmiddlewaretoken" {
			csrf, _ := s.Attr("value")
			dl.SetCSRF(csrf)
		}
	})
	log.Println(dl.CSRF)
}

func (dl *Downloader) DoLogin(email, password string) {
	dl.BaseURL.Path = "/join/login-popup/"
	urlStr := dl.BaseURL.String()
	values := url.Values{}
	values.Set("csrfmiddlewaretoken", dl.CSRF)
	values.Set("email", email)
	values.Set("password", password)
	values.Set("locale", "en_US")
	reqPOST, err := http.NewRequest("POST", urlStr, strings.NewReader(values.Encode()))
	if err != nil {
		log.Fatal(err.Error())
	}
	reqPOST.Header.Set("User-Agent", "StackOverflow")
	reqPOST.Header.Set("X-Requested-With", "XMLHttpRequest")
	reqPOST.Header.Set("Authorization", "Basic YWQxMmVjYTljYmUxN2FmYWM2MjU5ZmU1ZDk4NDcxYTY6YTdjNjMwNjQ2MzA4ODI0YjIzMDFmZGI2MGVjZmQ4YTA5NDdlODJkNQ==")
	reqPOST.Header.Set("Host", "www.udemy.com")
	reqPOST.Header.Set("Referer", "https://www.udemy.com/join/login-popup")
	reqPOST.Header.Set("Origin", "https://www.udemy.com")
	reqPOST.Header.Set("Accept", "application/json")
	reqPOST.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	respost, err := dl.Client.Do(reqPOST)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer respost.Body.Close()
	for _, v := range dl.Client.Jar.Cookies(dl.BaseURL) {
		if v.Name == "access_token" {
			dl.AccessToken = v.Value
			log.Println(v.Value)
		} else if v.Name == "client_id" {
			dl.ClientID = v.Value
		}
	}
}
