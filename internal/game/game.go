package game

import (
	"math/rand"
	"time"
)

type Card struct {
	Suit   string `json:"suit"`
	Value  string `json:"value"`
	Points int    `json:"points"`
}

type Game struct {
	PlayerHand []Card `json:"playerHand"`
	DealerHand []Card `json:"dealerHand"`
	Deck       []Card `json:"deck"`
	IsOver     bool   `json:"isOver"`
	Message    string `json:"message"`
	Chips      int    `json:"chips"`
	CurrentBet int    `json:"currentBet"`
	CanBet     bool   `json:"canBet"`
}

var suits = []string{"Herz", "Karo", "Pik", "Kreuz"}
var values = []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "Bube", "Dame", "König", "Ass"}

func NewDeck() []Card {
	var deck []Card
	for _, suit := range suits {
		for i, value := range values {
			points := i + 2
			if points > 10 && points < 14 {
				points = 10
			} else if points == 14 {
				points = 11
			}
			deck = append(deck, Card{Suit: suit, Value: value, Points: points})
		}
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
	return deck
}

func NewGame() *Game {
	deck := NewDeck()
	game := &Game{
		PlayerHand: make([]Card, 0),
		DealerHand: make([]Card, 0),
		Deck:       deck,
		IsOver:     false,
		Message:    "Platzieren Sie Ihren Einsatz",
		Chips:      1000,
		CurrentBet: 0,
		CanBet:     true,
	}
	return game
}

func (g *Game) Reset() {
	g.PlayerHand = make([]Card, 0)
	g.DealerHand = make([]Card, 0)
	g.Deck = NewDeck()
	g.IsOver = false
	g.Message = "Platzieren Sie Ihren Einsatz"
	g.CurrentBet = 0
	g.CanBet = true
}

func (g *Game) PlaceBet(amount int) bool {
	if !g.CanBet {
		return false
	}
	if amount > g.Chips {
		return false
	}
	g.CurrentBet = amount
	g.Chips -= amount
	g.CanBet = false

	// Erste zwei Karten austeilen
	g.PlayerHand = append(g.PlayerHand, g.Deck[0], g.Deck[1])
	g.DealerHand = append(g.DealerHand, g.Deck[2])
	g.Deck = g.Deck[3:]
	g.Message = "Spiel gestartet"

	return true
}

func (g *Game) CalculateScore(hand []Card) int {
	score := 0
	aces := 0

	for _, card := range hand {
		if card.Value == "Ass" {
			aces++
		} else {
			score += card.Points
		}
	}

	for i := 0; i < aces; i++ {
		if score+11 <= 21 {
			score += 11
		} else {
			score += 1
		}
	}

	return score
}

func (g *Game) Hit() {
	g.PlayerHand = append(g.PlayerHand, g.Deck[0])
	g.Deck = g.Deck[1:]

	score := g.CalculateScore(g.PlayerHand)
	if score > 21 {
		g.IsOver = true
		g.Message = "Bust! Sie haben verloren!"
		g.CurrentBet = 0
	}
}

func (g *Game) Stand() {
	// Dealer zieht Karten bis mindestens 17 Punkte erreicht sind
	for g.CalculateScore(g.DealerHand) < 17 {
		g.DealerHand = append(g.DealerHand, g.Deck[0])
		g.Deck = g.Deck[1:]
	}

	dealerScore := g.CalculateScore(g.DealerHand)
	playerScore := g.CalculateScore(g.PlayerHand)

	g.IsOver = true

	if dealerScore > 21 {
		g.Message = "Dealer bust! Sie haben gewonnen!"
		g.Chips += g.CurrentBet * 2
	} else if dealerScore > playerScore {
		g.Message = "Dealer gewinnt!"
		g.CurrentBet = 0
	} else if playerScore > dealerScore {
		g.Message = "Sie haben gewonnen!"
		g.Chips += g.CurrentBet * 2
	} else {
		g.Message = "Unentschieden!"
		g.Chips += g.CurrentBet
	}
	g.CurrentBet = 0
}
