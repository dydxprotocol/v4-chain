package constants

const (
	// Signature of the Bridge log, i.e. sha256("Bridge(uint256,uint256,bytes32,bytes)").
	BridgeEventSignature = "0xf8dd2841b36f876d311a264058cb076d68674181851a0688c405d2ae917a4fd2"

	// ABI (application binary interface) of the Bridge Event.
	BridgeEventABI = `[
		{
		  "anonymous": false,
		  "inputs": [
			{
			  "indexed": true,
			  "internalType": "uint256",
			  "name": "id",
			  "type": "uint256"
			},
			{
			  "indexed": false,
			  "internalType": "uint256",
			  "name": "amount",
			  "type": "uint256"
			},
			{
			  "indexed": false,
			  "internalType": "bytes32",
			  "name": "accAddress",
			  "type": "bytes32"
			},
			{
			  "indexed": false,
			  "internalType": "bytes",
			  "name": "data",
			  "type": "bytes"
			}
		  ],
		  "name": "Bridge",
		  "type": "event"
		}
	]`
)
