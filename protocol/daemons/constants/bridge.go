package constants

const (
	// Signature of the Bridge log, i.e. sha256("Bridge(uint256,uint256,address,bytes,bytes)").
	BridgeEventSignature = "0x498a04382650bc110983392ed12ab27595af8ece270a344fc70d773d2481043a"

	// ABI (application binary interface) of the Bridge Event.
	BridgeEventABI = `[
		{
			"anonymous":false,
			"inputs":[
				{
					"indexed":true,
					"internalType":"uint256",
					"name":"id",
					"type":"uint256"
				},
				{
					"indexed":false,
					"internalType":"uint256",
					"name":"amount",
					"type":"uint256"
				},
				{
					"indexed":false,
					"internalType":"address",
					"name":"from",
					"type":"address"
				},
				{
					"indexed":false,
					"internalType":"bytes",
					"name":"accAddress",
					"type":"bytes"
				},
				{
					"indexed":false,
					"internalType":"bytes",
					"name":"data",
					"type":"bytes"
				}
			],
			"name":"Bridge",
			"type":"event"
		}
	]`
)
