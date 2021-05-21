package own_deck

type Card struct {
	Name   CardName
	Family CardFamily
	Score  int
}

//go:generate stringer -type=CardName
type CardName int

//go:generate stringer -type=CardFamily
type CardFamily int

const (
	Spades  CardFamily = iota
	Diamons CardFamily = iota
	Clubs   CardFamily = iota
	Hearts  CardFamily = iota
)

const (
	_ CardName = iota
	A
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	J
	Q
	K
)

// categorias
// 1. SPADES 2.	DIAMONDS 3. CLUBS 4. HEARTS
// A,1,2,3,4...10,J,Q,K.

//permutations
//A:spades, 1:spades, 2:spades ... A diamons

func New() []Card {
	deck := make([]Card, 0)
	// cardFamilies := rand.Perm(4)
	// cardNames := rand.Perm(14)
	for i := 0; i < 4; i++ {
		for j := 0; j < 14; j++ {
			card := Card{Family: CardFamily(i), Name: CardName(j)}
			deck = append(deck, card)
		}
	}
	return deck
}
