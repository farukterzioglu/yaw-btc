package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcec"

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


func getAddressFromWif(wif *btcutil.WIF) btcutil.Address{
	addressPubKey, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeUncompressed(), &chaincfg.MainNetParams)
	PanicErr(err)
	sourceAddress, err := btcutil.DecodeAddress(addressPubKey.EncodeAddress(), &chaincfg.MainNetParams)
	PanicErr(err)

	return sourceAddress
}
func createCoinBaseInput(txHash string) *wire.TxIn {
	sourceUTxOHash,  err := chainhash.NewHashFromStr(txHash)
	PanicErr(err)

	sourceUTxO := wire.NewOutPoint(sourceUTxOHash, 0)
	sourceTxIn := wire.NewTxIn(sourceUTxO, nil, nil)

	return sourceTxIn
}
func createCoinbaseOutput(sourceAddress btcutil.Address, amount int64) *wire.TxOut{
	sourcePkScript, err := txscript.PayToAddrScript(sourceAddress)
	PanicErr(err)

	sourceTxOut := wire.NewTxOut(amount, sourcePkScript)
	return sourceTxOut
}
func createMsgTx(input *wire.TxIn, output *wire.TxOut) *wire.MsgTx {
	sourceTx := wire.NewMsgTx(wire.TxVersion)
	sourceTx.AddTxIn(input)
	sourceTx.AddTxOut(output)
	return sourceTx
}
func createRedeemTransaction (destinationAddress btcutil.Address, amount int64, sourceTxHash chainhash.Hash, sourceOutputIndex uint32) *wire.MsgTx {
	redeemTx := wire.NewMsgTx(wire.TxVersion)

	//model for previous transaction outputs
	prevOut := wire.NewOutPoint(&sourceTxHash, sourceOutputIndex)

	//Add input
	redeemTxIn := wire.NewTxIn(prevOut, nil, nil)
	redeemTx.AddTxIn(redeemTxIn)

	//Add output
	destinationPkScript, err  := txscript.PayToAddrScript(destinationAddress)
	PanicErr(err)

	redeemTxOut := wire.NewTxOut(amount, destinationPkScript)
	redeemTx.AddTxOut(redeemTxOut)
	return redeemTx
}

func signTransaction(transaction *wire.MsgTx, inputIndex int, output *wire.TxOut, privateKey *btcec.PrivateKey){
	// SignatureScript creates an input signature script for tx to spend BTC sent
	// from a previous output to the owner of privKey. tx must include all
	// transaction inputs and outputs, however txin scripts are allowed to be filled
	// or empty. The returned script is calculated to be used as the idx'th txin
	// sigscript for tx. subscript is the PkScript of the previous output being used
	// as the idx'th input. privKey is serialized in either a compressed or
	// uncompressed format based on compress. This format must match the same format
	// used to generate the payment address, or the script validation will fail.

	sigScript, err 						:= txscript.SignatureScript(
		transaction,
		inputIndex,
		output.PkScript,
		txscript.SigHashAll,
		privateKey, false)
	PanicErr(err)

	transaction.TxIn[inputIndex].SignatureScript = sigScript
}
func CreateTransaction(secret string, destination string, amount int64, txHash string) (Transaction, error) {
	var wif *btcutil.WIF
	wif, _ 								= btcutil.DecodeWIF(secret)

	var sourceAddress btcutil.Address 	= getAddressFromWif(wif)

	//Coinbase transaction
	var coinBaseInput 	*wire.TxIn 		= createCoinBaseInput(txHash)
	var coinBaseOutput 	*wire.TxOut 	= createCoinbaseOutput(sourceAddress, amount)
	var sourceTx 		*wire.MsgTx 	= createMsgTx(coinBaseInput, coinBaseOutput)
	var sourceTxHash 	chainhash.Hash	= sourceTx.TxHash()

	///Redeem transaction
	destinationAddress, _ 				:= btcutil.DecodeAddress(destination, &chaincfg.MainNetParams)
	var redeemTx 		*wire.MsgTx 	= createRedeemTransaction(destinationAddress, amount, sourceTxHash, 0)

	//Sign the every input of the transaction
	signTransaction(redeemTx, 0, sourceTx.TxOut[0], wif.PrivKey)

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

	var transaction Transaction
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