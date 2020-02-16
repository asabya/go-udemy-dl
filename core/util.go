package core

import (
	"net/http"
)

func (dl *Downloader) SetHeaders(request *http.Request) {
	if dl.AccessToken == "" {
		request.Header.Set("User-Agent", "StackOverflow")
		request.Header.Set("X-Requested-With", "XMLHttpRequest")
		request.Header.Set("Authorization", "Basic YWQxMmVjYTljYmUxN2FmYWM2MjU5ZmU1ZDk4NDcxYTY6YTdjNjMwNjQ2MzA4ODI0YjIzMDFmZGI2MGVjZmQ4YTA5NDdlODJkNQ==")
		request.Header.Set("Host", "www.udemy.com")
		request.Header.Set("Referer", "https://www.udemy.com/join/login-popup")
		request.Header.Set("Origin", "https://www.udemy.com")
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		request.Header.Set("User-Agent", "StackOverflow")
		request.Header.Set("X-Requested-With", "XMLHttpRequest")
		request.Header.Set("X-Udemy-Bearer-Token", dl.AccessToken)
		request.Header.Set("X-Udemy-Client-Id", dl.ClientID)
		request.Header.Set("X-Udemy-Cache-User", dl.ClientID)
		request.Header.Set("Authorization", "Bearer "+dl.AccessToken)
		request.Header.Set("X-Udemy-Authorization", "Bearer "+dl.AccessToken)
		request.Header.Set("Host", "www.udemy.com")
		request.Header.Set("Referer", "https://www.udemy.com/join/login-popup")
		request.Header.Set("Origin", "https://www.udemy.com")
		request.Header.Set("Accept", "application/json")
	}
}
