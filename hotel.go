package main 

import (
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"github.com/tidwall/gjson"
)

type Hotel struct {
	Name, Locality, Country string
	Price int64
	StarRating int64
	ApproxDistanceFromAirport float64
}

func buildHotelQuery(budget int, start, end, longitude, latitude, people string) string {
	return fmt.Sprintf("%v?&budget=%v&start=%v&end=%v&latitude=%v&longitude=%v&people=%v", HOTEL_SERVICE_ADDRESS, budget, start, end, latitude, longitude, people)
}

func decodeHotel(hotelS string) []Hotel {
	hotels := make([]Hotel, 0)
	// log.Printf("Hotel string: %v\n", hotelS)
	for _, v := range gjson.Parse(hotelS).Array() {
		h := Hotel{}
		// log.Printf("Item: %v\n", v)
		h.Name = v.Get("Name").String()
		h.Locality = v.Get("Locality").String()
		h.Country = v.Get("Country").String()
		h.Price = v.Get("Price").Int()
		h.StarRating = v.Get("StarRating").Int()
		h.ApproxDistanceFromAirport = v.Get("DistanceFromAirport").Float()
		hotels = append(hotels, h)
	}
	return hotels
}

func getHotels(budget int, start, end, longitude, latitude, people string) []Hotel {
	log.Printf("Getting hotels")
	hotelQuery := buildHotelQuery(budget, start, end, longitude, latitude, people)
	log.Printf("Hotel query: %v", hotelQuery)
	resp, err := http.Get(hotelQuery)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return decodeHotel(string(body))
}