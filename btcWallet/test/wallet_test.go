package test

import "testing"
import (
	btcWallet "github.com/farukterzioglu/yew-btc/btcWallet"
	_ "fmt"
)

func TestWalletImport (t *testing.T){
	wallet := btcWallet.Wallet{ Coins : make([]btcWallet.Coin, 1)}
	if wallet.Coins == nil {
		t.Error("it is nil")
	}
}
