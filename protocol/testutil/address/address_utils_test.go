package address

import (
	"fmt"
	"testing"
)

func TestAddressConversion(t *testing.T) {
	newAddr, err := ConvertAddressPrefix("dydx1x2hd82qerp7lc0kf5cs3yekftupkrl620te6u2", "klyra")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(newAddr)
}
