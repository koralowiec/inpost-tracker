package data

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type TrackingNumber string

type FileContent struct {
	TrackingNumbers []TrackingNumber `json:"tracking_numbers"`
}

func GetContentFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return home + "/.inpost-track/saved.json", nil
}

func createEmptyJsonFile(filePath string) (*FileContent, error) {
	dir := filepath.Dir(filePath)
	os.MkdirAll(dir, 0777)
	_, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	fileContent := FileContent{}
	fileContent.saveFileContent(filePath)

	return &fileContent, nil
}

func LoadFileContent(filePath string) (*FileContent, error) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		_, err = createEmptyJsonFile(filePath)
	}

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	fileContent := FileContent{}

	err = json.Unmarshal(file, &fileContent)
	if err != nil {
		return nil, err
	}

	return &fileContent, nil
}

func (f FileContent) saveFileContent(filePath string) error {
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func AppendTrackingNumber(filePath string, number TrackingNumber) ([]TrackingNumber, error) {
	_, err := os.Stat(filePath)
	var f *FileContent
	if err == nil {
		f, err = LoadFileContent(filePath)
		if err != nil {
			return nil, err
		}
	} else if os.IsNotExist(err) {
		f, err = createEmptyJsonFile(filePath)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	f.TrackingNumbers = append(f.TrackingNumbers, number)
	err = f.saveFileContent(filePath)
	if err != nil {
		return nil, err
	}

	return f.TrackingNumbers, nil
}
