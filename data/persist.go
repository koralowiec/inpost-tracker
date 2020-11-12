package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

func LoadFileContent(filePath string) *FileContent {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	fileContent := FileContent{}

	err = json.Unmarshal(file, &fileContent)
	if err != nil {
		log.Fatal(err)
	}

	return &fileContent
}

func (f FileContent) SaveFileContent(filePath string) error {
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))

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

func (f *FileContent) AppendTrackingNumber(number TrackingNumber) {
	f.TrackingNumbers = append(f.TrackingNumbers, number)
}
