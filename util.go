package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
)

type AcceptType struct {
	Types []string `json:"types"`
}

var acceptType *AcceptType

func initAcceptType() {
	bytes, err := ioutil.ReadFile("mime-type.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	var accept AcceptType
	err = json.Unmarshal(bytes, &accept)
	if err != nil {
		log.Fatal(err.Error())
	}

	acceptType = &accept
}

func (accept *AcceptType) String() string {
	var builder strings.Builder
	for _, value := range acceptType.Types {
		builder.WriteString(value)
	}
	return builder.String()
}
