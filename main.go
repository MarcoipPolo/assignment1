package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Article struct {
	Title   string `json:"Title"`
	Desc    string `json:"Desc"`
	Content string `json:"Content"`
}

type Articles []Article

func AllArticles(rw http.ResponseWriter, r *http.Request) {
	articles := Articles{
		Article{Title: "Article 1", Desc: "Description 1", Content: " Content 1"},
		Article{Title: "Article 2", Desc: "Description 2", Content: " Content 2"},
		Article{Title: "Article 3", Desc: "Description 3", Content: " Content 3"},
	}
	fmt.Println("Endpoint Hint: All Articles Endpoint")

	rw.Header().Set("Content-Type", "application/json")

	json.NewEncoder(rw).Encode(articles)
}

func homePage(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, "Hello Prog2005, Here we go")
}

func handelRequests() {
	/// We have two endpoints, for the main root, like localhost:4747, it runs homepage function and for localhost:4747/articles it executes AllArticles function
	http.HandleFunc("/", homePage)
	http.HandleFunc("/articles", AllArticles)
	log.Fatal(http.ListenAndServe(getport(), nil))
}

func main() {
	handelRequests()
}

//// Get Port if it is set by environment, else use a defined one like "4747"
func getport() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "4747"
	}
	return ":" + port
}
