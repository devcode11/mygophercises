package deck

import "testing"

func TestNew(t *testing.T) {
	cards := New(nil)
	if len(cards) != 52 {
		t.Errorf("Deck length is %d, should be 52", len(cards))
	}
}