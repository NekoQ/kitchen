package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Cook struct {
	ID          int    `json:"id"`
	Rank        int    `json:"rank"`
	Proficiency int    `json:"proficiency"`
	Name        string `json:"name"`
	CatchPhrase string `json:"catch-phrase"`
}

func UnmarshalCooks() {

	cooksJson, err := os.Open("cooks.json")
	if err != nil {
		fmt.Println(err)
	}
	defer cooksJson.Close()

	byteValue, _ := ioutil.ReadAll(cooksJson)

	err = json.Unmarshal(byteValue, &Cooks)
}

func FillRankChans() {

	for i := range RankChans {
		RankChans[i] = make(chan Cook, 100)
	}

	for _, cook := range Cooks {
		for i := 0; i < cook.Proficiency; i++ {
			RankChans[cook.Rank-1] <- cook
		}
	}
}
