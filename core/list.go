package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Sab94/go-udemy-dl/repo"
	"github.com/manifoldco/promptui"
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

	k := j.Results
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F449 {{ .Title | cyan }})",
		Inactive: "  {{ .Title | cyan }}",
		Selected: "\U0001F449 {{ .Title | green | cyan }}",
	}

	prompt := promptui.Select{
		Label:     "Select Course",
		Items:     k,
		Size:      50,
		Templates: templates,
	}

	i, _, err := prompt.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Downloading : %s\n", k[i].Title)
	dl.fetchCource(k[i])
}

func (dl *Downloader) fetchCource(course Course) {
	dl.BaseURL.Path = "/api-2.0/courses/" + fmt.Sprintf("%v", course.ID) + "/cached-subscriber-curriculum-items"
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
	courseId := course.ID
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
	resolutions := unique(resolutionChoices)
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F449 {{ . | cyan }})",
		Inactive: "  {{ . | cyan }}",
		Selected: "\U0001F449 {{ . | green | cyan }}",
	}

	prompt := promptui.Select{
		Label:     "Select Course",
		Items:     resolutions,
		Size:      50,
		Templates: templates,
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
	for _, v := range allVideosList {
		dl.readyDownload(v, result, course.Title)
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
