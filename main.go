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
	id, departure, arrival, arrivalLocation string 
	ticketCost int 
}

type Airport struct {
	name, latitude_deg, longitude_deg string
}

var airportToCords map[string]interface{}

func init() {
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
	// router := gin.Default()
	// router.GET("/compute", compute)
	// port := getPort()
	// log.Printf("About to listen on %s. Go to https://127.0.0.1%s/\n", port, port)
	// router.Run(port)
	fmt.Println(getAirportCoords("CDG"))
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
	budget := 500 
	start := "2022-04-10"
	end := "2022-04-17"
	startLocation := "NYC"
	log.Printf("Input data: budget = %v, start = %v, end = %v, startLocation = %v", budget, start, end, startLocation)

	

	c.String(http.StatusOK, "hi")
}

func getFlights() []Flight{
	flights := make([]Flight, 0)
	flights = append(flights, Flight{"sample-id", "2022-04-10-12:30pm", "2022-04-10-4:00pm", "CDG", 100})
	return flights
}