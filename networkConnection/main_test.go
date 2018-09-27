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

func TestWalletUseCases(t *testing.T) {
	tx := createTransactionTest(
		"FqBM82gKj2MXFurUjQaNVbbkXToU6k1fjvigSEMVf6WMDUd16jCj",
		"SdjQCwuYWpo6V2CSFVL1RnYsHanGoJh6ZQ",
		1,
		"c9705822b3650b9c2c34980770ea8b6f7b6297909281c4b1871edc14584897a9")

	println(tx)
}

func createTransactionTest(secret string, destination string, amount int64, txHash string) string {
	simNetParams := &chaincfg.SimNetParams

	sourceUtxOHash, err := chainhash.NewHashFromStr(txHash)
	var sourceUtxo *wire.OutPoint = wire.NewOutPoint(sourceUtxOHash, 0)
	//Create input by using pointed output above
	var sourceTxIn *wire.TxIn = wire.NewTxIn(sourceUtxo, nil, nil)

	wif, err := btcutil.DecodeWIF(secret)
	addresspubkey, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeUncompressed(), simNetParams)

	sourceAddress, err := btcutil.DecodeAddress(addresspubkey.EncodeAddress(), simNetParams)

	sourcePkScript, err := txscript.PayToAddrScript(sourceAddress)
	sourceTxOut := wire.NewTxOut(amount, sourcePkScript)

	var sourceTx *wire.MsgTx = wire.NewMsgTx(wire.TxVersion)
	sourceTx.AddTxIn(sourceTxIn)
	sourceTx.AddTxOut(sourceTxOut)
	sourceTxHash := sourceTx.TxHash()

	destinationAddress, err := btcutil.DecodeAddress(destination, simNetParams)

	//Redeem tx
	redeemTx := wire.NewMsgTx(wire.TxVersion)
	destinationUtxo := wire.NewOutPoint(&sourceTxHash, 0)

	var redeemTxIn *wire.TxIn = wire.NewTxIn(destinationUtxo, nil, nil)
	redeemTx.AddTxIn(redeemTxIn)

	var destinationPkScript []byte
	destinationPkScript, err = txscript.PayToAddrScript(destinationAddress)

	redeemTxOut := wire.NewTxOut(amount, destinationPkScript)
	redeemTx.AddTxOut(redeemTxOut)

	//Sign the transaction
	sigScript, err := txscript.SignatureScript(
		redeemTx, 0,
		sourceTx.TxOut[0].PkScript, txscript.SigHashAll, wif.PrivKey, false)

	redeemTx.TxIn[0].SignatureScript = sigScript

	flags := txscript.StandardVerifyFlags
	vm, err := txscript.NewEngine(sourceTx.TxOut[0].PkScript, redeemTx, 0, flags, nil, nil, amount)

	err = vm.Execute()
	if err != nil {
		panic("Not valid")
	}

	var signedTx bytes.Buffer
	redeemTx.Serialize(&signedTx)
	return hex.EncodeToString(signedTx.Bytes())
}
