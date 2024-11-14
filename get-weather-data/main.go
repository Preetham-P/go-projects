package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

const url string = "http://api.openweathermap.org/data/2.5/weather?APPID="

type openAPIkey struct {
	OpenWeatherAPIconfigKey string `json:OpenWeatherAPIconfigKey`
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func loadAPIConfig(fileName string) (openAPIkey, error) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		log.Println("there was an error", err)
	}
	var apiKey openAPIkey

	err = json.Unmarshal([]byte(file), &apiKey)

	if err != nil {
		return openAPIkey{}, err
	}

	return apiKey, nil
}

func main() {
	fmt.Println("Welcome to the weather App")
	apikey, err := loadAPIConfig(".apiconfig")
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/weather/", weatherHandler)
	http.HandleFunc("/", helloHandler)

	if err != nil {
		fmt.Println("there was an error", err)
	}

	fmt.Println("This is the apikey", apikey.OpenWeatherAPIconfigKey)

	http.ListenAndServe(":9090", nil)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from go\n"))
}

func query(city string) (weatherData, error) {
	apikey, err := loadAPIConfig(".apiconfig")
	if err != nil {
		fmt.Println("there was an error", err)
	}

	urlNew := url + apikey.OpenWeatherAPIconfigKey + "&q=" + city

	fmt.Println("urlNew ", urlNew)

	response, err := http.Get(urlNew)
	if err != nil {
		log.Println("there was an error fetching weather report", err)
	}

	defer response.Body.Close()

	var wd weatherData
	if err := json.NewDecoder(response.Body).Decode(&wd); err != nil {
		log.Println("Decoder error", err)
		return weatherData{}, err
	}

	return wd, nil
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	city := strings.SplitN(r.URL.Path, "/", 3)[2]

	data, err := query(city)

	if err != nil {
		log.Println("There was an error fetching weather data", err)
		w.Write([]byte("There was an error fetching weather data for the city"))
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)
}
