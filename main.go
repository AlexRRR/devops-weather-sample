package main

import (
	"encoding/json"
	owm "github.com/briandowns/openweathermap"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	station = NewWeatherStation()
)

type WeatherStation struct {
	token string
}

// NewWeatherStation returns a weather station with the Token
// set from the environment
func NewWeatherStation() WeatherStation {
	return WeatherStation{
		token: os.Getenv("OWM_API_KEY"),
	}
}

// WeatherSnap contains the current weather information
type WeatherSnap struct {
	Temperature float64 `json:"Temperature"`
	Humidity    int     `json:"Humidity"`
}

func (w *WeatherStation) getWeather(cityID int) *WeatherSnap {
	report, err := owm.NewCurrent("C", "en", w.token) // Celsius, english
	if err != nil {
		log.Fatalln(err)
	}
	report.CurrentByID(cityID)
	return &WeatherSnap{
		Temperature: report.Main.Temp,
		Humidity:    report.Main.Humidity,
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Still alive!")
}

func handleURLMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cityID, err := strconv.Atoi(vars["cityID"])
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("report for city: %v", cityID)
	currentWeather := station.getWeather(cityID)
	json.NewEncoder(w).Encode(currentWeather)
}

func main() {
	// city code for zurich 2657896
	var router = mux.NewRouter()
	router.HandleFunc("/healthcheck", healthCheck).Methods("GET")
	router.HandleFunc("/m/{cityID}", handleURLMessage).Methods("GET")
	log.Println("Running server port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
