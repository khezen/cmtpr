package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/google/go-github/v30/github"
)

func main() {
	eventPath := os.Getenv("GITHUB_EVENT_PATH")
	eventBytes, err := ioutil.ReadFile(eventPath)
	if err != nil {
		panic(err)
	}
	var event github.Event
	err = json.Unmarshal(eventBytes, &event)
	if err != nil {
		panic(err)
	}
	if len(os.Args) < 2 {
		panic(errors.New("no message trying to write on github PR"))
	}
	message = os.Args[1]

}
