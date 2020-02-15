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

	"github.com/Sab94/go-udemy-dl/repo"
	"github.com/gosuri/uiprogress"
)

type Course struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type ListResponse struct {
	Next    int      `json:"next"`
	Results []Course `json:"results"`
}

type CourseData struct {
	Results []Item `json:"results"`
}

type DownloadObject struct {
	Chapter   string
	CourseId  int64
	LectureId int64
	//added to get current video number
	VideoNumber int64
	Attachments map[string]interface{}
	Videos      []interface{}
	Type        string
}

type Item struct {
	Title          string                 `json:"title"`
	Class          string                 `json:"_class"`
	Asset          map[string]interface{} `json:"asset"`
	IsDownloadable bool                   `json:"is_downloadable"`
	Id             int64                  `json:"id"`
	ObjectIndex    int64                  `json:"object_index"`
}

type VDO struct {
	File  string `json:"file"`
	Label string `json:"label"`
}

func (dl *Downloader) List() {
	session, err := repo.GetSession(dl.Root)
	if err != nil {
		log.Fatal(err.Error())
	}
	dl.Client.Jar.SetCookies(dl.BaseURL, session.Cookies)
	dl.CSRF = session.CSRF
	dl.AccessToken = session.AccessToken
	dl.ClientID = session.ClientID
	dl.BaseURL.Path = "/api-2.0/users/me/subscribed-courses"
	urlStr := dl.BaseURL.String()
	req, err := http.NewRequest("GET", urlStr+"?page_size=500", nil)
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
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	var j ListResponse
	_ = json.Unmarshal(data, &j)
	// l, _ := json.MarshalIndent(j, "", " ")
	k := j.Results
	// l := k[0].(map[string]interface{})
	// for i, j := range l {
	// 	log.Println(i, j)
	// }
	dl.fetchCource(k[0].ID)

}

func (dl *Downloader) fetchCource(id int64) {
	dl.BaseURL.Path = "/api-2.0/courses/" + fmt.Sprintf("%v", id) + "/cached-subscriber-curriculum-items"
	urlStr := dl.BaseURL.String()
	url := urlStr + "?page_size=1400&fields[lecture]=@min,object_index,asset,supplementary_assets,sort_order,is_published,is_free&fields[quiz]=@min,object_index,title,sort_order,is_published&fields[practice]=@min,object_index,title,sort_order,is_published&fields[chapter]=@min,description,object_index,title,sort_order,is_published&fields[asset]=@min,title,filename,asset_type,external_url,download_urls,stream_urls,length,status"
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
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	var j CourseData
	_ = json.Unmarshal(data, &j)

	var allVideosList []DownloadObject
	var resolutionChoices []string
	courseId := id
	chapter := ""
	for _, v := range j.Results {
		item := v
		if item.Class == "chapter" {
			chapter = fmt.Sprintf("%v", item.ObjectIndex) + " - " + item.Title
			continue
		}
		if item.Class == "lecture" {
			asset := item.Asset
			if asset["asset_type"].(string) == "Video" {
				var videos []interface{}
				if item.IsDownloadable {
					v := asset["download_urls"].(map[string]interface{})
					videos = v["Video"].([]interface{})
				} else {
					v := asset["stream_urls"].(map[string]interface{})
					videos = v["Video"].([]interface{})
				}
				objects := DownloadObject{
					Chapter:   chapter,
					CourseId:  courseId,
					LectureId: item.Id,
					//added to get current video number
					VideoNumber: item.ObjectIndex,
					Videos:      videos,
					Type:        "v",
				}
				for _, v := range videos {
					u := v.(map[string]interface{})
					resolutionChoices = append(resolutionChoices, u["label"].(string))
				}
				allVideosList = append(allVideosList, objects)
			}
		}
	}
	// log.Printf("%+v \n %+v", allVideosList, unique(resolutionChoices))

	for _, v := range allVideosList {
		dl.readyDownload(v, "720")
	}
}

func unique(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func (dl *Downloader) readyDownload(item DownloadObject, selectedResolution string) {
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
				dl.download(o.File, fmt.Sprintf("%v", videoNumber)+" - "+i.Title, ext)
			}
			videoNumber++
		}
	}
}

func (dl *Downloader) download(url, fileName, ext string) {
	dstFile := dl.Root + string(os.PathSeparator) + fileName + "." + ext
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
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal(err.Error())
	}
	// done <- n
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
