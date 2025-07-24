import { create, findByEvmAddress, findBySvmAddress } from '@dydxprotocol-indexer/postgres/build/src/stores/turnkey-users-table';
import { TurnkeyUserFromDatabase } from '@dydxprotocol-indexer/postgres/build/src/types';
import { route, executeRoute, setClientOptions } from '@skip-go/client/cjs';
import { TurnkeyClient } from '@turnkey/http';
import { ApiKeyStamper, Turnkey as TurnkeyServerSDK } from '@turnkey/sdk-server';
import { createAccount } from '@turnkey/viem';
import { decode, encode } from 'bech32';
import express from 'express';
import { checkSchema } from 'express-validator';
import {
  Controller, Post, Query, Route,
} from 'tsoa';
import { http, createWalletClient } from 'viem';
import { mainnet, goerli, sepolia, arbitrum } from 'viem/chains';

import config from '../../../config';
import { handleControllerError } from '../../../lib/helpers';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { dbHelpers } from '@dydxprotocol-indexer/postgres';

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

async function getEvmSigner(suborgId: string) {

  // 1. Initialize Turnkey HTTP client
  const httpClient = new TurnkeyClient(
    { baseUrl: config.TURNKEY_API_BASE_URL },
    new ApiKeyStamper({
      apiPublicKey: config.TURNKEY_API_SENDER_PUBLIC_KEY!,
      apiPrivateKey: config.TURNKEY_API_SENDER_PRIVATE_KEY!,
    }),
  );

  // api client
  const apiClient = new TurnkeyServerSDK({
    apiBaseUrl: config.TURNKEY_API_BASE_URL,
    apiPrivateKey: config.TURNKEY_API_SENDER_PRIVATE_KEY,
    apiPublicKey: config.TURNKEY_API_SENDER_PUBLIC_KEY,
    defaultOrganizationId: config.TURNKEY_ORGANIZATION_ID,
  }).apiClient();

  // find the wallet account for the suborgId
  const wallets = await apiClient.getWallets({
    organizationId: suborgId,
  });

  // use the first wallet account.
  const wallet = wallets.wallets[0];
  if (!wallet) {
    throw new Error('No wallet found');
  }

  return async (chainId: string) => {
    const chain = {
      1: mainnet,
      5: goerli,
      11155111: sepolia,
    }[chainId];
    if (!chain) {
      throw new Error(`Unsupported chainId: ${chainId}`);
    }
    // 2. Create the Viem‐compatible Turnkey account
    const turnkeyAccount = await createAccount({
      client: httpClient,
      organizationId: config.TURNKEY_ORGANIZATION_ID!,
      signWith: wallet.walletId,
    });
    // 3. Create and return a WalletClient that uses Turnkey to sign
    return createWalletClient({
      account: turnkeyAccount,
      chain: mainnet,
      transport: http(`https://mainnet.infura.io/v3/${process.env.INFURA_KEY}`),
    });
  };
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
      record = await findByEvmAddress(userAddresses[0].address);
    }

    console.log('userAddresses is ', userAddresses);
    console.log('executing transaction...');

    await executeRoute({
      route: path,
      userAddresses,
      getEvmSigner: await getEvmSigner(record?.suborg_id || ''),
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
        evm_address: '0xEc58845E98c3C2bA4C83B83730EBD58C75433a97',
        svm_address: '7wMtPuTfupJBfqZdq5S8mTfaJeccwLdwPuvyV2jtZvko',
        dydx_address: 'dydx1sjssdnatk99j2sdkqgqv55a8zs97fcvstzreex',
        suborg_id: 'fed75fec-9243-4114-b58b-8bb0ce15cba9',
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
        '0xEc58845E98c3C2bA4C83B83730EBD58C75433a97',
        '2000000', // 2 USDC
        'arb_USDC',
        arbitrum.id.toString(),
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
