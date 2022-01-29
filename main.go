package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Movie struct {
	MovieID   string `json:"movieid"`
	MovieName string `json:"moviename"`
}

type JsonResponse struct {
	Type    string  `json:"type"`
	Data    []Movie `json:"data"`
	Message string  `json:"message"`
}

const (
	HOST        = "localhost"
	DB_USER     = "postgres"
	DB_PASSWORD = "0203"
	DB_NAME     = "postgres"
	PORT        = "5432"
)

func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("host=%s user= %s dbname = %s sslmode= disable password = %s  port=%s ", HOST, DB_USER, DB_NAME, DB_PASSWORD, PORT)

	db, err := sql.Open("postgres", dbinfo)

	if err != nil {
		fmt.Println("sql connection is broken")
		panic(err)

	} else {
		fmt.Println("succesfully connected to DB")
	}

	return db
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/movies/", GetMovies).Methods("GET")

	router.HandleFunc("/movies/", CreateMovie).Methods("POST")

	router.HandleFunc("/movies/{movieid}", DeleteMovie).Methods("DELETE")

	fmt.Println("Served on 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func DeleteMovie(rw http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	movieID := params["movieid"]
	var response = JsonResponse{}

	if movieID == "" {
		response = JsonResponse{
			Type:    "error",
			Message: "MovieID is empty!",
		}

	} else {
		db := setupDB()
		printMessage("Deleting movie from db")
		_, err := db.Exec("delete  from movies where movieID=%s", movieID)
		checkErr(err)

		response = JsonResponse{
			Type:    "success",
			Message: "Succesfully deleted from movies",
		}
	}
	json.NewEncoder(rw).Encode(response)
}

func CreateMovie(rw http.ResponseWriter, r *http.Request) {
	movieID := r.FormValue("movieid")
	movieName := r.FormValue("moviename")

	var response = JsonResponse{}

	if movieID == "" && movieName == "" {
		response = JsonResponse{
			Type:    "error",
			Message: "You are missing movieID or movieName paramater. ",
		}
	} else {
		db := setupDB()
		printMessage("Inserting movie into DB")
		fmt.Printf("Inserting new movie with ID : %s and name : %s", movieID, movieName)
		var lastInsertID int
		err := db.QueryRow("Insert into movies(movieID,movieName) values (%v , %v) returning id;", movieID, movieName).Scan(&lastInsertID)

		if err != nil {
			printMessage("Couldn't add to DB")
			return
		}

		response = JsonResponse{
			Type:    "success",
			Message: "The movie added to DB",
		}

		json.NewEncoder(rw).Encode(response)
	}
}

func GetMovies(rw http.ResponseWriter, r *http.Request) {
	db := setupDB()
	printMessage("Getting movies")

	rows, err := db.Query("Select * from movies")
	checkErr(err)

	var movies []Movie

	for rows.Next() {
		var id int
		var movieID string
		var movieName string

		err = rows.Scan(&id, &movieID, &movieName)

		checkErr(err)

		movies = append(movies, Movie{
			MovieID:   movieID,
			MovieName: movieName,
		})

		var response = JsonResponse{
			Type:    "success",
			Data:    movies,
			Message: "Succesfully get movies from db",
		}
		json.NewEncoder(rw).Encode(response)
		return
	}

}

func printMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")

}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
