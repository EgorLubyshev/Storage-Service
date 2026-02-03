package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

const portNum string = ":8000"

//go:embed site/*
var siteFiles embed.FS

func main() {
	log.Println("Starting our simple http server.")

	site, err := fs.Sub(siteFiles, "site")
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", http.FileServer(http.FS(site)))

	log.Println("Started on http://localhost" + portNum)
	log.Println("To close connection CTRL+C")

	// Spinning up the server.
	err = http.ListenAndServe(portNum, nil)
	if err != nil {
		log.Fatal(err)
	}
}
