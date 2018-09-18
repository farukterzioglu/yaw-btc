package btcWallet

import (
	"errors"
	"os"
)

type Wallet struct {
	Coins []Coin `json:"coins"`
}

func (wallet Wallet) Create(key string) error {
	err := wallet.EncryptFile(key)
	return err
}

func (wallet Wallet) Destroy() error {
	err := os.Remove("wallet.dat")
	return err
}

func (wallet Wallet) Import(coin Coin, passphrase string) error {
	if coin == (Coin{}) {
		return errors.New("the coin must be valid")
	}
	wallet.Coins = append(wallet.Coins, coin)
	wallet.EncryptFile(passphrase)
	return nil
}

func (wallet *Wallet) Dump(passphrase string) error {
	err := wallet.DecryptFile(passphrase)
	if err != nil {
		return errors.New(`{ "message": "The password is not correct" }`)
	}
	return nil
}

func (wallet *Wallet) GetAddresses(passphrase string) error {
	err := wallet.DecryptFile(passphrase)
	if err != nil {
		return errors.New(`{ "message": "The password is not correct" }`)
	}
	for index := range wallet.Coins {
		wallet.Coins[index].WIF = ""
	}
	return nil
}

func (wallet *Wallet) Authenticate(passphrase string) bool {
	err := wallet.DecryptFile(passphrase)
	if err != nil {
		return false
	}
	return true
}
