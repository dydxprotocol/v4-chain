package types

// DeepCopy returns a deep copy of the perpetual position.
func (p *PerpetualPosition) DeepCopy() PerpetualPosition {
	b, err := p.Marshal()
	if err != nil {
		panic(err)
	}
	newPerpetualPosition := PerpetualPosition{}
	if err := newPerpetualPosition.Unmarshal(b); err != nil {
		panic(err)
	}
	return newPerpetualPosition
}
