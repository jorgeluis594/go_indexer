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

func errorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func main() {
	zincRepository := initRepository()

	r := chi.NewRouter()
	r.Use(errorHandler)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		path := "./public/index.html"
		content, err := os.ReadFile(path)
		if err != nil {
			panic(err)
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
			response, err = zincRepository.GetAll(currentPage)
		} else {
			response, err = zincRepository.Search(q, currentPage)
		}

		if err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(response)
	})

	port := ":8080"
	log.Println("Server running on port: ", port)
	err := http.ListenAndServe(port, r)
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
