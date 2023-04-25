package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type movie struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	Cast   string `json:"cast"`
	Year   int    `json:"year"`
}

var ListOfMovies []movie

func GetMovies(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(ListOfMovies)
	if err != nil {
		return
	}
}

func PostMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newMovie movie
	err := json.NewDecoder(r.Body).Decode(&newMovie)
	if err != nil {
		http.Error(w, "Ошибка добавления", http.StatusBadRequest)
		return
	}
	ListOfMovies = append(ListOfMovies, newMovie)
	err1 := json.NewEncoder(w).Encode(newMovie)
	if err1 != nil {
		return
	}
}

func GetMoviesByNumber(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	number := vars["number"]

	for _, a := range ListOfMovies {
		if strconv.Itoa(a.Number) == number {
			err := json.NewEncoder(w).Encode(a)
			if err != nil {
				return
			}
			return
		}
	}
	http.Error(w, "movie not found", http.StatusNotFound)
}

func Sleep(w http.ResponseWriter, _ *http.Request) {
	time.Sleep(10 * time.Second)
	err := json.NewEncoder(w).Encode(map[string]string{"message": "Задержка в 10 секунд завершена"})
	if err != nil {
		return
	}
}

func DeleteByNumber(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	number := vars["number"]

	for i, a := range ListOfMovies {
		if strconv.Itoa(a.Number) == number {
			ListOfMovies = append(ListOfMovies[:i], ListOfMovies[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "movie not delete", http.StatusNotFound)
}

func UpdateFilm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	number := vars["number"]

	for i, a := range ListOfMovies {
		if strconv.Itoa(a.Number) == number {
			err := json.NewDecoder(r.Body).Decode(&a)
			if err != nil {
				return
			}
			ListOfMovies[i] = a
			err1 := json.NewEncoder(w).Encode(a)
			if err1 != nil {
				return
			}
			return
		}
	}
	http.Error(w, "movies not found", http.StatusNotFound)
}

func GetGif(w http.ResponseWriter, r *http.Request) {
	fileName := "gif.gif"
	http.ServeFile(w, r, fileName)
}

/*func PostGif(w http.ResponseWriter, r *http.Request) {
	imageBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err = ioutil.WriteFile("example.gif", imageBytes, 0644)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err1 := json.NewEncoder(w).Encode(map[string]string{"message": "Image uploaded successfully"})
	if err1 != nil {
		return
	}
}*/

func PostGif1(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("gif")
	if err != nil {
		http.Error(w, "Failed to upload gif", http.StatusBadRequest)
		return
	}
	defer file.Close()

	f, err := os.OpenFile("./uploads/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		http.Error(w, "Failed to save gif", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, "Failed to save gif", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Gif uploaded successfully")
}

func main() {
	server := mux.NewRouter()
	server.HandleFunc("/movies", GetMovies).Methods(http.MethodGet)
	server.HandleFunc("/new_movie", PostMovies).Methods(http.MethodPost)
	server.HandleFunc("/movies/{number}", GetMoviesByNumber).Methods(http.MethodGet)
	server.HandleFunc("/sleep", Sleep).Methods(http.MethodGet)
	server.HandleFunc("/delete/{number}", DeleteByNumber).Methods(http.MethodDelete)
	server.HandleFunc("/update/{number}", UpdateFilm).Methods(http.MethodPut)
	server.HandleFunc("/gif", GetGif).Methods(http.MethodGet)
	//server.HandleFunc("/new_gif", PostGif).Methods(http.MethodPost)
	server.HandleFunc("/new_gif", PostGif1).Methods(http.MethodPost)
	err := http.ListenAndServe(":8080", server)
	if err != nil {
		return
	}
}
