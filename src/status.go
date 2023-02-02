package src

import (
	"encoding/json"
	"os"
	"time"
)

type StatusEntry struct {
	Timestamp time.Time
	Url       string
	Success   bool
	Err       string
}

func MakeStatusEntry(url string, err error) StatusEntry {
	var errorText string
	if err != nil {
		errorText = err.Error()
	}

	return StatusEntry{
		Timestamp: time.Now().UTC(),
		Url:       url,
		Success:   err == nil,
		Err:       errorText,
	}
}

func SaveStatus(path string, status []StatusEntry) error {
	data, err := json.Marshal(status)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, data, 0644)
	return err
}
