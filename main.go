package main

import (
	"fmt"
	"log"

	"gitlab.com/koralowiec/inpost-track/data"
)

const filePath = "/home/arek/.inpost-track/saved.json"

func main() {
	a := data.LoadFileContent()
	fmt.Println(a.TrackingNumbers)
	a.AppendTrackingNumber("2137")
	fmt.Println(a.TrackingNumbers)
	err := a.SaveFileContent()
	if err != nil {
		log.Fatal(err)
	}
}
