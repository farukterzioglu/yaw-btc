package btcWallet

import (
	"fmt"
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

	//Get addresses
	walletOfAddresses := walletToAuthenticate
	walletOfAddresses.GetAddresses(passPhrase)
	for _, address := range walletOfAddresses.Coins{
		fmt.Printf("%s : %s \n", address.Symbol, address.UncompressedAddress)
	}
}

func TestTransaction(T *testing.T){
	passPhrase := "mySecretKey"
	btcNetwork := networks[strings.ToLower("BTC")]

	sender := Coin{ Name: "Bitcoin", Symbol: "BTC"}
	sender.Generate(btcNetwork)

	receiver := Coin{ Name: "Bitcoin", Symbol: "BTC"}
	receiver.Generate(btcNetwork)

	wallet := Wallet{ }
	wallet.Create(passPhrase)
	wallet.Import(sender, passPhrase)
	wallet.Import(receiver, passPhrase)

	var walletWithDetails Wallet
	walletWithDetails.Dump(passPhrase)

	//TODO : Get transaction from endpoint
	transaction := Transaction{
		Amount:10,
		DestinationAddress:sender.UncompressedAddress,
		SourceAddress:receiver.UncompressedAddress,
		TxId:"",
	}

	//Select btc
	var coinToBeTransfered *Coin
	for _, coin := range walletWithDetails.Coins {
		if coin.UncompressedAddress == transaction.SourceAddress {
			coinToBeTransfered = &coin
		}
	}

	if coinToBeTransfered == nil {
		return
	}

	//Check if coin type exist in wallet
	if transaction.SourceAddress == coinToBeTransfered.UncompressedAddress {
		tx, err := CreateTransaction(
			networks[strings.ToLower(coinToBeTransfered.Symbol)],
			coinToBeTransfered.WIF,
			transaction.DestinationAddress,
			transaction.Amount,
			transaction.TxId)

		if err != nil {
			panic(err)
		}

		fmt.Printf("Transaction %s", tx.TxId)
	}
}
