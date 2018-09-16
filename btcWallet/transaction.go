package btcWallet


import (
	"bytes"
	"encoding/hex"

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

func CreateTransaction(
	network Network,
	secret string,
	destination string,
	amount int64,
	txHash string) (transaction Transaction, err error) {
	transaction = Transaction{}

	wif, err := btcutil.DecodeWIF(secret)
	if err != nil { return }

	//A new AddressPubKey which represents a pay-to-pubkey address
	var addresspubkey *btcutil.AddressPubKey
	addresspubkey, err = btcutil.NewAddressPubKey(
		wif.PrivKey.PubKey().SerializeUncompressed(),
		network.GetNetworkParams())
	if err != nil { return }

	var sourceUtxOHash *chainhash.Hash
	sourceUtxOHash, err = chainhash.NewHashFromStr(txHash)
	if err != nil { return }

	//OutPoint defines a bitcoin data type that is used to track previous transaction outputs
	var sourceUtxo *wire.OutPoint 	= wire.NewOutPoint(sourceUtxOHash, 0)
	//Create input by using pointed output above
	var sourceTxIn *wire.TxIn 		= wire.NewTxIn(sourceUtxo, nil, nil)

	var destinationAddress 	btcutil.Address
	var sourceAddress 		btcutil.Address
	destinationAddress, err	= btcutil.DecodeAddress(destination, network.GetNetworkParams())
	sourceAddress, err 		= btcutil.DecodeAddress(addresspubkey.EncodeAddress(), network.GetNetworkParams())
	if err != nil { return }

	var sourcePkScript 		[]byte
	sourcePkScript, err 		= txscript.PayToAddrScript(sourceAddress)
	if err != nil { return }

	var sourceTxOut *wire.TxOut = wire.NewTxOut(amount, sourcePkScript)

	//Create source tx
	var sourceTx *wire.MsgTx 	= wire.NewMsgTx(wire.TxVersion)
	sourceTx.AddTxIn(sourceTxIn)
	sourceTx.AddTxOut(sourceTxOut)

	var sourceTxHash chainhash.Hash = sourceTx.TxHash()

	//New transaction to send
	redeemTx := wire.NewMsgTx(wire.TxVersion)
	//Create outpoint from senders output. "0"th output from tx "sourceTxHash"
	destinationUtxo := wire.NewOutPoint(&sourceTxHash, 0)

	var redeemTxIn *wire.TxIn 		= wire.NewTxIn(destinationUtxo, nil, nil)
	redeemTx.AddTxIn(redeemTxIn)

	//Receiver address
	var destinationPkScript	[]byte
	destinationPkScript, err	= txscript.PayToAddrScript(destinationAddress)
	if err != nil { return }

	//Create output
	var redeemTxOut *wire.TxOut		= wire.NewTxOut(amount, destinationPkScript)
	redeemTx.AddTxOut(redeemTxOut)

	//Sign the transaction
	sigScript, err := txscript.SignatureScript(
		redeemTx, 0,
		sourceTx.TxOut[0].PkScript, txscript.SigHashAll, wif.PrivKey, false)
	if err != nil { return }

	redeemTx.TxIn[0].SignatureScript = sigScript

	//validate
	flags := txscript.StandardVerifyFlags
	vm, err := txscript.NewEngine(sourceTx.TxOut[0].PkScript, redeemTx, 0, flags, nil, nil, amount)
	if err != nil { return }

	err = vm.Execute()
	if err != nil { return }

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