package crypto

import (
	"github.com/larscom/bitvavo-go/internal/test"
	"testing"
)

func TestCreateSignature(t *testing.T) {
	body := "{\"market\":\"ETH-EUR\",\"amount\":1.5,\"price\":2500.5 }"
	timestamp := int64(1721452468484)
	sig := CreateSignature("GET", "/test", []byte(body), timestamp, "API_SECRET")
	test.AssertEqual(t, "d922f806412a560232d5326d95c389893432325f0e89f303f8ed5c9c04cc242b", sig)
}

func TestCreateSignatureNoBody(t *testing.T) {
	timestamp := int64(1721452468484)
	sig := CreateSignature("GET", "/test", nil, timestamp, "API_SECRET")
	test.AssertEqual(t, "dce7f0d49d559d6012733af234fa2bdef5a8492842726405e5c0b514f9bf1f55", sig)
}
