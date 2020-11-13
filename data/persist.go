package data

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type TrackingNumber string

type FileContent struct {
	TrackingNumbers []TrackingNumber `json:"tracking_numbers"`
}

// grabbed from: https://dev.to/christalib/append-data-to-json-in-go-5gbj
func checkFile(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		_, err := os.Create(filename)
		if err != nil {
			return err

		}
	}
	return nil
}

func LoadFileContent(filePath string) (*FileContent, error) {
	err := checkFile(filePath)
	if err != nil {
		return nil, err
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

func (f FileContent) SaveFileContent(filePath string) error {
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}

	err = checkFile(filePath)
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
	f, err := LoadFileContent(filePath)
	if err != nil {
		return nil, err
	}

	f.TrackingNumbers = append(f.TrackingNumbers, number)
	err = f.SaveFileContent(filePath)
	if err != nil {
		return nil, err
	}

	return f.TrackingNumbers, nil
}
