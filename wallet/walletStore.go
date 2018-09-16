package btcWallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

func (wallet Wallet) EncryptFile(passphrase string) error {
	file, err := os.Create("wallet.dat")
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := wallet.encrypt(passphrase)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	return err
}
func (wallet *Wallet) DecryptFile(passphrase string) error {
	ciphertext, err := ioutil.ReadFile("wallet.dat")
	if err != nil {
		return err
	}
	err = wallet.decrypt(ciphertext, passphrase)
	if err != nil {
		return err
	}
	return nil
}

func (wallet Wallet) createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
func (wallet Wallet) encrypt(passphrase string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(wallet.createHash(passphrase)))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	data, err := json.Marshal(wallet)
	if err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}
func (wallet *Wallet) decrypt(data []byte, passphrase string) error {
	block, err := aes.NewCipher([]byte(wallet.createHash(passphrase)))
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}
	json.Unmarshal(plaintext, wallet)
	return nil
}

