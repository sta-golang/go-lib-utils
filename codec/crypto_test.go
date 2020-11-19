package codec

import (
	"fmt"
	"testing"
)

func TestCryptoCodec_Check(t *testing.T) {
	kk := "d0dcbf0d12a6b1e7fbfa2ce5848f3eff"
	fmt.Println(API.CryptoAPI.Check("qq123456", kk))
}
