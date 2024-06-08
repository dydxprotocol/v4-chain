export const external = `import { QueryClient, createProtobufRpcClient, ProtobufRpcClient } from '@cosmjs/stargate'
import { Tendermint34Client, HttpEndpoint } from "@cosmjs/tendermint-rpc";

const _rpcClients: Record<string, ProtobufRpcClient> = {};

export const getRpcEndpointKey = (rpcEndpoint: string | HttpEndpoint) => {
    if (typeof rpcEndpoint === 'string') {
        return rpcEndpoint;
    } else if (!!rpcEndpoint) {
        //@ts-ignore 
        return rpcEndpoint.url;
    }
}

export const getRpcClient = async (rpcEndpoint: string | HttpEndpoint) => {
    const key = getRpcEndpointKey(rpcEndpoint);
    if (!key) return;
    if (_rpcClients.hasOwnProperty(key)) {
        return _rpcClients[key];
    }
    const tmClient = await Tendermint34Client.connect(rpcEndpoint);
    //@ts-ignore
    const client = new QueryClient(tmClient);
    const rpc = createProtobufRpcClient(client);
    _rpcClients[key] = rpc;
    return rpc;
}
`;