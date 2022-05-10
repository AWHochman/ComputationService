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
	"time"
	"sync"
)

type Vacation struct {
	Lodging []Hotel 
	Transportation RoundTrip 
	TotalPrice []int
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
		FLIGHT_SERVICE_ADDRESS = "https://cloudflightservice.azurewebsites.net/api/QueryFlights"
		HOTEL_SERVICE_ADDRESS = "https://hotel-service.azurewebsites.net/api/query-hotels"
	}
	FLIGHT_SERVICE_ADDRESS = "https://cloudflightservice.azurewebsites.net/api/QueryFlights"
	HOTEL_SERVICE_ADDRESS = "https://hotel-service.azurewebsites.net/api/query-hotels"
	// airportToCords = make(map[string]interface{})
	plan, err := ioutil.ReadFile("Datasets/airports.json")
	if err != nil {
		log.Fatalln(err)
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
	budget := c.DefaultQuery("budget", "1000") // shortcut for c.Request.URL.Query().Get("budget")
	start := c.DefaultQuery("start", "-1")
	end := c.DefaultQuery("end", time.Now().Format("2006-01-02"))
	home := c.DefaultQuery("home", "-1")
	people := c.DefaultQuery("people", "1")
	preference := c.DefaultQuery("preference", "major")
	exclude := c.DefaultQuery("exclude", "[]")
	list := c.DefaultQuery("list", "true")

	if badInput(end, home) {
		log.Printf("BAD INPUT SUPPLIED")
		c.String(http.StatusOK, "Please specify an end date and the airport you would like to depart from")
		return 
	}
	

	// log.Printf("Input data: budget = %v, start = %v, end = %v, startLocation = %v, people = %v", budget, start, end, home, people)
	
	log.Printf("Getting round trips\n")
	roundTrips := getFlight(start, end, people, home, preference, exclude, list)
	log.Printf("Round trips successfully aquired\n")

	vacations := make([]Vacation, len(roundTrips))

	var wg sync.WaitGroup 
	for i, v := range roundTrips {
		wg.Add(1)
		budgetI, err := strconv.Atoi(budget)
		if err != nil {
			log.Fatalln(err)
		}

		go func(j int, v RoundTrip) {
			hotelThread(j, budgetI, start, end, people, &v, vacations)
			wg.Done()
		}(i, v)
	}
	wg.Wait()
	c.PureJSON(http.StatusOK, vacations)
}

func badInput(end, home string) bool {
	return end == "-1" || home == "-1"
}

func calculateCost(hotels []Hotel, transportation RoundTrip, start, end string) []int {
	totalPrices := make([]int, 0)
	log.Printf("here")
	tStart, err := time.Parse("2006-01-02", start)
	if err != nil {
		log.Fatalln(err)
	}
	tEnd, err := time.Parse("2006-01-02", end) 
	if err != nil {
		log.Fatalln(err)
	}
	numDays := int(tEnd.Sub(tStart).Hours()/24)
	for _, v := range hotels {
		totalPrices = append(totalPrices, numDays*int(v.Price) + int(transportation.Cost))
	}
	return totalPrices
}