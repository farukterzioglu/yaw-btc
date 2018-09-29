btcctl -C .\btcctl.conf getblock 3b8a87dc55b02e15454789f4fadc9c4cecdbfc854b4d6073df704afc8cbe3026  
{
  "hash": "3b8a87dc55b02e15454789f4fadc9c4cecdbfc854b4d6073df704afc8cbe3026",
  "confirmations": 1306,
  "strippedsize": 188,
  "size": 188,
  "weight": 752,
  "height": 1,
  "version": 536870912,
  "versionHex": "20000000",
  "merkleroot": "19e952a5d0b57d30a35465626ca6a88a7784e491eeb1b550da731b188bde8c69",
  "tx": [
    "19e952a5d0b57d30a35465626ca6a88a7784e491eeb1b550da731b188bde8c69"
  ],
  "time": 1538073661,
  "nonce": 0,
  "bits": "207fffff",
  "difficulty": 1,
  "previousblockhash": "683e86bd5c6d110d91b94b97137ba6bfe02dbbdb8e3dff722a669b5d69d77af6",
  "nextblockhash": "580ce2f916598487b8f9b74db6c4f75479a177347cb53670f6cd1c4891e0b69e"
}  

btcctl -C .\btcctl.conf --wallet gettransaction 19e952a5d0b57d30a35465626ca6a88a7784e491eeb1b550da731b188bde8c69  
{
  "amount": 50,
  "confirmations": 1306,
  "blockhash": "3b8a87dc55b02e15454789f4fadc9c4cecdbfc854b4d6073df704afc8cbe3026",
  "blockindex": 0,
  "blocktime": 1538073661,
  "txid": "19e952a5d0b57d30a35465626ca6a88a7784e491eeb1b550da731b188bde8c69",
  "walletconflicts": [],
  "time": 1538073661,
  "timereceived": 1538073661,
  "details": [
    {
      "account": "default",
      "address": "SZ8hKFdRsiTVedW4rzx8CWfK1wLAFjfQ4A",
      "amount": 50,
      "category": "generate",
      "vout": 0
    }
  ],
  "hex": "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff165108cb8279c4c50b27a40b2f503253482f627463642fffffffff0100f2052a010000001976a91481def9eeefa2f84fa436a9d75da6336ae8318c1d88ac00000000"
}  

btcctl -C .\btcctl.conf --wallet decoderawtransaction 01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff165108cb8279c4c50b27a40b2f503253482f627463642fffffffff0100f2052a010000001976a91481def9eeefa2f84fa436a9d75da6336ae8318c1d88ac00000000  
{
  "txid": "19e952a5d0b57d30a35465626ca6a88a7784e491eeb1b550da731b188bde8c69",
  "version": 1,
  "locktime": 0,
  "vin": [
    {
      "coinbase": "5108cb8279c4c50b27a40b2f503253482f627463642f",
      "sequence": 4294967295
    }
  ],
  "vout": [
    {
      "value": 50,
      "n": 0,
      "scriptPubKey": {
        "asm": "OP_DUP OP_HASH160 81def9eeefa2f84fa436a9d75da6336ae8318c1d OP_EQUALVERIFY OP_CHECKSIG",
        "hex": "76a91481def9eeefa2f84fa436a9d75da6336ae8318c1d88ac",
        "reqSigs": 1,
        "type": "pubkeyhash",
        "addresses": [
          "SZ8hKFdRsiTVedW4rzx8CWfK1wLAFjfQ4A"
        ]
      }
    }
  ]
}  

//Created a transaction in main_test.go with informations above and print the hex of it (below)

btcctl -C .\btcctl.conf --wallet sendrawtransaction 0100000001698cde8b181b73da50b5b1ee91e484778aa8a66c626554a3307db5d0a552e919000000008a47304402207e07eb658b28ae7484dcb05c910f28447ba5ec4221f9f9406a366869c257e8f7022063b1ad97faadf59d1b6a7dc5cf5d26f194041cfd8bb5e58b1e7066c7995acab90141044abe53b41d5923945575c705b99c0806d54b4735b1a7f0c00e030df7121828552459d171445835ea145ebe3cc30f51922a801e67fd7216a0eb05d8cd914989b4ffffffff01f0ca052a010000001976a914b44fc1a78817026ce866baf41b45efe25b848d6688ac00000000  
