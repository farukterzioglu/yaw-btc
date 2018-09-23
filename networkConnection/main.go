package networkConnection

import (
	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"log"
)

type Transaction struct {
	TxId               string `json:"txid"`
	SourceAddress      string `json:"source_address"`
	DestinationAddress string `json:"destination_address"`
	Amount             int64  `json:"amount"`
	UnsignedTx         string `json:"unsignedtx"`
	SignedTx           string `json:"signedtx"`
}

func main() {
	tx := createTransaction(Transaction{})
	createBlock(tx)
	mineBlock()
}

func createTransaction(transaction Transaction) *wire.MsgTx {
	return wire.NewMsgTx(1)
}

func createBlock(*wire.MsgTx) {

}

func mineBlock() {
	createCoinbaseTx()
	createCoincaseTx()
}

func createCoinbaseTx() {

}

func createCoincaseTx() {

}

func panicIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
