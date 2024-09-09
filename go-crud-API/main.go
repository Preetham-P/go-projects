package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	ISBN     string    `json:"isbn"`
	Director *Director `json:"director"`
}

type Director struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var movies []Movie

func seedmovies() {
	movies = append(movies, Movie{ID: "1", ISBN: "23432434", Title: "KGF Chapter 1", Director: &Director{ID: "1", Name: "Prashant Neel"}})
	movies = append(movies, Movie{ID: "2", ISBN: "23432435", Title: "KGF Chapter 2", Director: &Director{ID: "1", Name: "Prashant Neel"}})
	movies = append(movies, Movie{ID: "3", ISBN: "23432436", Title: "Maharaja", Director: &Director{ID: "2", Name: "Some guy"}})
	movies = append(movies, Movie{ID: "4", ISBN: "23432437", Title: "Mr and MrsRamachari", Director: &Director{ID: "3", Name: "Santhosh Anadram"}})
	movies = append(movies, Movie{ID: "5", ISBN: "23432438", Title: "Masterpiece", Director: &Director{ID: "4", Name: "Dr Suri"}})
	movies = append(movies, Movie{ID: "6", ISBN: "23432439", Title: "Toxic", Director: &Director{ID: "5", Name: "Geetu Mohandas"}})
}

func getmovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func getmovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var movie Movie

	_ = json.NewDecoder(r.Body).Decode(&movie)

	movie.ID = strconv.Itoa(rand.IntN((10000000)))

	movies = append(movies, movie)

	json.NewEncoder(w).Encode(movie)
}

func main() {
	fmt.Println("Hi this is from CRUD API")
	seedmovies()
	r := mux.NewRouter()
	r.HandleFunc("/movies", getmovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getmovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	//r.HandleFunc("/movies/{id}".updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Println("starting the server on the port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
