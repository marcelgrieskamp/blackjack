package handlers

import (
	"blackjack/internal/game"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

var (
	// Spielzustände mit Mutex für Thread-Sicherheit
	games = make(map[string]*game.Game)
	mutex sync.RWMutex
)

type NewGameRequest struct {
	Reset bool `json:"reset"`
}

// Hilfsfunktion zum Abrufen oder Erstellen einer Session-ID
func getOrCreateSessionID(w http.ResponseWriter, r *http.Request) string {
	// Prüfen, ob bereits ein Cookie existiert
	sessionCookie, err := r.Cookie("blackjack_session")

	// Falls nicht, erstelle einen neuen
	if err != nil || sessionCookie.Value == "" {
		sessionID := uuid.New().String()
		cookie := &http.Cookie{
			Name:     "blackjack_session",
			Value:    sessionID,
			Path:     "/",
			MaxAge:   86400 * 30, // 30 Tage
			HttpOnly: true,
			Secure:   false, // Auf true setzen, wenn HTTPS verwendet wird
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, cookie)
		return sessionID
	}

	return sessionCookie.Value
}

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	// Stelle sicher, dass der Benutzer eine Session-ID hat
	getOrCreateSessionID(w, r)
	http.ServeFile(w, r, "static/index.html")
}

func NewGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Nur POST-Methode erlaubt", http.StatusMethodNotAllowed)
		return
	}

	sessionID := getOrCreateSessionID(w, r)

	var req NewGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err == nil && req.Reset {
		mutex.Lock()
		currentGame := game.NewGame() // Neues Spiel mit frischem Guthaben
		games[sessionID] = currentGame
		mutex.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(currentGame)
		return
	}

	mutex.Lock()
	currentGame, exists := games[sessionID]
	if exists {
		currentGame.Reset()
	} else {
		currentGame = game.NewGame()
		games[sessionID] = currentGame
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

	sessionID := getOrCreateSessionID(w, r)

	var betReq BetRequest
	if err := json.NewDecoder(r.Body).Decode(&betReq); err != nil {
		http.Error(w, "Ungültiger Einsatz", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	currentGame, exists := games[sessionID]
	if !exists {
		// Falls kein Spiel existiert, erstelle ein neues
		currentGame = game.NewGame()
		games[sessionID] = currentGame
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

	sessionID := getOrCreateSessionID(w, r)

	mutex.Lock()
	currentGame, exists := games[sessionID]
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

	sessionID := getOrCreateSessionID(w, r)

	mutex.Lock()
	currentGame, exists := games[sessionID]
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
