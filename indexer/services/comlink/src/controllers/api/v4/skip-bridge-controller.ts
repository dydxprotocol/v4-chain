import { dbHelpers } from '@dydxprotocol-indexer/postgres';
import { create, findByEvmAddress, findBySvmAddress } from '@dydxprotocol-indexer/postgres/build/src/stores/turnkey-users-table';
import { TurnkeyUserFromDatabase } from '@dydxprotocol-indexer/postgres/build/src/types';
import {
  route, executeRoute, setClientOptions,
} from '@skip-go/client/cjs';
import { Adapter } from '@solana/wallet-adapter-base';
import { Keypair, Transaction } from '@solana/web3.js';
// import { TurnkeyClient } from '@turnkey/http';
import { Turnkey } from '@turnkey/sdk-server';
import { TurnkeySigner } from '@turnkey/solana';
// import { createAccount } from '@turnkey/viem';
import { decode, encode } from 'bech32';
import bs58 from 'bs58';
import express from 'express';
import { checkSchema } from 'express-validator';
import fetch, { Headers, Request, Response } from 'node-fetch';
import {
  Controller, Post, Query, Route,
} from 'tsoa';
import nacl from 'tweetnacl';
import {
  mainnet, arbitrum,
} from 'viem/chains';

import config from '../../../config';
import { handleControllerError } from '../../../lib/helpers';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
// import { privateKeyToAccount } from 'viem/accounts';

// Polyfill Headers for Node.js 16
if (!globalThis.fetch) {
  // @ts-ignore
  globalThis.fetch = fetch;
  // @ts-ignore
  globalThis.Headers = Headers;
  // @ts-ignore
  globalThis.Request = Request;
  // @ts-ignore
  globalThis.Response = Response;
}

// Add global declaration for types
declare global {
  interface Array<T> {
    findLast(predicate: (value: T, index: number, obj: T[]) => unknown): T | undefined,
  }
}

// Polyfill
if (!Array.prototype.findLast) {
  // eslint-disable-next-line no-extend-native
  Array.prototype.findLast = function <T>(
    this: T[],
    callback: (element: T, index: number, array: T[]) => unknown,
  ): T | undefined {
    if (this == null) {
      throw new TypeError('this is null or not defined');
    }
    const len = this.length;
    for (let i = len - 1; i >= 0; i--) {
      if (callback(this[i], i, this)) {
        return this[i];
      }
    }
    return undefined;
  };
}

const router = express.Router();
const controllerName: string = 'bridging-controller';

setClientOptions();

interface BridgeResponse {
  toAddress: string,
  amount: string,
  asset: string,
}

const assetMap: Record<string, string> = {
  dydx_USDC: 'ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5', // usdc on dydx.
  sol_USDC: 'EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v', // usdc on solana.
  eth_USDC: '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48', // usdc on ethereum mainnet.
  arb_USDC: '0xaf88d065e77c8cC2239327C5EDb3A432268e5831', // usdc on arbitrum.
};

async function getDydxAddress(address: string, chainId: string): Promise<string> {
  let dydxAddress = '';
  if (chainId === mainnet.id.toString() || chainId === arbitrum.id.toString()) {
    // look up in turnkey table
    console.log('looking up in turnkey table for evm address', address);
    const record = await findByEvmAddress(address);
    dydxAddress = record?.dydx_address || '';
  } else if (chainId === 'solana') {
    // look up in turnkey table
    console.log('looking up in turnkey table for svm address', address);
    const record = await findBySvmAddress(address);
    dydxAddress = record?.dydx_address || '';
  } else {
    throw new Error(`Unsupported chainId: ${chainId}`);
  }
  return dydxAddress;
}

function getAddress(
  chainId: string,
  sourceAddress: string,
  dydxAddress: string,
): string {
  switch (chainId) {
    case 'solana':
    case '42161': // arbitrum
    case '11155111': // sepolia
    case '1': // mainnet
    case '5': // goerli
    case '43114': // avalanche
    case '137': // polygon
    case '8453': // base
    case '10': // optimism
      return sourceAddress;
    case 'noble-1':
      return toNobleAddress(dydxAddress) || '';
    case 'osmosis-1':
      return toOsmosisAddress(dydxAddress) || '';
    case 'neutron':
      return toNeutronAddress(dydxAddress) || '';
    case 'dydx-mainnet-1':
      return dydxAddress;
    default:
      throw new Error(`Unsupported chain ID: ${chainId}`);
  }
}

// const amount = "1100000"; // 1 USDC in smallest denomination (6 decimals for USDC)

