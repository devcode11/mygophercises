package deck

import (
	"math/rand"
	"time"
)

const deckLen = 52

type Deck []Card

type Options struct {
	Jokers uint8
	ExtraDecks uint8
}

func New(opt *Options) Deck {
	if opt == nil {
		opt = &Options{}
	}
	var cards Deck = make(Deck, 0, (1 + opt.ExtraDecks) * deckLen + opt.Jokers)

	for s := Spade; s <= Heart; s++ {
		for r := Ace; r <= King; r++ {
			cards = append(cards, Card{Suit: s, Rank: r})
		}
	}

	for i := 0; uint8(i) < opt.ExtraDecks; i++ {
		cards = append(cards, cards...)
	}

	for i := 0; uint8(i) < opt.Jokers; i++ {
		cards = append(cards, Card{Suit: Joker})
	}

	return cards
}

func Shuffle(d Deck) Deck {
	shuffled := make(Deck, len(d))
	rSrc := rand.NewSource(time.Now().Unix())
	for i, j := range rand.New(rSrc).Perm(len(d)) {
		shuffled[i] = d[j]
	}
	return shuffled
}

func Filter(d Deck, predicate func(c Card) bool) Deck {
	filtered := make(Deck, 0, len(d)/2)
	for _, c := range d {
		if predicate(c) {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

func Less(c, d Card) bool {
	switch {
	case c.Rank < d.Rank:
		return true
	case c.Rank == d.Rank:
		return c.Suit < d.Suit
	default:
		return false
	}
}
