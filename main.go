package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
	"os"
	"encoding/json"
	"io/ioutil"
)

type Flight struct {
	Id, Departure, Arrival, ArrivalLocation string 
	TicketCost int 
}

type Airport struct {
	Name, Latitude_deg, Longitude_deg string
}

var airportToCords map[string]interface{}
var FLIGHT_SERVICE_ADDRESS string 
var LOCAL = true 

func init() {
	if LOCAL {
		FLIGHT_SERVICE_ADDRESS = "http://localhost:8081/api/query-hotels"
	} else {
		FLIGHT_SERVICE_ADDRESS = "https://hotel-service.azurewebsites.net/api/query-hotels"
	}
	// airportToCords = make(map[string]interface{})
	plan, err := ioutil.ReadFile("Datasets/airports.json")
	if err != nil {
		panic(err)
	}
	var data interface{}
	err = json.Unmarshal(plan, &data)
	airportToCords = data.(map[string]interface{})
}

func main() {
	router := gin.Default()
	router.GET("/api/compute", compute)
	port := getPort()
	log.Printf("About to listen on %s. Go to https://127.0.0.1%s/\n", port, port)
	router.Run(port)
}

func getAirportCoords(code string) (string, string) {
	ap := airportToCords[code].(map[string]interface{})
	lat := fmt.Sprintf("%f", ap["latitude_deg"].(float64)) 
	long := fmt.Sprintf("%f", ap["longitude_deg"].(float64)) 
	return lat, long
}

func getPort() string {
    port := ":8080"
    if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
        port = ":" + val
    }
    return port
}

// /compute?home=SAF&budget=500&start=2022-03-26&end=2022-03-27
func compute(c *gin.Context) {
	// budget := c.Query("budget") // shortcut for c.Request.URL.Query().Get("budget")
	// start := c.Query("start")
	// end := c.Query("end")
	// home := c.Query("home")
	// people := c.Query("people")
	budget := 1500 
	start := "2022-04-10"
	end := "2022-04-17"
	startLocation := "JFK"
	people := 3
	log.Printf("Input data: budget = %v, start = %v, end = %v, startLocation = %v, people = %v", budget, start, end, startLocation, people)
	
	flights := getFlights()
	curFlight := flights[0]
	lat, long := getAirportCoords(curFlight.ArrivalLocation)
	
	hotels := getHotels(budget - flights[0].TicketCost, start, end, long, lat, people)
	c.String(http.StatusOK, hotels)
}

func getFlights() []Flight{
	flights := make([]Flight, 0)
	flights = append(flights, Flight{"sample-id", "2022-04-10-12:30pm", "2022-04-10-4:00pm", "BOS", 100})
	return flights
}

func buildHotelQuery(budget int, start, end, longitude, latitude string, people int) string {
	return fmt.Sprintf("%v?&budget=%v&start=%v&end=%v&latitude=%v&longitude=%v&people=%v", FLIGHT_SERVICE_ADDRESS, budget, start, end, latitude, longitude, people)
}

func getHotels(budget int, start, end, longitude, latitude string, people int) string {
	log.Printf("Getting hotels")
	resp, err := http.Get(buildHotelQuery(budget, start, end, longitude, latitude, people))
	if err != nil {
	log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return string(body)
}