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
}

func buildHotelQuery(budget int, start, end, longitude, latitude, people string) string {
	return fmt.Sprintf("%v?&budget=%v&start=%v&end=%v&latitude=%v&longitude=%v&people=%v", HOTEL_SERVICE_ADDRESS, budget, start, end, latitude, longitude, people)
}

func decodeHotel(hotelS string) Hotel {
	hotel := Hotel{}
	hotel.Name = gjson.Get(hotelS, "Name").String()
	hotel.Locality = gjson.Get(hotelS, "Locality").String()
	hotel.Country = gjson.Get(hotelS, "Country").String()
	hotel.Price = gjson.Get(hotelS, "Price").Int()
	hotel.StarRating = gjson.Get(hotelS, "StarRating").Int()
	return hotel
}

func getHotels(budget int, start, end, longitude, latitude, people string) Hotel {
	log.Printf("Getting hotels")
	resp, err := http.Get(buildHotelQuery(budget, start, end, longitude, latitude, people))
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return decodeHotel(string(body))
}