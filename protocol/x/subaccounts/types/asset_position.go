package types

// DeepCopy returns a deep copy of the asset position.
func (p *AssetPosition) DeepCopy() AssetPosition {
	b, err := p.Marshal()
	if err != nil {
		panic(err)
	}
	position := AssetPosition{}
	if err := position.Unmarshal(b); err != nil {
		panic(err)
	}
	return position
}
