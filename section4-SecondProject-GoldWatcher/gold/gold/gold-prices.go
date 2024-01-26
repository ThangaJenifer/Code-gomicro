package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// User can chose there own currency
var currency = "INR"

// This type will store the JSON data coming from external site to these fields type of Gold
type Gold struct {
	//Prices are coming as array item of slices has a JSON tag of items
	Prices []Price `json:"items"`
	//Creating pointer to http client because we want to write tests for this file
	Client *http.Client
}

// Create price struct and keep a structure which JSON Value we care about
type Price struct {
	Currency      string    `json:"currency"`
	Price         float64   `json:"xauPrice"`
	Change        float64   `json:"chgXau"`
	PreviousClose float64   `json:"xauClose"`
	Time          time.Time `json:"-"`
}

func (g *Gold) GetPrices() (*Price, error) {
	//First we need to call external client for getting json, so we check g.Client == nil which if not been set, then give a default value g.Client = &http.Client{} which can connect to outside world
	if g.Client == nil {
		g.Client = &http.Client{}
	}
	//create a client, url and request , use client.Do(req) to get response back
	client := g.Client
	//Copy the https://data-asg.goldprice.org/dbXRates/ API endpoint with placeholder for currency
	url := fmt.Sprintf("https://data-asg.goldprice.org/dbXRates/%s", currency)
	//Create a request
	req, _ := http.NewRequest("GET", url, nil)
	//get response by using client.Do(req)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("error contacting goldprice.org", err)
		return nil, err
	}
	defer resp.Body.Close()
	//read body of response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("error reading json", err)
		return nil, err
	}
	//Create variable of gold
	gold := Gold{}
	//Create three vars below of float64 type
	var previous, current, change float64
	//Unmarshalling body bytes of response into gold var
	err = json.Unmarshal(body, &gold)
	if err != nil {
		log.Println("error unmarshalling", err)
		return nil, err
	}
	//Assiging gold.Prices[0].PreviousClose, gold.Prices[0].Price, gold.Prices[0].Change
	//to created variable, as always we get first index
	previous, current, change = gold.Prices[0].PreviousClose, gold.Prices[0].Price, gold.Prices[0].Change
	//We will return the gold prices to this function , once we build it
	var currentInfo = Price{
		Currency:      currency,
		Price:         current,
		Change:        change,
		PreviousClose: previous,
		Time:          time.Now(),
	}

	return &currentInfo, nil
}
