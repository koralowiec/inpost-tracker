package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const apiUrl = "https://api-shipx-pl.easypack24.net/v1/"
const trackingEndpoint = apiUrl + "tracking/"
const statusesEndpoint = apiUrl + "statuses/"

type Status struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Statuses []Status

func (statues Statuses) findMatchingStatus(status string) Status {
	for _, s := range statues {
		if s.Name == status {
			return s
		}
	}

	return Status{}
}

type StatusesResponse struct {
	Items Statuses `json:"items"`
}

type TrackingDetail struct {
	Status   Status
	DateTime time.Time `json:"datetime"`
}

func (t *TrackingDetail) UnmarshalJSON(data []byte) error {
	var v map[string]string
	if err := json.Unmarshal(data, &v); err != nil {
		return err

	}

	datetime, err := time.Parse(time.RFC3339, v["datetime"])
	if err != nil {
		log.Fatal(err)
	}
	t.DateTime = datetime

	s, err := GetStatuses()
	if err != nil {
		log.Fatal(err)
	}
	t.Status = s.findMatchingStatus(v["status"])

	return nil
}

type TrackingResponse struct {
	TrackingNumber  string           `json:"tracking_number"`
	TrackingDetails []TrackingDetail `json:"tracking_details"`
}

var statuses Statuses

func GetStatuses() (Statuses, error) {
	if statuses != nil {
		return statuses, nil
	}

	return fetchStatuses()
}

func fetchStatuses() ([]Status, error) {
	res, err := http.Get(statusesEndpoint)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resJson StatusesResponse
	if err := json.Unmarshal(data, &resJson); err != nil {
		return nil, err
	}
	statuses = resJson.Items

	return statuses, nil
}

type PackageNotFoundError struct {
	trackingNumber string
}

func (e *PackageNotFoundError) Error() string {
	return fmt.Sprintf("Inpost API error: package with number: %s not found", e.trackingNumber)
}

func GetTrackingInfo(trackingNumber string) (*TrackingResponse, error) {
	if statuses != nil {
		if _, err := GetStatuses(); err != nil {
			return nil, err
		}
	}

	url := trackingEndpoint + "/" + trackingNumber
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		err := &PackageNotFoundError{trackingNumber}
		return nil, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resJson TrackingResponse
	if err := json.Unmarshal(data, &resJson); err != nil {
		return nil, err
	}

	return &resJson, nil
}
