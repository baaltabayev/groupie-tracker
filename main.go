package main

import (
	f "groupie-tracker/cmd"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", f.HomeHandler)
	mux.HandleFunc("/artist/", f.ArtistPage)
	mux.HandleFunc("/filters/", f.Filters)
	mux.HandleFunc("/search/", f.Search)
	mux.Handle("/ui/static/", http.StripPrefix("/ui/static/", http.FileServer(http.Dir("ui/static"))))
	log.Println("Запуск веб-сервера на http://localhost:8070/ ")
	err := http.ListenAndServe(":8070", mux)
	if err != nil {
		log.Fatal(err)
	}
}
