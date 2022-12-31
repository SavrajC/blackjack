package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Card represents a playing card with a suit and rank
type Card struct {
	Suit string
	Rank int
}

func (c Card) String() string {
	// Use a map to convert the rank to a string representation
	rankStrings := map[int]string{
		1:  "A",
		11: "J",
		12: "Q",
		13: "K",
	}

	// If the rank is not in the map, use the integer value of the rank
	rankString, ok := rankStrings[c.Rank]
	if !ok {
		rankString = fmt.Sprintf("%d", c.Rank)
	}

	return rankString + c.Suit
}

// Deck represents a deck of cards
type Deck []Card

// Shuffle shuffles the deck using the Fisher-Yates shuffle algorithm
func (d Deck) Shuffle() {
	// Seed the random number generator with the current time
	rand.Seed(time.Now().UnixNano())

	for i := range d {
		j := rand.Intn(i + 1)
		d[i], d[j] = d[j], d[i]
	}
}

// Deal deals a number of cards from the top of the deck
func (d *Deck) Deal(n int) []Card {
	cards := (*d)[:n]
	*d = (*d)[n:]
	return cards
}

// Hand represents a hand of cards in a game of blackjack
type Hand struct {
	Cards []Card
	Total int
	Soft  bool
}

func (h Hand) String() string {
	cardStrings := make([]string, len(h.Cards))
	for i, card := range h.Cards {
		cardStrings[i] = card.String()
	}
	return strings.Join(cardStrings, ", ")
}

// Score calculates the total score of the hand, taking into account
// whether the hand contains an Ace that should be counted as 11 points
// (a "soft" hand) or 1 point (a "hard" hand)
func (h *Hand) Score() {
	h.Total = 0
	h.Soft = false
	for _, card := range h.Cards {
		if card.Rank == 1 {
			h.Total += 11
			h.Soft = true
		} else if card.Rank > 10 {
			h.Total += 10
		} else {
			h.Total += card.Rank
		}
	}
	if h.Total > 21 && h.Soft {
		h.Total -= 10
		h.Soft = false
	}
}

// Game represents a game of blackjack
type Game struct {
	Deck   Deck
	Dealer Hand
	Player Hand
	Done   bool
}

// NewGame creates a new game of blackjack with a shuffled deck of cards
func NewGame() *Game {
	deck := make(Deck, 52)
	suits := []string{"♠", "♥", "♦", "♣"}
	index := 0
	for _, suit := range suits {
		for rank := 1; rank <= 13; rank++ {
			deck[index] = Card{Suit: suit, Rank: rank}
			index++
		}
	}

	// Shuffle the deck
	deck.Shuffle()

	return &Game{
		Deck:   deck,
		Dealer: Hand{},
		Player: Hand{},
		Done:   false,
	}
}

// Hit deals a card to the hand and updates the total score
func (g *Game) Hit(h *Hand) {
	card := g.Deck.Deal(1)[0]
	h.Cards = append(h.Cards, card)
	h.Score()
}

// Stand ends the player's turn and starts the dealer's turn
func (g *Game) Stand() {
	g.Done = true
	for g.Dealer.Total < 17 || (g.Dealer.Total == 17 && g.Dealer.Soft) {
		g.Hit(&g.Dealer)
	}
}

// Result returns the result of the game as a string
func (g *Game) Result() string {

	if g.Dealer.Total > 21 {
		return "You win! Dealer bust."
	} else if g.Player.Total == 21 && g.Dealer.Total == 21 {
		return "You both hit blackjack! You Tie."
	} else if g.Player.Total > 21 {
		return "You lose! You bust."
	} else if g.Player.Total > 21 && g.Dealer.Total != 21 {
		return "You win! Dealer busts."
	} else if g.Dealer.Total > 21 && g.Dealer.Total != 21 {
		return "You win! Dealer busts."
	} else if g.Dealer.Total > g.Player.Total {
		return "You lose! Dealer has a higher score."
	} else if g.Dealer.Total < g.Player.Total {
		return "You win! You have a higher score."
	}
	return "Push. You have the same score as the dealer."
}

func main() {
	g := NewGame()

	// Deal the initial cards to the player and dealer
	g.Hit(&g.Player)
	g.Hit(&g.Dealer)
	g.Hit(&g.Player)
	g.Hit(&g.Dealer)

	for !g.Done {
		fmt.Println("Your hand:", g.Player, "\nTotal:", g.Player.Total)
		fmt.Println("Dealer's hand:", g.Dealer, "\nTotal:", g.Dealer.Total)
		fmt.Println("Enter 'h' to hit or 's' to stand:")
		var input string
		fmt.Scanln(&input)
		if input == "h" {
			g.Hit(&g.Player)
		} else if input == "s" {
			g.Stand()
		}
		if g.Player.Total > 21 {
			fmt.Println("Player Busts")
			break
		}
	}

	fmt.Println("Dealer's hand:", g.Dealer, "\nDealer Total:", g.Dealer.Total)
	fmt.Println("Your hand:", g.Player, "\nYour Total:", g.Player.Total)
	fmt.Println(g.Result())
}
