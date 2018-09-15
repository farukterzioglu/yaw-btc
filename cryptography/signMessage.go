package cryptography

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

//https://godoc.org/github.com/btcsuite/btcd/btcec#example-package--SignMessage
func signMessage(){
	// Decode a hex-encoded private key.
	pkBytes, err := hex.DecodeString("22a47fa09a223f2aa079edf85a7c2d4f87" +
		"20ee63e502ee2869afab7de234b80c")
	if err != nil {
		fmt.Println(err)
		return
	}
	privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)

	// Sign a message using the private key.
	message := "test message"
	messageHash := chainhash.DoubleHashB([]byte(message))
	signature, err := privKey.Sign(messageHash)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Serialize and display the signature.
	fmt.Printf("Serialized Signature: %x\n", signature.Serialize())

	// Verify the signature for the message using the public key.
	verified := signature.Verify(messageHash, pubKey)
	fmt.Printf("Signature Verified? %v\n", verified)
}
