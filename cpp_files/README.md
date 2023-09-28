added folders amino , cosmos, cosmos_proto , gogproto,google for compiling protos to cpp files
you will find cpp files in every folder of dydxprotocol folder

example:
to find cpp files for subaccounts folder
cd  dydxprotocol/subaccounts/dydxprotocol/subaccounts


protoc command for ref:

protoc --proto_path=/Users/akshayhatkar/v4-chain/proto/ --cpp_out=. /Users/akshayhatkar/v4-chain/proto/dydxprotocol/bridge/*.proto
