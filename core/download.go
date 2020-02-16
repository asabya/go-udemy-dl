package core

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gosuri/uiprogress"
)

func (dl *Downloader) readyDownload(item DownloadObject, selectedResolution, title string) {
	videoNumber := 1
	CourseId := item.CourseId
	LectureId := item.LectureId

	dl.BaseURL.Path = "/api-2.0/users/me/subscribed-courses/" + fmt.Sprintf("%v", CourseId) + "/lectures/" + fmt.Sprintf("%v", LectureId)
	urlStr := dl.BaseURL.String()
	url := urlStr + "?fields[asset]=@min,download_urls,stream_urls,external_url,slide_urls,captions,tracks&fields[course]=id,is_paid,url&fields[lecture]=@default,view_html,course&page_config=ct_v4"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
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
	data, _ := ioutil.ReadAll(resp.Body)
	var i Item
	_ = json.Unmarshal(data, &i)
	asset := i.Asset
	if asset["asset_type"].(string) == "Video" {
		var videos []VDO
		if i.IsDownloadable {
			v := asset["download_urls"].(map[string]interface{})
			a, _ := json.Marshal(v["Video"])
			json.Unmarshal(a, &videos)
		} else {
			v := asset["stream_urls"].(map[string]interface{})
			a, _ := json.Marshal(v["Video"])
			json.Unmarshal(a, &videos)
		}
		for _, o := range videos {
			if o.Label == selectedResolution {
				urlArr := strings.Split(o.File, "/")
				q := urlArr[len(urlArr)-1]
				p := strings.Split(q, ".")[1]
				ext := strings.Split(p, "?")[0]
				dl.download(o.File, fmt.Sprintf("%v", videoNumber)+" - "+i.Title, ext, title)
			}
			videoNumber++
		}
	}
}

func (dl *Downloader) download(url, fileName, ext, courseDir string) {
	location := dl.Root + string(os.PathSeparator) + courseDir
	_, err := os.Stat(location)
	if err != nil {
		err := os.Mkdir(location, 0775)
		if err != nil {
			log.Fatalf("Failed creating repo directory Err:%s", err.Error())
		}
	}

	dstFile := location + string(os.PathSeparator) + fileName + "." + ext
	log.Println(url, dstFile, ext)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
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
	resp, err := dl.Client.Head(url)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer resp.Body.Close()

	size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		log.Fatal(err.Error())
	}

	out, err := os.OpenFile(dstFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer out.Close()

	fInfo, err := out.Stat()
	if err == nil && fInfo.Size() == int64(size) {
		return
	}

	resp, err = dl.Client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer resp.Body.Close()
	done := make(chan int64)
	go func() {
		_ = printDownloadPercent(done, out, int64(size), fileName)
	}()
	n, err := io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal(err.Error())
	}
	done <- n
}

func printDownloadPercent(done chan int64, file *os.File, total int64, prepend string) error {
	uiprogress.Start()                              // start rendering
	bar := uiprogress.AddBar(100).AppendCompleted() // Add a new bar
	// prepend Downloading IPFS
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		return prepend
	})

	var stop bool = false
	for {
		select {
		case <-done:
			bar.Set(100)
			uiprogress.Stop()
			stop = true
		default:
			fi, err := file.Stat()
			if err != nil {
				return err
			}
			size := fi.Size()
			if size == 0 {
				size = 1
			}

			var percent float64 = float64(size) / float64(total) * 100
			bar.Set(int(percent))
		}
		if stop {
			break
		}
		time.Sleep(time.Second)
	}
	return nil
}