// async function getSkipRouteData(sourceAddress: string): Promise<{data: string,
// toAddress: string}>{
//   const routeResult = await route({
//     amountIn: amount, // Desired amount in smallest denomination (e.g., uatom)
//     sourceAssetDenom: assetMap.eth_USDC, // USDC on Base
//     sourceAssetChainId: mainnet.id.toString(),
//     destAssetDenom: "ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5",
//     destAssetChainId: "dydx-mainnet-1",
//     cumulativeAffiliateFeeBps: '0',
//     goFast: true,
//   });
//   if (!routeResult) {
//     throw new Error("Failed to find a route");
//   }
//   console.log('Route Result:', routeResult);

//   const userAddresses = await Promise.all(
//     routeResult!.requiredChainAddresses.map(async (chainId) => ({
//       chainId,
//       address: await getAddress(chainId, sourceAddress),
//     }))
//   );

//   // getting the route data as shown here: https://github.com/skip-mev/skip-go/blob/a8907b389fa27fa942c42abca9181c0b0eee98e1/packages/client/src/public-functions/executeRoute.ts#L126
//   let addressList: string[] = [];
//   userAddresses.forEach((userAddress, index) => {
//     const requiredChainAddress = routeResult.requiredChainAddresses[index];

//     if (requiredChainAddress === userAddress?.chainId) {
//       addressList.push(userAddress.address);
//     }
//   });

//   if (addressList.length !== routeResult.requiredChainAddresses.length) {
//     addressList = userAddresses.map((x) => x.address);
//   }

//   const validLength =
//     addressList.length === routeResult.requiredChainAddresses.length ||
//     addressList.length === routeResult.chainIds?.length;

//   if (!validLength) {
//     throw new Error("executeRoute error: invalid address list");
//   }

//   const timeoutSeconds = '60'; // Set a timeout for the messages request
//   const response = await messages({
//     timeoutSeconds,
//     amountIn: routeResult?.amountIn,
//     amountOut: routeResult.estimatedAmountOut || '0',
//     sourceAssetChainId: routeResult?.sourceAssetChainId,
//     sourceAssetDenom: routeResult?.sourceAssetDenom,
//     destAssetChainId: routeResult?.destAssetChainId,
//     destAssetDenom: routeResult?.destAssetDenom,
//     operations: routeResult?.operations,
//     addressList,
//     slippageTolerancePercent: '1',
//   });

//   console.log('Skip Route Data:', response);

//   let data = '';
//   let toAddress = '';
//   response?.msgs?.forEach((msg, index) => {
//     if ('evmTx' in msg) {
//       console.log(`Message ${index + 1} EVM Transaction:`, msg.evmTx);
//       data = msg.evmTx.data || '';
//       toAddress = msg.evmTx.to || '';
//     } else if ('svmTx' in msg) {
//       console.log(`Message ${index + 1} SVM Transaction:`, msg.svmTx);
//     } else if ('multiChainMsg' in msg) {
//       console.log(`Message ${index + 1} Multi-Chain Message:`, msg.multiChainMsg);
//     }
//   });
//   return { data, toAddress };
// }

// const srcAccountPrivateKey = '0x11ecbbfcc6fa1a7c1c0b00f59423b6934376a20f26430f9f15d2eafb0643d7a5'

// const eip7702Account = privateKeyToAccount(
//   sourceAccountPrivateKey, // generatePrivateKey() ?? (process.env.PRIVATE_KEY as Hex)
// );

function toNobleAddress(address: string): string | null {
  try {
    const decoded = decode(address);
    if (decoded.prefix !== 'dydx') {
      return null;
    }
    return encode('noble', decoded.words);
  } catch (e) {
    return null;
  }
}

function toOsmosisAddress(address: string): string | null {
  try {
    const decoded = decode(address);
    if (decoded.prefix !== 'dydx') {
      return null;
    }
    return encode('osmo', decoded.words);
  } catch (e) {
    return null;
  }
}

function toNeutronAddress(address: string): string | null {
  try {
    const decoded = decode(address);
    if (decoded.prefix !== 'dydx') {
      return null;
    }
    return encode('neutron', decoded.words);
  } catch (e) {
    return null;
  }
}

// function getEvmSigner(suborgId: string, signWith: string) {
//   return async (chainId: string) => {
//     // 1. Initialize Turnkey HTTP client
//     // const httpClient = new TurnkeyClient(
//     //   { baseUrl: config.TURNKEY_API_BASE_URL as string },
//     //   new ApiKeyStamper({
//     //     apiPublicKey: config.TURNKEY_API_SENDER_PUBLIC_KEY as string,
//     //     apiPrivateKey: config.TURNKEY_API_SENDER_PRIVATE_KEY as string,
//     //   }),
//     // );

//     const serverClient = new Turnkey({
//       apiBaseUrl: config.TURNKEY_API_BASE_URL as string,
//       apiPublicKey: config.TURNKEY_API_SENDER_PUBLIC_KEY as string,
//       apiPrivateKey: config.TURNKEY_API_SENDER_PRIVATE_KEY as string,
//       defaultOrganizationId: suborgId,
//     });

