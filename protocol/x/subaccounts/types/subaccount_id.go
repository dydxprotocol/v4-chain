package types

import (
	"sort"
)

// ToStateKey returns a bytes representation of a SubaccountId for use as a state key.
// The key uses the proto marshaling of the object such that it can be unmarshalled in
// the same way if it needs to be.
func (id *SubaccountId) ToStateKey() []byte {
	b, err := id.Marshal()
	if err != nil {
		panic(err)
	}
	return b
}

// SortedSubaccountIds is type alias for `[]SubaccountId` which supports deterministic
// sorting. SubaccountIds are first ordered by string comparison
// of their `Owner`, followed by integer comparison of their
// `Number`. If two `SubaccountId` have equal Owners, and Numbers, they
// are assumed to be equal, and their sorted order is not deterministic.
type SortedSubaccountIds []SubaccountId

// The below methods are required to implement `sort.Interface` for sorting using the sort package.
var _ sort.Interface = SortedSubaccountIds{}

func (s SortedSubaccountIds) Len() int {
	return len(s)
}

func (s SortedSubaccountIds) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortedSubaccountIds) Less(i, j int) bool {
	si := s[i]
	sj := s[j]

	if si.Owner != sj.Owner {
		return si.Owner < sj.Owner
	}

	if si.Number != sj.Number {
		return si.Number < sj.Number
	}

	return false
}
