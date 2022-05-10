package main 

import (
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"github.com/tidwall/gjson"
	"sync"
)

type Hotel struct {
	Name, Locality, Country string
	Price int64
	StarRating int64
	ApproxDistanceFromAirport float64
}

func hotelThread(i, budget int, start, end, people string, roundTrip *RoundTrip, vacations []Vacation, channel chan ReturnChan, mu *sync.Mutex) {
	log.Printf("Getting longitude and latitude of %v\n", roundTrip.DestinationAirport)
	lat, long := getAirportCoords(roundTrip.DestinationAirport)
	log.Printf("Coordinate long: %v, lat: %v\n", lat, long)
	hotel := getHotel(budget, start, end, long, lat, people, mu)
	log.Printf("About to calculate cost")
	totalCost := calculateCost(hotel, *roundTrip, start, end)
	log.Printf("Hotel: %v", hotel)
	log.Printf("About to initialize vacation object")
	vacation := Vacation{hotel, *roundTrip, totalCost}
	channel <- ReturnChan{vacation, i}
	// vacations[i] = vacation
}

func buildHotelQuery(budget int, start, end, longitude, latitude, people string) string {
	return fmt.Sprintf("%v?&budget=%v&start=%v&end=%v&latitude=%v&longitude=%v&people=%v", HOTEL_SERVICE_ADDRESS, budget, start, end, latitude, longitude, people)
}

func decodeHotel(hotelS string) Hotel {
	v := gjson.Parse(hotelS)
	hotel := Hotel{}
	// log.Printf("Item: %v\n", v)
	hotel.Name = v.Get("Name").String()
	hotel.Locality = v.Get("Locality").String()
	hotel.Country = v.Get("Country").String()
	hotel.Price = v.Get("Price").Int()
	hotel.StarRating = v.Get("StarRating").Int()
	hotel.ApproxDistanceFromAirport = v.Get("DistanceFromAirport").Float()
	return hotel
}

func getHotel(budget int, start, end, longitude, latitude, people string, mu *sync.Mutex) Hotel {
	log.Printf("Getting hotel")
	hotelQuery := buildHotelQuery(budget, start, end, longitude, latitude, people)
	// log.Printf("Hotel query: %v", hotelQuery)
	// mu.Lock()
	resp, err := http.Get(hotelQuery)
	// mu.Unlock()
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Body: %v, query: %v", string(body), hotelQuery)
	return decodeHotel(string(body))
}