//     const chain = {
//       [mainnet.id.toString()]: mainnet,
//       [sepolia.id.toString()]: sepolia,
//       [arbitrum.id.toString()]: arbitrum,
//     }[chainId];
//     if (!chain) {
//       throw new Error(`Unsupported chainId: ${chainId}`);
//     }

//     const turnkeyAccount = await createAccount({
//       client: serverClient.apiClient(),
//       organizationId: suborgId,
//       signWith,
//     });

//     const rpcUrls: Record<string, string> = {
//       [mainnet.id.toString()]: `https://eth-mainnet.g.alchemy.com/v2/${config.ALCHEMY_KEY}`,
//       [sepolia.id.toString()]: `https://eth-sepolia.g.alchemy.com/v2/${config.ALCHEMY_KEY}`,
//       [arbitrum.id.toString()]: `https://arb-mainnet.g.alchemy.com/v2/${config.ALCHEMY_KEY}`,
//     };

//     const rpcUrl = rpcUrls[chainId];
//     if (!rpcUrl) {
//       throw new Error(`No RPC URL configured for chainId: ${chainId}`);
//     }

//     return createWalletClient({
//       account: turnkeyAccount,
//       chain,
//       transport: http(rpcUrl),
//     });
//   };
// }

function getSvmSigner(suborgId: string, signWith: string) {
  const serverClient = new Turnkey({
    apiBaseUrl: config.TURNKEY_API_BASE_URL as string,
    apiPublicKey: config.TURNKEY_API_SENDER_PUBLIC_KEY as string,
    apiPrivateKey: config.TURNKEY_API_SENDER_PRIVATE_KEY as string,
    defaultOrganizationId: suborgId,
  });

  const turnkeySigner = new TurnkeySigner({
    organizationId: suborgId,
    client: serverClient.apiClient(),
  });

  // eslint-disable-next-line @typescript-eslint/require-await
  return async () => ({
    publicKey: {
      toString: () => signWith,
      toBase58: () => signWith,
    },
    signTransaction: async (tx: Transaction) => {
      try {
        return await turnkeySigner.signTransaction(tx, signWith);
      } catch (error) {
        throw new Error(`Failed to sign transaction with TurnkeySigner: ${error.message}`);
      }
    },
    signAllTransactions: async (txs: Transaction[]) => {
      try {
        return await turnkeySigner.signAllTransactions(txs, signWith);
      } catch (error) {
        throw new Error(`Failed to sign transactions with TurnkeySigner: ${error.message}`);
      }
    },
  } as Adapter);
}

@Route('bridging')
class BridgeController extends Controller {
  @Post('/startBridge')
  async startBridge(
    @Query() fromAddress: string,
      @Query() amount: string,
      @Query() asset: string,
      @Query() chainId: string,
  ): Promise<BridgeResponse> {
    console.log('calculating path...');
    const path = await route({
      goFast: true,
      amountIn: amount,
      sourceAssetDenom: assetMap[asset],
      sourceAssetChainId: chainId,
      destAssetChainId: 'dydx-mainnet-1',
      cumulativeAffiliateFeeBps: '0',
      destAssetDenom: assetMap.dydx_USDC,
    });

    console.log('path is ', path);

    const dydxAddress = await getDydxAddress(fromAddress, chainId);
    if (!dydxAddress) {
      throw new Error('Failed to derive dYdX address');
    }

    if (!path) {
      throw new Error('Failed to create route');
    }

    // the end user must be able to sign these intermediate addresses as
    // they are the addresses that the funds will be deposited into in case
    // of a failure.
    const userAddresses = await Promise.all(
      path.requiredChainAddresses.map((chain: string) => ({
        chainId: chain,
        address: getAddress(chain, fromAddress, dydxAddress),
      })),
    );

    // search for suborgId,
    let record: TurnkeyUserFromDatabase | undefined;
    if (chainId === mainnet.id.toString()) {
      record = await findByEvmAddress(fromAddress);
    } else if (chainId === 'solana') {
      record = await findBySvmAddress(fromAddress);
    }

    console.log('userAddresses is ', userAddresses);
    console.log('executing transaction...');
    // Replace with your own private key
    const solanaSponsorPrivateKey = '3UJeupkPcz7Xc3QLQ96UDTv1N2rHpKXFRgXS3MW6XiqADGFHBdm7eS5G5aCVd8Nnf5xEnGLVv76dkosXx98Pjnwo';
    const sponsorKeypair = Keypair.fromSecretKey(
      bs58.decode(solanaSponsorPrivateKey),
    );
    console.log('sponsorKeypair is ', sponsorKeypair.publicKey.toString());
    await executeRoute({
      route: path,
      userAddresses,
      simulate: false,
      // getEvmSigner: getEvmSigner(record?.suborg_id || '', fromAddress),
      getSvmSigner: getSvmSigner(record?.suborg_id || '', fromAddress),
      svmFeePayer: {
        address: sponsorKeypair.publicKey.toString(), // Replace with the fee payer's Solana address
        signTransaction: (dataToSign: Buffer) => {
          const data = new Uint8Array(dataToSign);
          return Promise.resolve(nacl.sign.detached(data, sponsorKeypair.secretKey));
        },
      },
      onTransactionBroadcast: async ({ chainId: c, txHash }) => {
        await console.log(`Broadcasted on ${c}: ${txHash}`);
      },
      onTransactionCompleted: async ({ chainId: c, txHash, status }) => {
        await console.log(`Completed on ${c}: ${txHash} (Status: ${status})`);
      },
      onTransactionTracked: async ({ chainId: c, txHash, explorerLink }) => {
        await console.log(`Tracking ${c}: ${txHash} (Explorer: ${explorerLink})`);
      },
      onTransactionSignRequested: async ({ chainId: c, signerAddress }) => {
        await console.log(`Sign requested for ${c}`, signerAddress);
      },
      onValidateGasBalance: async ({ chainId: c, txIndex, status }) => {
        await console.log('validate: ', c, txIndex, status);
      },
    });

    // TODO: Implement bridge creation logic
    const bridge: BridgeResponse = {
      toAddress: fromAddress,
      amount,
      asset,
    };
    return bridge;
  }
}

