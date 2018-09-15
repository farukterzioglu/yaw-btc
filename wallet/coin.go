package main

import (
	"errors"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

type Coin struct {
	Name                string `json:"name"`
	Symbol              string `json:"symbol"`
	WIF                 string `json:"wif,omitempty"`
	UncompressedAddress string `json:"uncompressed_address"`
	CompressedAddress   string `json:"compressed_address"`
}

//Generate a new set of keys for a network
func (coin *Coin) Generate(network Network) (err error) {
	if network == (Network{}) {
		return errors.New("unsupported cryptocurrency symbol provided")
	}

	//Generete new private key with secp256k1
	var privateKey *btcec.PrivateKey
	privateKey, err = btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return err
	}

	//Create wallet import format from private key
	wif, err := btcutil.NewWIF(privateKey, network.GetNetworkParams(), false)
	if err != nil {
		return err
	}
	coin.WIF = wif.String()

	var uncompressedAddress *btcutil.AddressPubKey
	uncompressedAddress, err = btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeUncompressed(), network.GetNetworkParams())
	if err != nil {
		return err
	}

	var compressedAddress *btcutil.AddressPubKey
	compressedAddress, err = btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), network.GetNetworkParams())
	if err != nil {
		return err
	}

	coin.UncompressedAddress = uncompressedAddress.EncodeAddress()
	coin.CompressedAddress = compressedAddress.EncodeAddress()
	coin.Name = network.name
	coin.Symbol = network.symbol

	return
}

func (coin *Coin) Import(network Network) error {
	if network == (Network{}) {
		return errors.New("unsupported cryptocurrency symbol provided")
	}
	wif, err := btcutil.DecodeWIF(coin.WIF)
	if err != nil {
		return err
	}
	if !wif.IsForNet(network.GetNetworkParams()) {
		return errors.New("The WIF string is not valid for the `" + network.name + "` network")
	}
	uncompressedAddress, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeUncompressed(), network.GetNetworkParams())
	if err != nil {
		return err
	}
	compressedAddress, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), network.GetNetworkParams())
	if err != nil {
		return err
	}
	coin.WIF = wif.String()
	coin.UncompressedAddress = uncompressedAddress.EncodeAddress()
	coin.CompressedAddress = compressedAddress.EncodeAddress()
	coin.Name = network.name
	coin.Symbol = network.symbol

	return nil
}
