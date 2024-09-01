package crypto

import (
	"github.com/larscom/bitvavo-go/v2/internal/test"
	"testing"
)

func TestCreateSignature(t *testing.T) {
	body := "{\"market\":\"ETH-EUR\",\"amount\":1.5,\"price\":2500.5}"
	timestamp := int64(1721452468484)
	sig := CreateSignature("GET", "/test", []byte(body), timestamp, "API_SECRET")
	test.AssertEqual(t, "cf9f81048eccf714305dfd0147252a38de6788ec343f4466a124ffe7c524ded8", sig)
}

func TestCreateSignatureNoBody(t *testing.T) {
	timestamp := int64(1721452468484)
	sig := CreateSignature("GET", "/test", nil, timestamp, "API_SECRET")
	test.AssertEqual(t, "dce7f0d49d559d6012733af234fa2bdef5a8492842726405e5c0b514f9bf1f55", sig)
}
