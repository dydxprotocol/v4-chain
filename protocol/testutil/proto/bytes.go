package proto

// MustFirst is used for returning the first value of the `Marshal`
// method on a protobuf. This will panic if the conversion fails.
func MustFirst(bytes []byte, err error) []byte {
	if err != nil {
		panic(err)
	}
	return bytes
}
