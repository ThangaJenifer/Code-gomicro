package main

import (
	"net/http"
	"os"
	"testing"
)

// Create a var testApp type of Config
var testApp Config

// Inside TestMain we will setup things want to be in place before the actual tests are run
func TestMain(m *testing.M) {
	//os.Exit(m.Run()) this statement will run our tests
	os.Exit(m.Run())
}

// get the json from website and copy it to create var jsonToReturn
var jsonToReturn = `
{
	"ts": 1654782060772,
	"tsj": 1654782056216,
	"date": "Jun 9th 2022, 09:40:56 am NY",
	"items": [
	  {
		"curr": "USD",
		"xauPrice": 1849,
		"xagPrice": 21.9115,
		"chgXau": -3.735,
		"chgXag": -0.1425,
		"pcXau": -0.2016,
		"pcXag": -0.6461,
		"xauClose": 1852.735,
		"xagClose": 22.054
	  }
	]
  }
`

// Create a type RoundTripFunc which has func with parameter request and returns pointer to *http.Response
type RoundTripFunc func(req *http.Request) *http.Response

// Method of RoundTripFunc type RoundTrip takes request and gives out response and error
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// Constructor function which takes RoundTripFunc variable and gives back http.client
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}
