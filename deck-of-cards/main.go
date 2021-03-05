package main

import "fmt"
import "main/deck"
// import "sort"

func main() {
  d := deck.New(nil)
//   fmt.Println(d)
//   sort.Slice(d, func (i, j int) bool {
// 	  return deck.Less(d[i], d[j])
//   })

  d = deck.Filter(d, func (c deck.Card) bool {
	  return c.Rank == deck.Ace || c.Suit == deck.Spade
  })

  fmt.Println(d)
}