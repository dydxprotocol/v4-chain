package eth

const (
	MinAddrLen = 20
	MaxAddrLen = 32
)

// PadOrTruncateAddress right-pads an address with zeros if it's shorter than `MinAddrLen` or
// takes the first `MaxAddrLen` if it's longer than that.
func PadOrTruncateAddress(address []byte) []byte {
	if len(address) > MaxAddrLen {
		return address[:MaxAddrLen]
	} else if len(address) < MinAddrLen {
		return append(address, make([]byte, MinAddrLen-len(address))...)
	}
	return address
}
