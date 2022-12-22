package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/jorgeluis594/go_indexer/repository"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	zincRepository := initRepository()
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		path := "./public/index.html"
		content, err := os.ReadFile(path)
		if err != nil {
			log.Fatal("Error reading ", path)
		}

		fmt.Fprintf(w, string(content))
	})

	r.Get("/api/search", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		q := r.URL.Query().Get("q")
		page := r.URL.Query().Get("page")

		if page == "" {
			page = "1"
		}
		currentPage, err := strconv.Atoi(page)
		if err != nil {
			currentPage = 1
		}

		var response *repository.SearchResponse
		if q == "" {
			response = zincRepository.GetAll(currentPage)
		} else {
			response = zincRepository.Search(q, currentPage)
		}

		json.NewEncoder(w).Encode(response)
	})

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}

func initRepository() repository.Repository {
	host := flag.String("host", "", "host of Zinc Search client")
	username := flag.String("username", "", "username of db")
	password := flag.String("password", "", "password of db")
	flag.Parse()

	clientHttp := repository.InitHttpClient(*host, *username, *password)
	return repository.InitRepository(clientHttp, "email_copy")
}
