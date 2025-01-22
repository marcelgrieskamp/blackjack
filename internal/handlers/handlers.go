package handlers

import (
	"blackjack/internal/game"
	"encoding/json"
	"net/http"
	"sync"
)

var (
	// Spielzustände mit Mutex für Thread-Sicherheit
	games = make(map[string]*game.Game)
	mutex sync.RWMutex
)

type NewGameRequest struct {
	Reset bool `json:"reset"`
}

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func NewGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Nur POST-Methode erlaubt", http.StatusMethodNotAllowed)
		return
	}

	var req NewGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err == nil && req.Reset {
		mutex.Lock()
		currentGame := game.NewGame() // Neues Spiel mit frischem Guthaben
		games["game1"] = currentGame
		mutex.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(currentGame)
		return
	}

	mutex.Lock()
	currentGame, exists := games["game1"]
	if exists {
		currentGame.Reset()
	} else {
		currentGame = game.NewGame()
		games["game1"] = currentGame
	}
	mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentGame)
}

type BetRequest struct {
	Amount int `json:"amount"`
}

func PlaceBet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Nur POST-Methode erlaubt", http.StatusMethodNotAllowed)
		return
	}

	var betReq BetRequest
	if err := json.NewDecoder(r.Body).Decode(&betReq); err != nil {
		http.Error(w, "Ungültiger Einsatz", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	currentGame, exists := games["game1"]
	if !exists {
		mutex.Unlock()
		http.Error(w, "Kein aktives Spiel gefunden", http.StatusNotFound)
		return
	}

	if !currentGame.PlaceBet(betReq.Amount) {
		mutex.Unlock()
		http.Error(w, "Einsatz nicht möglich", http.StatusBadRequest)
		return
	}

	mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentGame)
}

func Hit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Nur POST-Methode erlaubt", http.StatusMethodNotAllowed)
		return
	}

	mutex.Lock()
	currentGame, exists := games["game1"]
	if !exists {
		mutex.Unlock()
		http.Error(w, "Kein aktives Spiel gefunden", http.StatusNotFound)
		return
	}

	currentGame.Hit()
	mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentGame)
}

func Stand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Nur POST-Methode erlaubt", http.StatusMethodNotAllowed)
		return
	}

	mutex.Lock()
	currentGame, exists := games["game1"]
	if !exists {
		mutex.Unlock()
		http.Error(w, "Kein aktives Spiel gefunden", http.StatusNotFound)
		return
	}

	currentGame.Stand()
	mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentGame)
}