router.post(
  '/startBridge',
  ...checkSchema({
    // Validate the event object structure
    event: {
      in: 'body',
      isObject: true,
      errorMessage: 'Event must be an object',
    },
    'event.activity': {
      in: 'body',
      isArray: true,
      errorMessage: 'Event.activity must be an array',
    },
    'event.activity.*.fromAddress': {
      in: 'body',
      isString: true,
      errorMessage: 'Activity fromAddress must be a string',
    },
    'event.activity.*.toAddress': {
      in: 'body',
      isString: true,
      optional: true,
      errorMessage: 'Activity toAddress must be a string',
    },
    'event.activity.*.asset': {
      in: 'body',
      isString: true,
      optional: true,
      errorMessage: 'Activity asset must be a string',
    },
    'event.activity.*.value': {
      in: 'body',
      isNumeric: true,
      optional: true,
      errorMessage: 'Activity value must be a number',
    },
    'event.network': {
      in: 'body',
      isString: true,
      optional: true,
      errorMessage: 'Event network must be a string',
    },
    // Webhook metadata
    id: {
      in: 'body',
      isString: true,
      optional: true,
    },
    type: {
      in: 'body',
      isString: true,
      optional: true,
    },
    webhookId: {
      in: 'body',
      isString: true,
      optional: true,
    },
  }),
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    try {

      console.log({
        apiPublicKey: config.TURNKEY_API_SENDER_PUBLIC_KEY,
        apiPrivateKey: config.TURNKEY_API_SENDER_PRIVATE_KEY,
        organizationId: config.TURNKEY_ORGANIZATION_ID,
        apiBaseUrl: config.TURNKEY_API_BASE_URL,
      });
      // create a sample db record.
      await dbHelpers.clearData();
      await create({
        evm_address: '0x46c9E748dfb814Da6577fD4ceF8f785CE7bB4Be7',
        svm_address: 'AuV1WxiP1bswKykhC9KB5J1Ek1xmq9AdZANWGP97hsPh',
        dydx_address: 'dydx1sjssdnatk99j2sdkqgqv55a8zs97fcvstzreex',
        suborg_id: '70528f95-da66-49f4-a096-fe50a19f2e6d',
        salt: '1234567890',
        created_at: new Date().toISOString(),
      });

      console.log({
        apiPublicKey: config.TURNKEY_API_SENDER_PUBLIC_KEY,
        apiPrivateKey: config.TURNKEY_API_SENDER_PRIVATE_KEY,
        organizationId: config.TURNKEY_ORGANIZATION_ID,
        apiBaseUrl: config.TURNKEY_API_BASE_URL,
      });

      const bridgeController = new BridgeController();
      console.log('starting bridge...');
      await bridgeController.startBridge(
        'AuV1WxiP1bswKykhC9KB5J1Ek1xmq9AdZANWGP97hsPh',
        '1000000', // 2 USDC
        'sol_USDC',
        'solana',
      );

      return res.status(200).send({
        success: true,
      });
    } catch (error) {
      return handleControllerError(
        'BridgeController POST /startBridge',
        'Bridge start error',
        error,
        req,
        res,
      );
    }
  },
);

export default router;
