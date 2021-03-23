package own_deck

type Card struct {
	name   CardName
	family CardFamily
	score  int
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
	A     CardName = iota
	One   CardName = iota
	Two   CardName = iota
	Three CardName = iota
	Four  CardName = iota
	Five  CardName = iota
	Six   CardName = iota
	Seven CardName = iota
	Eight CardName = iota
	Nine  CardName = iota
	Ten   CardName = iota
	J     CardName = iota
	Q     CardName = iota
	K     CardName = iota
)

// categorias
// 1. SPADES 2.	DIAMONDS 3. CLUBS 4. HEARTS
// A,1,2,3,4...10,J,Q,K.

//permutations
//A:spades, 1:spades, 2:spades ... A diamons

func New() []Card {
	deck := make([]Card, 52)
	// cardFamilies := rand.Perm(4)
	// cardNames := rand.Perm(14)
	for i := 0; i < 4; i++ {
		for j := 0; j < 14; j++ {
			card := Card{family: CardFamily(i), name: CardName(j)}
			deck = append(deck, card)
		}
	}
	return deck
}
