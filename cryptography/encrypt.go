package cryptography

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
)

//https://godoc.org/github.com/btcsuite/btcd/btcec
func encrypt(){
	// Decode the hex-encoded pubkey of the recipient.
	pubKeyBytes, err := hex.DecodeString("04115c42e757b2efb7671c578530ec191a1" +
		"359381e6a71127a9d37c486fd30dae57e76dc58f693bd7e7010358ce6b165e483a29" +
		"21010db67ac11b1b51b651953d2") // uncompressed pubkey
	if err != nil {
		fmt.Println(err)
		return
	}
	pubKey, err := btcec.ParsePubKey(pubKeyBytes, btcec.S256())
	if err != nil {
		fmt.Println(err)
		return
	}

	// Encrypt a message decryptable by the private key corresponding to pubKey
	message := "test message"
	ciphertext, err := btcec.Encrypt(pubKey, []byte(message))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Decode the hex-encoded private key.
	pkBytes, err := hex.DecodeString("a11b0a4e1a132305652ee7a8eb7848f6ad" +
		"5ea381e3ce20a2c086a2e388230811")
	if err != nil {
		fmt.Println(err)
		return
	}
	// note that we already have corresponding pubKey
	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)

	// Try decrypting and verify if it's the same message.
	plaintext, err := btcec.Decrypt(privKey, ciphertext)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(plaintext))
}
