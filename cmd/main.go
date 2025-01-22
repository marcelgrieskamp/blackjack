package main

import (
	"blackjack/internal/handlers"
	"log"
	"net/http"
)

func main() {
	// Statische Dateien servieren
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// API-Routen
	http.HandleFunc("/api/game/new", handlers.NewGame)
	http.HandleFunc("/api/game/bet", handlers.PlaceBet)
	http.HandleFunc("/api/game/hit", handlers.Hit)
	http.HandleFunc("/api/game/stand", handlers.Stand)

	// Index-Seite servieren
	http.HandleFunc("/", handlers.ServeIndex)

	log.Println("Server startet auf http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
