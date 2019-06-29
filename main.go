package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/status", statusHandler)
	router.HandleFunc("/{digest:[0-9a-f]+}", digestHandler)
	router.HandleFunc("/", indexHandler)

	http.Handle("/", router)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err.Error())
	}

}
