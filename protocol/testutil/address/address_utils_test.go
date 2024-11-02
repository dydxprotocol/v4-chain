package address

import (
	"fmt"
	"testing"
)

func TestAddressConversion(t *testing.T) {
	newAddr, err := ConvertAddressPrefix("klyravaloper10fx7sy6ywd5senxae9dwytf8jxek3t2gytemk4", "klyravalcons")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(newAddr)
}
