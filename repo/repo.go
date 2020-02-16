package repo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"log"
)

var pkgLock sync.Mutex

type Session struct {
	Business    string
	BaseURL     string
	Username    string
	AccessToken string
	CSRF        string
	Cookies     []*http.Cookie
	ClientID    string
}

func Init(repoPath, username, clientId, accessToken, csrf, baseURL, business string, cookies []*http.Cookie) error {
	pkgLock.Lock()
	defer pkgLock.Unlock()

	// Remove old repo
	RemoveRepo(repoPath)

	// Init repo
	err := initRepo(repoPath)
	if err != nil {
		return err
	}

	// Create session file
	err = createSession(repoPath)
	if err != nil {
		return err
	}

	session := Session{
		Username:    username,
		ClientID:    clientId,
		Cookies:     cookies,
		AccessToken: accessToken,
		BaseURL:     baseURL,
		CSRF:        csrf,
		Business:    business,
	}

	// Write to session
	err = writeSession(session, repoPath)
	if err != nil {
		return err
	}
	return nil
}

func IsInitialized(repoPath string) bool {
	return isInitialized(repoPath)
}

func isInitialized(repoPath string) bool {
	if _, e := os.Stat(repoPath); e != nil {
		return false
	}

	return true
}

func initRepo(repoPath string) error {
	// Check if directory is initialized
	if _, e := os.Stat(repoPath); e != nil && !os.IsNotExist(e) {
		log.Fatalf("Incorrect error in opening repo Err:%s", e.Error())
		return e
	} else if os.IsNotExist(e) {
		// Try to create the repo directory
		e = os.Mkdir(repoPath, 0775)
		if e != nil {
			log.Fatalf("Failed creating repo directory Err:%s", e.Error())
			return e
		}
	}
	return nil
}

func createSession(repoPath string) error {
	_, err := os.Create(repoPath + string(os.PathSeparator) + "session")
	if err != nil {
		return err
	}
	return nil
}

func RemoveRepo(repoPath string) error {
	// Check if directory is initialized
	if _, e := os.Stat(repoPath); e == nil && !os.IsNotExist(e) {
		// Try to remove the repo directory
		e = os.RemoveAll(repoPath)
		if e != nil {
			log.Fatalf("Failed removing repo directory Err:%s", e.Error())
			return e
		}
	}
	return nil
}

func writeSession(session Session, repoPath string) error {
	file, err := json.MarshalIndent(session, "", "\t")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(repoPath+string(os.PathSeparator)+"session", file, 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetSession(repoPath string) (*Session, error) {
	file, _ := ioutil.ReadFile(repoPath + string(os.PathSeparator) + "session")

	var data *Session

	err := json.Unmarshal([]byte(file), &data)

	return data, err
}
