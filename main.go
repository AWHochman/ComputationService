package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
	"os"
	"encoding/json"
	"io/ioutil"
	"strconv"
)

type Vacation struct {
	Lodging []Hotel 
	Transportation RoundTrip 
}

var airportToCords map[string]interface{}
var HOTEL_SERVICE_ADDRESS string 
var FLIGHT_SERVICE_ADDRESS string 
var LOCAL = true 

func init() {
	if LOCAL {
		HOTEL_SERVICE_ADDRESS = "http://localhost:8081/api/query-hotels"
		FLIGHT_SERVICE_ADDRESS = "http://localhost:1989"
	} else {
		HOTEL_SERVICE_ADDRESS = "https://hotel-service.azurewebsites.net/api/query-hotels"
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
	// router.GET("/api/compute", compute)
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

// /compute?home=SAF&budget=1000&start=2022-05-10&end=2022-05-17&people=2&preference=tropical
func compute(c *gin.Context) {
	budget := c.Query("budget") // shortcut for c.Request.URL.Query().Get("budget")
	start := c.Query("start")
	end := c.Query("end")
	home := c.Query("home")
	people := c.Query("people")
	preference := c.Query("preference")

	// log.Printf("Input data: budget = %v, start = %v, end = %v, startLocation = %v, people = %v", budget, start, end, home, people)
	
	log.Printf("Getting round trip\n")
	roundTrip := getFlight(start, end, people, home, preference)
	log.Printf("Round trip successfully aquired\n")

	log.Printf("Getting longitude and latitude of %v\n", roundTrip.DestinationAirport)
	lat, long := getAirportCoords(roundTrip.DestinationAirport)
	log.Printf("Coordinate long: %v, lat: %v\n", lat, long)
	
	budgetI, err := strconv.Atoi(budget)
	if err != nil {
		panic(err)
	}
	hotels := getHotels(budgetI, start, end, long, lat, people)
	vacation := Vacation{hotels, *roundTrip}
	c.PureJSON(http.StatusOK, vacation)
}