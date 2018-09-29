package networkConnection

import (
	"bytes"
	"encoding/hex"
	"testing"

	_ "github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	_ "github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

const txFee = 10000

type utxo struct {
	Address     string
	TxID        string
	OutputIndex uint32
	Script      string
	Satoshis    int64
	Height      int64
}

func TestWalletUseCases(t *testing.T) {
	unspentTx := utxo{
		Address:     "SZ8hKFdRsiTVedW4rzx8CWfK1wLAFjfQ4A",
		TxID:        "19e952a5d0b57d30a35465626ca6a88a7784e491eeb1b550da731b188bde8c69",
		OutputIndex: 0,
		Script:      "76a91481def9eeefa2f84fa436a9d75da6336ae8318c1d88ac",
		Satoshis:    5000000000,
	}

	tx := createTransactionTest(
		"FqBM82gKj2MXFurUjQaNVbbkXToU6k1fjvigSEMVf6WMDUd16jCj",
		"SdjQCwuYWpo6V2CSFVL1RnYsHanGoJh6ZQ",
		1,
		unspentTx)

	println(tx)
}

func createTransactionTest(secret string, destination string, amount int64, unspentTx utxo) string {
	simNetParams := &chaincfg.SimNetParams
	wif, err := btcutil.DecodeWIF(secret)

	destinationAddress, err := btcutil.DecodeAddress(destination, simNetParams)

	hash, err := chainhash.NewHashFromStr(unspentTx.TxID)

	//Redeem tx
	redeemTx := wire.NewMsgTx(wire.TxVersion)
	outPoint := wire.NewOutPoint(hash, unspentTx.OutputIndex)

	txIn := wire.NewTxIn(outPoint, nil, nil)
	redeemTx.AddTxIn(txIn)

	var destinationPkScript []byte
	destinationPkScript, err = txscript.PayToAddrScript(destinationAddress)

	// Pay the minimum network fee so that nodes will broadcast the tx.
	//TODO : ????
	outCoin := unspentTx.Satoshis - txFee

	redeemTxOut := wire.NewTxOut(outCoin, destinationPkScript)
	redeemTx.AddTxOut(redeemTxOut)

	addresspubkey, err := btcutil.NewAddressPubKey(
		wif.PrivKey.PubKey().SerializeUncompressed(),
		simNetParams)
	sourceAddress, err := btcutil.DecodeAddress(addresspubkey.EncodeAddress(), simNetParams)
	sourcePkScript, err := txscript.PayToAddrScript(sourceAddress)

	//Sign the transaction
	sigScript, err := txscript.SignatureScript(
		redeemTx, 0,
		sourcePkScript,
		txscript.SigHashAll,
		wif.PrivKey,
		false)

	redeemTx.TxIn[0].SignatureScript = sigScript

	flags := txscript.StandardVerifyFlags
	vm, err := txscript.NewEngine(sourcePkScript, redeemTx, 0, flags, nil, nil, amount)

	err = vm.Execute()
	if err != nil {
		panic("Not valid")
	}

	var signedTx bytes.Buffer
	redeemTx.Serialize(&signedTx)
	return hex.EncodeToString(signedTx.Bytes())
}
