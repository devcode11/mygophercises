package deck

import (
	"fmt"
)

type Suit uint8
type Rank uint8
type Card struct {
	Suit
	Rank
}

const (
	Joker Suit = iota
	Spade
	Diamond
	Club
	Heart
)

const (
	Ace Rank = iota + 1
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

func (s Suit) String() string {
	switch s {
		case Joker: return "Joker"
		case Heart: return "Heart"
		case Spade: return "Spade"
		case Diamond: return "Diamond"
		case Club: return "Club"
		default: return string(uint8(s))
	}
}

func (r Rank) String() string {
	switch r {
		case Ace: return "Ace"
		case Two: return "Two"
		case Three: return "Three"
		case Four: return "Four"
		case Five: return "Five"
		case Six: return "Six"
		case Seven: return "Seven"
		case Eight: return "Eight"
		case Nine: return "Nine"
		case Ten: return "Ten"
		case Jack: return "Jack"
		case Queen: return "Queen"
		case King: return "King"
		default: return string(uint8(r))
	}
}

func (c Card) String() string {
	return fmt.Sprintf("%s of %ss", c.Rank, c.Suit)
}