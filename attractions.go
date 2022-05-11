package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"io/ioutil"
	"github.com/tidwall/gjson"
)

type Attraction struct {
	Price, RecommendationScore, Title string
}

var ATTRACTION_ADDRESS string = "https://cloudflightservice.azurewebsites.net/api/queryattractions"

func getAttractions(roundTrip *RoundTrip, startDate, endDate string) []Attraction {
	lenOutbound := len(roundTrip.Outbound.Flights)
	location := roundTrip.Outbound.Flights[lenOutbound-1].ArriveAt
	query := buildAttractionQuery(strings.ReplaceAll(location, " ", ""), startDate, endDate)
	log.Printf("Built attraction query: %v", query)
	resp, err := http.Get(query)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return buildAttraction(string(body))
}

func buildAttraction(body string) []Attraction {
	attractions := make([]Attraction, 0)
	for _, v := range gjson.Parse(body).Array() {
		a := Attraction{}
		a.Price = v.Get("price").String()
		a.RecommendationScore = v.Get("recommendationScore").String()
		a.Title = v.Get("title").String()
		attractions = append(attractions, a)
	}
	return attractions
}

func buildAttractionQuery(location, startDate, endDate string) string {
	return fmt.Sprintf("%v?location=%v&startDate=%v&endDate=%v", ATTRACTION_ADDRESS, location, startDate, endDate)
}
