package btcWallet

import (
	"strings"
	"testing"
)

func TestWalletUseCases (t *testing.T){
	passPhrase := "mySecretKey"

	btcNetwork := networks[strings.ToLower("BTC")]

	//Create a wallet then encrypt it and store
	wallet := Wallet{ }
	wallet.Create(passPhrase)

	//Create & store new coin type in wallet
	coin := Coin{ Name: "Bitcoin", Symbol: "BTC"}
	//Creates private key, address and store inside Coin type
	coin.Generate(btcNetwork)

	//Open wallet
	wallet.DecryptFile(passPhrase)
	//Import generated coin type
	wallet.Import(coin, passPhrase)
	wif := coin.WIF

	//Import existing WIF to wallet
	coinToImport := Coin{ WIF: wif, Symbol: "BTC"}
	coinToImport.Import(btcNetwork)

	var existingWallet Wallet
	existingWallet.DecryptFile(passPhrase)
	wallet.Import(coinToImport, passPhrase)

	//Check if te passphrase is valid
	var walletToAuthenticate Wallet
	res := walletToAuthenticate.Authenticate(passPhrase)
	if !res {
		panic("Not authenticated!")
	}


}

