package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

type Transaction struct {
	TxId               string `json:"txid"`
	SourceAddress      string `json:"source_address"`
	DestinationAddress string `json:"destination_address"`
	Amount             int64  `json:"amount"`
	UnsignedTx         string `json:"unsignedtx"`
	SignedTx           string `json:"signedtx"`
}

func CreateTransaction(secret string, destination string, amount int64, txHash string) (Transaction, error) {
	var transaction Transaction

	wif, err := btcutil.DecodeWIF(secret)
	PanicErr(err)

	addressPubKey, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeUncompressed(), &chaincfg.MainNetParams)
	PanicErr(err)

	sourceTx := wire.NewMsgTx(wire.TxVersion)

	//transaction id for the redeem transaction
	sourceUTxOHash,  err := chainhash.NewHashFromStr(txHash)
	PanicErr(err)

	sourceUTxO := wire.NewOutPoint(sourceUTxOHash, 0)
	sourceTxIn := wire.NewTxIn(sourceUTxO, nil, nil)

	destinationAddress, err := btcutil.DecodeAddress(destination, &chaincfg.MainNetParams)
	PanicErr(err)

	sourceAddress, err := btcutil.DecodeAddress(addressPubKey.EncodeAddress(), &chaincfg.MainNetParams)
	PanicErr(err)

	destinationPkScript, err  := txscript.PayToAddrScript(destinationAddress)
	PanicErr(err)

	sourcePkScript, err := txscript.PayToAddrScript(sourceAddress)
	PanicErr(err)

	sourceTxOut := wire.NewTxOut(amount, sourcePkScript)

	sourceTx.AddTxIn(sourceTxIn)
	sourceTx.AddTxOut(sourceTxOut)

	sourceTxHash := sourceTx.TxHash()

	///Redeem transaction
	redeemTx := wire.NewMsgTx(wire.TxVersion)
	//model for previous transaction outputs
	prevOut := wire.NewOutPoint(&sourceTxHash, 0)

	//Add input
	redeemTxIn := wire.NewTxIn(prevOut, nil, nil)
	redeemTx.AddTxIn(redeemTxIn)

	//Add output
	redeemTxOut := wire.NewTxOut(amount, destinationPkScript)
	redeemTx.AddTxOut(redeemTxOut)

	//Sign the transaction
	sigScript, err := txscript.SignatureScript(redeemTx, 0, sourceTx.TxOut[0].PkScript, txscript.SigHashAll, wif.PrivKey, false)
	if err != nil {
		return Transaction{}, err
	}
	redeemTx.TxIn[0].SignatureScript = sigScript

	//Verify transaction
	flags := txscript.StandardVerifyFlags
	vm, err := txscript.NewEngine(sourceTx.TxOut[0].PkScript, redeemTx, 0, flags, nil, nil, amount)
	if err != nil {
		return Transaction{}, err
	}
	if err := vm.Execute(); err != nil {
		return Transaction{}, err
	}

	var unsignedTx bytes.Buffer
	var signedTx bytes.Buffer
	sourceTx.Serialize(&unsignedTx)
	redeemTx.Serialize(&signedTx)
	transaction.TxId = sourceTxHash.String()
	transaction.UnsignedTx = hex.EncodeToString(unsignedTx.Bytes())
	transaction.Amount = amount
	transaction.SignedTx = hex.EncodeToString(signedTx.Bytes())
	transaction.SourceAddress = sourceAddress.EncodeAddress()
	transaction.DestinationAddress = destinationAddress.EncodeAddress()
	return transaction, nil
}

func main() {
	transaction, err := CreateTransaction(
		"5HusYj2b2x4nroApgfvaSfKYZhRbKFH41bVyPooymbC6KfgSXdD",
		"1KKKK6N21XKo48zWKuQKXdvSsCf95ibHFa",
		91234,
		"81b4c832d70cb56ff957589752eb4125a4cab78a25a8fc52d6a09e5bd4404d48")
	PanicErr(err)

	data, _ := json.Marshal(transaction)
	fmt.Println(string(data))
}

func PanicErr(err error){
	if err != nil {
		panic(err)
	}
}