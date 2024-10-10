import { Network } from '../src/clients/constants';
import { UserError } from '../src/clients/lib/errors';
import {
  connect,
  connectWallet,
  deposit,
  getAccountBalances,
  getOptimalNode,
  getUserStats,
  simulateDeposit,
  simulateTransferNativeToken,
  simulateWithdraw,
  transferNativeToken,
  withdraw,
  withdrawToIBC,
  wrappedError,
} from '../src/clients/native';
import { DYDX_TEST_ADDRESS, DYDX_TEST_MNEMONIC } from './constants';

async function test(): Promise<void> {
  try {
    const paramsInJson = `{
      "endpointUrls":[
        "https://dydx-testnet.nodefleet.org",
        "https://test-dydx.kingnodes.com",
        "https://dydx-rpc.liquify.com/api=8878132/dydx"
      ],
      "chainId":"dydx-testnet-4"
    }`;
    const result = await getOptimalNode(paramsInJson);
    console.log(result);

    const wallet = await connectWallet(DYDX_TEST_MNEMONIC);
    console.log(wallet);

    const address = await connect(Network.testnet(), DYDX_TEST_MNEMONIC);
    console.log(address);

    const payload = `{ "address": "${DYDX_TEST_ADDRESS}" }`;
    const userStats = await getUserStats(payload);
    console.log(userStats);

    const sendTokenPayload = {
      subaccountNumber: 0,
      amount: '10',   // Dydx Token
      recipient: 'dydx15ndn9c895f8ntck25qughtuck9spv2d9svw5qx',
    };
    const fees = await simulateTransferNativeToken(JSON.stringify(sendTokenPayload));
    console.log(fees);

    let tx = await transferNativeToken(JSON.stringify(sendTokenPayload));
    console.log(tx);

    let balances = await getAccountBalances();
    console.log(balances);

    const simulatePayload = {
      subaccountNumber: 0,
      amount: 20,   // In TDAI i.e. $20.00
    };
    let stdFee = await simulateWithdraw(JSON.stringify(simulatePayload));
    console.log(stdFee);

    const withdrawlPayload = {
      subaccountNumber: 0,
      amount: 20,
    };
    tx = await withdraw(JSON.stringify(withdrawlPayload));
    console.log(tx);

    balances = await getAccountBalances();
    console.log(balances);

    await new Promise((resolve) => setTimeout(resolve, 1000));

    const depositPayload = {
      subaccountNumber: 0,
      amount: 20,
    };
    const stringVal = JSON.stringify(depositPayload);

    stdFee = await simulateDeposit(stringVal);
    console.log(stdFee);

    tx = await deposit(stringVal);
    console.log(tx);

    // Use sDAI denom for ibc transfers
    const squidPayload = `
    {
      "msgTypeUrl": "/ibc.applications.transfer.v1.MsgTransfer",
      "msg": {
          "sourcePort": "transfer",
          "sourceChannel": "channel-0",
          "token": {
              "denom": "ibc/DEEFE2DEFDC8EA8879923C4CCA42BB888C3CD03FF7ECFEFB1C2FEC27A732ACC8",
              "amount": "10000000000000000000"
          },
          "sender": "dydx16zfx8g4jg9vels3rsvcym490tkn5la304c57e9",
          "receiver": "noble16zfx8g4jg9vels3rsvcym490tkn5la305z0jpu",
          "timeoutTimestamp": {
              "low": -1208865792,
              "high": 393844701,
              "unsigned": false
          },
          "memo": "{\\"forward\\":{\\"receiver\\":\\"osmo1zl9ztmwe2wcdvv9std8xn06mdaqaqm789rutmazfh3z869zcax4sv0ctqw\\",\\"port\\":\\"transfer\\",\\"channel\\":\\"channel-10\\",\\"next\\":{\\"wasm\\":{\\"contract\\":\\"osmo1zl9ztmwe2wcdvv9std8xn06mdaqaqm789rutmazfh3z869zcax4sv0ctqw\\",\\"msg\\":{\\"swap_with_action\\":{\\"swap_msg\\":{\\"token_out_min_amount\\":\\"26039154\\",\\"path\\":[{\\"pool_id\\":\\"46\\",\\"token_out_denom\\":\\"ibc/6F34E1BD664C36CE49ACC28E60D62559A5F96C4F9A6CCE4FC5A67B2852E24CFE\\"}]},\\"after_swap_action\\":{\\"ibc_transfer\\":{\\"receiver\\":\\"axelar1dv4u5k73pzqrxlzujxg3qp8kvc3pje7jtdvu72npnt5zhq05ejcsn5qme5\\",\\"channel\\":\\"channel-3\\",\\"next_memo\\":{\\"destination_chain\\":\\"ethereum-2\\",\\"destination_address\\":\\"0x481A2AAE41cd34832dDCF5A79404538bb2c02bC8\\",\\"payload\\":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,64,0,0,0,0,0,0,0,0,0,0,0,0,49,218,140,19,124,78,181,10,51,8,133,105,138,128,201,57,254,53,175,138,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,3,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,96,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,96,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,2,224,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,3,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,160,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,192,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,32,0,0,0,0,0,0,0,0,0,0,0,0,37,77,6,243,59,220,91,142,224,91,46,164,114,16,126,48,2,38,101,154,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,37,77,6,243,59,220,91,142,224,91,46,164,114,16,126,48,2,38,101,154,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,160,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,32,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,68,9,94,167,179,0,0,0,0,0,0,0,0,0,0,0,0,122,37,13,86,48,180,207,83,151,57,223,44,93,172,180,198,89,242,72,141,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,64,0,0,0,0,0,0,0,0,0,0,0,0,37,77,6,243,59,220,91,142,224,91,46,164,114,16,126,48,2,38,101,154,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,122,37,13,86,48,180,207,83,151,57,223,44,93,172,180,198,89,242,72,141,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,160,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,224,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,4,24,203,175,229,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,148,53,113,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,29,55,213,140,136,107,36,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,160,0,0,0,0,0,0,0,0,0,0,0,0,49,218,140,19,124,78,181,10,51,8,133,105,138,128,201,57,254,53,175,138,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,137,180,65,29,110,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,2,0,0,0,0,0,0,0,0,0,0,0,0,37,77,6,243,59,220,91,142,224,91,46,164,114,16,126,48,2,38,101,154,0,0,0,0,0,0,0,0,0,0,0,0,180,251,242,113,20,63,79,191,123,145,165,222,211,24,5,228,43,34,8,214,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,64,0,0,0,0,0,0,0,0,0,0,0,0,37,77,6,243,59,220,91,142,224,91,46,164,114,16,126,48,2,38,101,154,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],\\"type\\":2}}},\\"local_fallback_address\\":\\"osmo16zfx8g4jg9vels3rsvcym490tkn5la3056f20q\\"}}}}}}"
      }
  }`;
    console.log(squidPayload);

    const encode = (str: string):string => Buffer.from(str, 'binary').toString('base64');
    const encoded = encode(squidPayload);

    tx = await withdrawToIBC(0, '13', encoded);
    console.log(tx);
  } catch (error) {
    console.log(error.message);
  }
}

test().then(() => {
}).catch((error) => {
  console.log(error.message);
});

const error = new UserError('client is not connected. Call connectClient() first');
const text = wrappedError(error);
console.log(text);
