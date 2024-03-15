package types

func (id *VaultId) ToStateKey() []byte {
	b, err := id.Marshal()
	if err != nil {
		panic(err)
	}
	return b
}
