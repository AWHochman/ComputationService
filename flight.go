package main

import (
	"github.com/tidwall/gjson"
	"log"
	"net/http"
	"io/ioutil"
	"fmt"
	"strconv"
)

type RoundTrip struct {
	Outbound, Inbound *Trip 
	DestinationAirport string
	Cost float64
}

type Trip struct {
	Flights []Flight
}

type Flight struct {
	LegNumber, Airline string 
	ArrivalAirport, ArriveAt, ArrivalTime string 
	DepartFrom, DepartTime, DepartAirport string 
}

// type Airport struct {
// 	Name, Latitude_deg, Longitude_deg string
// }


// http://localhost:1989/departAirport=LGA&departDate=2022-04-22&returnDate=2022-04-29&numTravelers=2&preference=tropical
func getFlight(start, end, people, home, preference, exclude, list string) []RoundTrip {
	query := buildFlightQuery(start, end, people, home, preference, exclude, list)
	log.Printf("Flight query built: %v", query)
	resp, err := http.Get(query)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return buildRoundTrips(string(body))
}

func buildRoundTrips(body string) []RoundTrip {
	roundTrips := make([]RoundTrip, 0)
	for _, v := range gjson.Parse(body).Array() {
		roundTrip := RoundTrip{}
		roundTrip.Cost = v.Get("TotalRoundTripPrice").Float()
		roundTrip.Outbound, roundTrip.DestinationAirport = buildTrip(v.Get("outboundTrip").Array(), true)
		roundTrip.Inbound, _ = buildTrip(v.Get("returnTrip").Array(), false)
		roundTrips = append(roundTrips, roundTrip)
	}
	return roundTrips
}

func buildTrip(trip []gjson.Result, outbound bool) (*Trip, string) {
	flights := make([]Flight, 0)
	var roundTripDestination string
	for _, v := range trip {
		f := Flight{}
		f.LegNumber = v.Get("legNumber").String()
		f.ArrivalAirport = v.Get("arrival.arrivalAirportCode").String()
		f.ArriveAt = v.Get("arrival.arriveAt").String()
		f.ArrivalTime = v.Get("arrival.arrivalTime").String()
		f.DepartFrom = v.Get("departure.departFrom").String()
		f.DepartTime = v.Get("departure.departTime").String()
		f.DepartAirport = v.Get("departure.departAirportCode").String()
		f.Airline = v.Get("airline").String()
		// log.Printf("Leg number: %v, num legs: %v\n", f.LegNumber, len(trip))
		if legNum, _ := strconv.Atoi(f.LegNumber); outbound && len(trip) == legNum{
			roundTripDestination = f.ArrivalAirport
			// log.Printf("roundTripDestination: %v\n", roundTripDestination)
		}
		flights = append(flights, f)
	}
	return &Trip{flights}, roundTripDestination
}

func getFlights() []Flight{
	flights := make([]Flight, 0)
	flights = append(flights, Flight{})
	return flights
}

func buildFlightQuery(start, end, people, home, preference, exclude, list string) string {
	return fmt.Sprintf("%v?departAirport=%v&departDate=%v&returnDate=%v&numTravelers=%v&preference=%v&exclude=%v&list=%v", FLIGHT_SERVICE_ADDRESS, home, start, end, people, preference, exclude, list)
}