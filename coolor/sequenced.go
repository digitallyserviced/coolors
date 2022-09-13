package coolor

type SequencedCausality interface {
	Before(Sequence) bool
	After(Sequence) bool
	Is(Sequence) bool
}

type Sequenced interface {
	Current() uint
}

type Sequence interface {
	Sequenced
	SequencedCausality
}

// type Sequence
type SequencedItem struct {
	seqNo uint
}

// After implements SequencedCausality
func (sqi *SequencedItem) After(s Sequence) bool {
	return sqi.Current() > s.Current()
}

// Before implements SequencedCausality
func (sqi *SequencedItem) Before(s Sequence) bool {
	return sqi.Current() < s.Current()
}

// Is implements SequencedCausality
func (sqi *SequencedItem) Is(s Sequence) bool {
	return sqi.Current() == s.Current()
}

func (o *SequencedItem) Current() uint {
	return o.seqNo
}

func (o *SequencedItem) Next() uint {
	o.seqNo++
	return o.seqNo
}
