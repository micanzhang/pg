package pg

import (
	"encoding/json"
	"io/ioutil"
)

type EntryPersistant interface {
	Restore() (entries []Entry, err error)
	Save(entries []Entry) error
}

type FileEntryPersistant struct {
	path string
}

func NewFileEntryPersistant(path string) EntryPersistant {
	return &FileEntryPersistant{
		path: path,
	}
}

func (f *FileEntryPersistant) Restore() (entries []Entry, err error) {
	data, err := ioutil.ReadFile(f.path)
	if err != nil || len(data) == 0 {
		return
	}

	err = json.Unmarshal(data, &entries)
	return
}

func (f *FileEntryPersistant) Save(entries []Entry) error {
	data, err := json.Marshal(entries)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(f.path, data, 0600)
	return err
}
