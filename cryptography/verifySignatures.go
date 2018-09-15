package cryptography

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

func verify(){
	// Decode hex-encoded serialized public key.
	pubKeyBytes, err := hex.DecodeString("02a673638cb9587cb68ea08dbef685c" +
		"6f2d2a751a8b3c6f2a7e9a4999e6e4bfaf5")
	if err != nil {
		fmt.Println(err)
		return
	}
	pubKey, err := btcec.ParsePubKey(pubKeyBytes, btcec.S256())
	if err != nil {
		fmt.Println(err)
		return
	}

	// Decode hex-encoded serialized signature.
	sigBytes, err := hex.DecodeString("30450220090ebfb3690a0ff115bb1b38b" +
		"8b323a667b7653454f1bccb06d4bbdca42c2079022100ec95778b51e707" +
		"1cb1205f8bde9af6592fc978b0452dafe599481c46d6b2e479")

	if err != nil {
		fmt.Println(err)
		return
	}
	signature, err := btcec.ParseSignature(sigBytes, btcec.S256())
	if err != nil {
		fmt.Println(err)
		return
	}

	// Verify the signature for the message using the public key.
	message := "test message"
	messageHash := chainhash.DoubleHashB([]byte(message))
	verified := signature.Verify(messageHash, pubKey)
	fmt.Println("Signature Verified?", verified)
}
