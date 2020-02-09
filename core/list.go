package core

import (
	"io/ioutil"
	"log"
	"net/http"
)

func (dl *Downloader) List() {
	dl.BaseURL.Path = "/api-2.0/users/me/subscribed-courses"
	urlStr := dl.BaseURL.String()
	req, err := http.NewRequest("GET", urlStr+"?page_size=500", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(urlStr)
	req.Header.Set("User-Agent", "StackOverflow")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("X-Udemy-Bearer-Token", dl.AccessToken)
	req.Header.Set("X-Udemy-Client-Id", dl.ClientID)
	req.Header.Set("X-Udemy-Cache-User", dl.ClientID)
	req.Header.Set("Authorization", "Bearer "+dl.AccessToken)
	req.Header.Set("X-Udemy-Authorization", "Bearer "+dl.AccessToken)
	req.Header.Set("Host", "www.udemy.com")
	req.Header.Set("Referer", "https://www.udemy.com/join/login-popup")
	req.Header.Set("Origin", "https://www.udemy.com")
	req.Header.Set("Accept", "application/json")
	resp, err := dl.Client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	log.Println(resp.Status)
	log.Println(string(data))
}
