package pg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
)

func Path() (string, error) {
	path, err := dropboxPath()
	if err != nil {
		return path, err
	}

	if path == "" {
		path = os.Getenv("HOME")
	}

	return fmt.Sprintf("%s/.pg", path), nil
}

// https://www.dropbox.com/help/desktop-web/find-folder-paths
func dropboxPath() (string, error) {
	var infoPath string
	switch runtime.GOOS {
	case "windows":
		infoPath = "%APPDATA%\\Dropbox\\info.json"
		if _, err := os.Stat(infoPath); err != nil {
			infoPath = "%LOCALAPPDATA%\\Dropbox\\info.json"
			if _, err := os.Stat(infoPath); err != nil {
				infoPath = ""
			}
		}
	default:
		infoPath = fmt.Sprintf("%s/.dropbox/info.json", os.Getenv("HOME"))
	}

	if infoPath == "" {
		return "", nil
	}

	_, err := os.Stat(infoPath)
	if err != nil && (err == os.ErrNotExist || os.IsNotExist(err)) {
		return "", nil
	}

	data, err := ioutil.ReadFile(infoPath)
	if err != nil {
		return "", err
	}

	var info dropboxInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return "", err
	}

	return info.Path(), nil
}

type dropboxInfo struct {
	Personal *pathInfo `json:"personal,omitempty"`
	Business *pathInfo `json:"business,omitempty"`
}

func (i *dropboxInfo) Path() string {
	if i.Personal != nil {
		return i.Personal.Path
	}

	if i.Business != nil {
		return i.Business.Path
	}

	return ""
}

type pathInfo struct {
	Host int    `json:"host"`
	Path string `json:"path"`
}
