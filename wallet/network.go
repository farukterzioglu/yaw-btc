package main

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
)

type Network struct {
	name        string
	symbol      string
	xpubkey     byte
	xprivatekey byte
	magic       wire.BitcoinNet
}

//Returns a Bitcoin network parameters.
func (network Network) GetNetworkParams() *chaincfg.Params {
	networkParams := &chaincfg.MainNetParams
	networkParams.Name = network.name
	networkParams.Net = network.magic
	networkParams.PubKeyHashAddrID = network.xpubkey
	networkParams.PrivateKeyID = network.xprivatekey
	return networkParams
}

var networks = map[string]Network{
	"btc":  {name: "bitcoin", symbol: "btc", xpubkey: 0x00, xprivatekey: 0x80, magic: 0xf9beb4d9},
}


