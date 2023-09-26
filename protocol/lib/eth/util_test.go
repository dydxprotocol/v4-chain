package eth_test

import (
	"encoding/json"
	"fmt"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"io/ioutil"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	libeth "github.com/dydxprotocol/v4-chain/protocol/lib/eth"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	ethcoretypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

func AddrToPort(addr string) string {
	parts := strings.Split(addr, ":")
	return parts[len(parts)-1]
}

func TestGroupBy2(t *testing.T) {
	nodes := strings.Split(
		"6f8f55881a894740b8bd2b362ac811237317ab32,4efb2ce99d6259d330627b3a124f550d5fb150b6,178b7abe7b6fbde8620588246ee7b63ed58feae1,58606527923333b4a7d685c3c3ad0027df8d2045,e4347bf1400d447fb0aeb980e93ee706988b2c7d,604c3f8abb235fa6b245484b01f3080e48b7456a,de1b4d964ecd006ec2d7f2eb12f5db0e67ff14cc,c3d51ae4cd2ef3537076cbbe204df259127a08e4,a7c50e19f82ad7b3180050393604d45676c793ca,835589aab7e1abeb8151e8f21b59e728d9b40f3b,dfa67970296bbecce14daba6cb0da516ed60458a,816b4e5cbaae22ffaf96106ca80b442873228a6c,3a68c220faf8dd5a64e4de01ca4fabc68d1a6acf,3ae00d90f0ba59f2ed8c3fd870b6f3b57be9246a,81b0f152201f839eab18e210340d5131cab8ec6d,8d011764ce872ec49c21b65403350a4700a3157a,2677b13ecdefddd54f36405e1e2bd1ec7478e8b2,67dc2eda6761aa9d18104896702e8b1e5c1dcb20,a1a4783fc9e9db9365b10317b5672c9343c5568c,f3a15e80d5bdb03a313fd86bff0215c048513824,0359d938a4c38adc766657c21dce54ba412281d7,a988534ab1e4bc42aad26ea7ec7bdc7d5415a14c,dcdbb554907adbcb3725fd5d547d5c7c3fa21250,ed4a9b9ae4f6e7aab14a7b28e4049de384dc5487,3f667030ddd9c561ec66f35e8221be0178cf62c4,6b40b34ec9d2ab7a491d2472a59bc26ebebeb442,bc1360f22de0c741bca96535ecdfc17362dc98d9,b2e6f21f76952fe5efb6f45bd732ef13bdbfa288,4dc9beaf9bc48dd7a164a1f9d2539caa44b3c377,15910c22a895fda3b125b6238c1a682db6435a62,44aa51186450156fae53595556f95fe74d0fd2e3,bdb7205533df19ec9bbbeb2e15cb57e64fe6afb2,afdf731f85291b6fe68e5ba0f18a60246908b6c1,befd2ecdafd23a4ad244e8db473ffc94c4307b35,4f0a61ee857009bf5cf71656ad02fdcc9cec42c5,d501756d3b2b9f1c7e23371ce1b3e4f23591d340,81147b86e209f8a1e9e82e7deabcaacbf6c37577,b94e2ef8c332c386e934882ac860c03e6c9979db,07237b2cc6baaaba09462993849cce9cc0c7a06f,a4d1a43a207a3f394625d49be9f73baad1329203,95c52cb8e00df207ea1f51b4bc081f057acf5b85,8dff35714decffdbe8e9e326a22bc0009b3a81cc,f4b35c56a77b6e389983e4e3902d0914af74b2e9,acd1db409f6bc98f23e20673a84a5bcca0b3c7da,7b075e8e37a97c23d001204e30bbd9995132906a,48d3c07b6a67dcec527f00fb8a5a5ebb2f490410,20b4f9207cdc9d0310399f848f057621f7251846,810f10402c42b6fc136f8f38652dde737965a569,c9a0b182ab5d10f758628d579b25d31846a2e549,5d842f0658a6cc7645f17061cdbf987a389202de,dc09ac8961cd8087fe0355818b1903b726c6a449,a26d6ea49d85f90aabff4698609f2beee9ebf03d,15a1c6920e2f65956f96126584d311b00e231d2f,064d1c8d8fee5baf7161db294ef9de9f91f66037,c4ed621f5b34c513d0ddd9420c0776b75d32bac2,19594e42a39d7cbea1dc38649e60951fe4bd8448,6bb0f9a20934b0ef503cbd7f07177f649f3726d3,1ae0bff4921df547efe6a01a441b584d041af8a3,19e3f6a648a6eb38f4bb789b27821a8bf02cb545,1ca78b330a3b00885f3aa0eb7849341afc12be6e,206a664050a7aa33f6b7229e47920bb9a7172392,47cb90109f169ac4a18154b16517d38c19830773,adc750b54aa31fef773f9cd40c210d32e0f361bd,aa0113770c1c906dd14dbc84a170186d2d412207,74a6c612c4abe9788693801d83ce19f4d2288af7,d99017279a3b1e05c2577b0d9bff4568f24ac5df,a21ca0f93f5735034fbeed8e76f1984e8be6ef57,69010d57a8227a5d6585bb5ea4337fce3507583e,9c506dcbd413e86031450befdf052ab44cf43dc7,cc3031901beda3da17c30b9bf9db2da8bbabf110,9472dde362189eb369e98246fe8c2ea5cb6521ab,dd3ed1acd8ab7ff404142256bcf89dc6d3ca1e41,6276dc0a4a5efa9b9ee328b0d03d15fdfe636035,150a3b8a6a25032bfcf43d9ea97d1f25f714e4d3,e9150ff614917ca2899a015eadb76a413f958b4e,26dd04963f0e8c56bde77d0a5414156d65a047f1,0eb076ed78e66e19f631b6f522e3cbf8e3f532a1,d812c7c0a5b4c4e9229c706241c0aa9a44e09378,320339c162630010dfcfde649cd9c51b0fd78550,47817889aece4dd89753173a3e4d1ee8daeb4898,938F8478449F1045861BB9B1130AFCE9A5875F42,4521ef33e881699ef1f97e2a8bbb0c9167a7673e,8d7353164ca242ebc2bde0c94e5e7dc8ad12491a",
		",")
	ips := strings.Split(
		"64.176.48.105,52.199.92.130,54.199.218.236,34.84.56.197,147.75.93.91,125.131.181.20,52.68.32.178,45.250.255.99,20.41.104.10,109.123.230.196,54.64.201.139,65.108.13.110,65.108.105.48,65.109.53.22,65.109.23.114,34.88.180.144,65.108.226.183,135.181.216.54,65.109.80.150,65.109.97.62,65.109.112.162,172.111.52.52,190.124.251.29,3.139.127.183,3.22.142.107,34.75.252.225,34.152.16.59,34.174.247.174,54.39.16.240,160.202.129.79,167.71.248.97,35.247.200.116,35.230.117.136,93.115.25.18,34.116.141.89,146.59.71.233,195.14.6.186,195.14.6.184,195.14.6.185,165.232.88.32,88.198.14.106,185.119.116.239,74.118.136.153,146.59.69.126,141.164.51.182,34.80.18.27,158.247.233.222,34.81.53.149,211.219.19.81,170.64.216.189,35.213.143.229,64.120.114.5,34.92.198.135,13.209.114.32,15.165.189.236,13.229.142.253,178.23.126.123,51.195.61.9,62.210.145.130,51.210.146.190,46.4.39.178,85.207.33.82,88.99.143.105,5.9.100.25,46.4.81.211,165.22.73.158,173.212.223.233,51.79.229.127,139.180.141.248,54.178.51.9,18.138.11.164,103.106.228.57,[2400:d320:2146:5223::2],13.114.96.36,188.42.220.52,194.233.95.115,52.74.20.233,104.156.238.56,149.28.17.195,45.32.36.82,35.78.229.77,13.213.5.36,15.235.204.15",
		",")
	ports := strings.Split(
		"26656,26656,26656,26656,26656,26656,26656,26656,26656,26656,26656,26656,23856,55836,23856,26656,23856,3090,2690,26656,26656,32667,26670,26656,26656,26656,26656,26656,26656,26656,26180,26656,26656,54656,26656,26656,26656,26656,26656,26180,26656,26656,26656,32602,26656,26656,26656,26656,26656,26180,26656,23856,26656,26656,26656,26656,31310,26856,26656,26656,26656,26656,26656,26656,26656,26656,10456,26656,55836,26656,26656,26656,26656,26656,26656,26657,26656,26656,26656,26656,26656,26656,26656",
		",")
	clusters := strings.Split(
		"0,0,0,0,0,0,0,0,0,0,0,1,1,1,1,1,1,1,1,1,1,2,2,2,2,2,2,2,2,2,2,2,2,3,3,3,3,3,3,3,3,3,3,3,4,4,4,4,4,4,4,4,4,4,4,4,5,5,5,5,5,5,5,5,5,5,5,6,6,6,6,6,6,6,6,6,6,6,6,0,0,4,4",
		",")
	ipToNode := make(map[string]string)
	clusterToIps := make(map[string][]string)
	ipToPort := make(map[string]string)

	for i, ip := range ips {
		ipToNode[ip] = nodes[i]
		ipToPort[ip] = ports[i]
		clusterToIps[clusters[i]] = append(clusterToIps[clusters[i]], ip)
	}

	for i := range ips {
		foo := make([]string, 0)
		for _, key := range []string{"0", "1", "2", "3", "4", "5", "6"} {
			if key != clusters[i] {
				value := clusterToIps[key]
				ip := value[i%len(value)]
				foo = append(foo, fmt.Sprintf("%s@%s:%s", ipToNode[ip], ip, ipToPort[ip]))
			}
		}
		for _, ip := range clusterToIps[clusters[i]] {
			if ips[i] != ip {
				foo = append(foo, fmt.Sprintf("%s@%s:%s", ipToNode[ip], ip, ipToPort[ip]))
			}
		}
		println(strings.Join(foo, ","))
	}

}

func TestGroupBy(t *testing.T) {
	file, _ := ioutil.ReadFile("foo.json")
	nodes := make([]coretypes.Peer, 0)
	_ = json.Unmarshal(file, &nodes)

	//ports := strings.Split("//0.0.0.0:26656,51.195.61.9:26856,//0.0.0.0:26656,195.14.6.186:26656,//0.0.0.0:26656,//0.0.0.0:26656,172.111.52.52:32667,//0.0.0.0:26656,//0.0.0.0:26656,//93.115.25.18:54656,//18.182.95.191:26656,65.108.13.110:26656,165.232.88.32:26180,//0.0.0.0:26656,195.14.6.184:26656,//0.0.0.0:26656,//0.0.0.0:26656,65.108.105.48:23856,65.109.53.22:55836,170.64.216.189:26180,//0.0.0.0:26656,65.109.23.114:23856,190.124.251.29:26670,//0.0.0.0:26656,//0.0.0.0:26656,//0.0.0.0:26656,85.207.33.82:26656,//3.139.127.183:26656,//0.0.0.0:26656,//0.0.0.0:26656,//52.68.32.178:26656,//3.20.153.106:26656,//0.0.0.0:26656,65.108.226.183:23856,//0.0.0.0:26656,//0.0.0.0:26656,//0.0.0.0:26656,146.59.71.233:26656,5.9.100.25:26656,//0.0.0.0:26656,//54.39.16.240:26656,//0.0.0.0:3090,185.119.116.239:26656,//0.0.0.0:26656,//0.0.0.0:2690,65.109.97.62:26656,64.120.114.5:23856,195.14.6.185:26656,167.71.248.97:26180,//0.0.0.0:26656,//0.0.0.0:26656,//0.0.0.0:26656,//0.0.0.0:26656,//0.0.0.0:26656,//0.0.0.0:26656,20.41.104.10:26656,//0.0.0.0:26656,//0.0.0.0:26656,109.123.230.196:26656,//13.230.189.28:26656,//0.0.0.0:26656", ",")
	//for i, port := range ports {
	//	ports[i] = strings.Split(port, ":")[1]
	//}
	//fmt.Printf("%+v\n", ports)
	//nodes := strings.Split("6f8f55881a894740b8bd2b362ac811237317ab32,4efb2ce99d6259d330627b3a124f550d5fb150b6,178b7abe7b6fbde8620588246ee7b63ed58feae1,58606527923333b4a7d685c3c3ad0027df8d2045,e4347bf1400d447fb0aeb980e93ee706988b2c7d,604c3f8abb235fa6b245484b01f3080e48b7456a,de1b4d964ecd006ec2d7f2eb12f5db0e67ff14cc,c3d51ae4cd2ef3537076cbbe204df259127a08e4,a7c50e19f82ad7b3180050393604d45676c793ca,835589aab7e1abeb8151e8f21b59e728d9b40f3b,dfa67970296bbecce14daba6cb0da516ed60458a,816b4e5cbaae22ffaf96106ca80b442873228a6c,3a68c220faf8dd5a64e4de01ca4fabc68d1a6acf,3ae00d90f0ba59f2ed8c3fd870b6f3b57be9246a,81b0f152201f839eab18e210340d5131cab8ec6d,8d011764ce872ec49c21b65403350a4700a3157a,2677b13ecdefddd54f36405e1e2bd1ec7478e8b2,67dc2eda6761aa9d18104896702e8b1e5c1dcb20,a1a4783fc9e9db9365b10317b5672c9343c5568c,f3a15e80d5bdb03a313fd86bff0215c048513824,a988534ab1e4bc42aad26ea7ec7bdc7d5415a14c,dcdbb554907adbcb3725fd5d547d5c7c3fa21250,ed4a9b9ae4f6e7aab14a7b28e4049de384dc5487,3f667030ddd9c561ec66f35e8221be0178cf62c4,6b40b34ec9d2ab7a491d2472a59bc26ebebeb442,bc1360f22de0c741bca96535ecdfc17362dc98d9,b2e6f21f76952fe5efb6f45bd732ef13bdbfa288,4dc9beaf9bc48dd7a164a1f9d2539caa44b3c377,15910c22a895fda3b125b6238c1a682db6435a62,44aa51186450156fae53595556f95fe74d0fd2e3,bdb7205533df19ec9bbbeb2e15cb57e64fe6afb2,afdf731f85291b6fe68e5ba0f18a60246908b6c1,befd2ecdafd23a4ad244e8db473ffc94c4307b35,4f0a61ee857009bf5cf71656ad02fdcc9cec42c5,d501756d3b2b9f1c7e23371ce1b3e4f23591d340,81147b86e209f8a1e9e82e7deabcaacbf6c37577,b94e2ef8c332c386e934882ac860c03e6c9979db,07237b2cc6baaaba09462993849cce9cc0c7a06f,a4d1a43a207a3f394625d49be9f73baad1329203,95c52cb8e00df207ea1f51b4bc081f057acf5b85,8dff35714decffdbe8e9e326a22bc0009b3a81cc,f4b35c56a77b6e389983e4e3902d0914af74b2e9,7b075e8e37a97c23d001204e30bbd9995132906a,48d3c07b6a67dcec527f00fb8a5a5ebb2f490410,20b4f9207cdc9d0310399f848f057621f7251846,810f10402c42b6fc136f8f38652dde737965a569,c9a0b182ab5d10f758628d579b25d31846a2e549,5d842f0658a6cc7645f17061cdbf987a389202de,dc09ac8961cd8087fe0355818b1903b726c6a449,a26d6ea49d85f90aabff4698609f2beee9ebf03d,6bb0f9a20934b0ef503cbd7f07177f649f3726d3,15a1c6920e2f65956f96126584d311b00e231d2f,1ae0bff4921df547efe6a01a441b584d041af8a3,19e3f6a648a6eb38f4bb789b27821a8bf02cb545,1ca78b330a3b00885f3aa0eb7849341afc12be6e,206a664050a7aa33f6b7229e47920bb9a7172392,47cb90109f169ac4a18154b16517d38c19830773,adc750b54aa31fef773f9cd40c210d32e0f361bd,aa0113770c1c906dd14dbc84a170186d2d412207,74a6c612c4abe9788693801d83ce19f4d2288af7,d99017279a3b1e05c2577b0d9bff4568f24ac5df", ",")
	ips := strings.Split("64.176.48.105,52.199.92.130,54.199.218.236,34.84.56.197,147.75.93.91,125.131.181.20,52.68.32.178,45.250.255.99,20.41.104.10,109.123.230.196,54.64.201.139,65.108.13.110,65.108.105.48,65.109.53.22,65.109.23.114,34.88.180.144,65.108.226.183,135.181.216.54,65.109.80.150,65.109.97.62,172.111.52.52,190.124.251.29,3.139.127.183,3.22.142.107,34.75.252.225,34.152.16.59,34.174.247.174,54.39.16.240,160.202.129.79,167.71.248.97,35.247.200.116,35.230.117.136,93.115.25.18,34.116.141.89,146.59.71.233,195.14.6.186,195.14.6.184,195.14.6.185,165.232.88.32,88.198.14.106,185.119.116.239,74.118.136.153,141.164.51.182,34.80.18.27,158.247.233.222,34.81.53.149,211.219.19.81,170.64.216.189,35.213.143.229,64.120.114.5,34.93.129.140,34.92.198.135,51.195.61.9,62.210.145.130,51.210.146.190,46.4.39.178,85.207.33.82,88.99.143.105,5.9.100.25,46.4.81.211,165.22.73.158", ",")
	clusters := strings.Split("0,0,0,0,0,0,0,0,0,0,0,1,1,1,1,1,1,1,1,1,2,2,2,2,2,2,2,2,2,2,2,2,3,3,3,3,3,3,3,3,3,3,4,4,4,4,4,4,4,4,4,4,5,5,5,5,5,5,5,5,5", ",")

	clusterToIps := make(map[string][]string)
	ipToNode := make(map[string]coretypes.Peer)

	for _, node := range nodes {
		ipToNode[node.RemoteIP] = node
	}
	for _, cluster := range clusters {
		clusterToIps[cluster] = make([]string, 0)
	}
	for i, ip := range ips {
		node := ipToNode[ip]
		fmt.Printf("%s,%s,%s,%s\n", node.NodeInfo.DefaultNodeID, node.RemoteIP, AddrToPort(node.NodeInfo.ListenAddr), node.NodeInfo.Moniker)
		clusterToIps[clusters[i]] = append(clusterToIps[clusters[i]], ip)
	}
	print("\n\n\n")
	for i := range ips {
		foo := make([]string, 0)
		for _, key := range []string{"0", "1", "2", "3", "4", "5"} {
			if key != clusters[i] {
				value := clusterToIps[key]
				ip := value[i%len(value)]
				foo = append(foo, fmt.Sprintf("%s@%s:%s", ipToNode[ip].NodeInfo.DefaultNodeID, ip, AddrToPort(ipToNode[ip].NodeInfo.ListenAddr)))
			}
		}
		for _, ip := range clusterToIps[clusters[i]] {
			if ips[i] != ip {
				foo = append(foo, fmt.Sprintf("%s@%s:%s", ipToNode[ip].NodeInfo.DefaultNodeID, ip, AddrToPort(ipToNode[ip].NodeInfo.ListenAddr)))
			}
		}
		println(strings.Join(foo, ","))
	}
}

func TestBridgeLogToEvent(t *testing.T) {
	tests := map[string]struct {
		inputLog   ethcoretypes.Log
		inputDenom string

		expectedEvent bridgetypes.BridgeEvent
	}{
		"Success: event ID 0": {
			inputLog:   constants.EthLog_Event0,
			inputDenom: "dv4tnt",
			expectedEvent: bridgetypes.BridgeEvent{
				Id: 0,
				Coin: sdk.NewCoin(
					"dv4tnt",
					sdk.NewInt(12345),
				),
				Address:        "dydx1qqgzqvzq2ps8pqys5zcvp58q7rluextx92xhln",
				EthBlockHeight: 3872013,
			},
		},
		"Success: event ID 1 - empty address": {
			inputLog:   constants.EthLog_Event1,
			inputDenom: "test-token",
			expectedEvent: bridgetypes.BridgeEvent{
				Id: 1,
				Coin: sdk.NewCoin(
					"test-token",
					sdk.NewInt(55),
				),
				// address shorter than 20 bytes is padded with zeros.
				Address:        "dydx1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq66wm82",
				EthBlockHeight: 3969937,
			},
		},
		"Success: event ID 2": {
			inputLog:   constants.EthLog_Event2,
			inputDenom: "test-token",
			expectedEvent: bridgetypes.BridgeEvent{
				Id: 2,
				Coin: sdk.NewCoin(
					"test-token",
					sdk.NewInt(777),
				),
				// 32 bytes * 8 bits / 5 bits = 51.2 characters ~ 52 bech32 characters
				Address:        "dydx1qqgzqvzq2ps8pqys5zcvp58q7rluextxzy3rx3z4vemc3xgq42as94fpcv",
				EthBlockHeight: 4139345,
			},
		},
		"Success: event ID 3": {
			inputLog:   constants.EthLog_Event3,
			inputDenom: "test-token-2",
			expectedEvent: bridgetypes.BridgeEvent{
				Id: 3,
				Coin: sdk.NewCoin(
					"test-token-2",
					sdk.NewInt(888),
				),
				// address data is 62 bytes but we take the first 32 bytes only.
				// 32 bytes * 8 bits / 5 bits ~ 52 bech32 characters
				Address:        "dydx124n92ej4ve2kv4tx24n92ej4ve2kv4tx24n92ej4ve2kv4tx24nq8exmjh",
				EthBlockHeight: 4139348,
			},
		},
		"Success: event ID 4": {
			inputLog:   constants.EthLog_Event4,
			inputDenom: "dv4tnt",
			expectedEvent: bridgetypes.BridgeEvent{
				Id: 4,
				Coin: sdk.NewCoin(
					"dv4tnt",
					sdk.NewInt(1234123443214321),
				),
				// address shorter than 20 bytes is padded with zeros.
				Address:        "dydx1zg6pydqqqqqqqqqqqqqqqqqqqqqqqqqqm0r5ra",
				EthBlockHeight: 4139349,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			event := libeth.BridgeLogToEvent(tc.inputLog, tc.inputDenom)
			require.Equal(t, tc.expectedEvent, event)
		})
	}
}
