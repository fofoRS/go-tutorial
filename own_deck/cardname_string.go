// Code generated by "stringer -type=CardName"; DO NOT EDIT.

package own_deck

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[A-0]
	_ = x[One-1]
	_ = x[Two-2]
	_ = x[Three-3]
	_ = x[Four-4]
	_ = x[Five-5]
	_ = x[Six-6]
	_ = x[Seven-7]
	_ = x[Eight-8]
	_ = x[Nine-9]
	_ = x[Ten-10]
	_ = x[J-11]
	_ = x[Q-12]
	_ = x[K-13]
}

const _CardName_name = "AOneTwoThreeFourFiveSixSevenEightNineTenJQK"

var _CardName_index = [...]uint8{0, 1, 4, 7, 12, 16, 20, 23, 28, 33, 37, 40, 41, 42, 43}

func (i CardName) String() string {
	if i < 0 || i >= CardName(len(_CardName_index)-1) {
		return "CardName(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _CardName_name[_CardName_index[i]:_CardName_index[i+1]]
}
