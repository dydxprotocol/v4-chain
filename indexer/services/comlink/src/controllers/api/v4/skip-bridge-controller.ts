import { dbHelpers } from '@dydxprotocol-indexer/postgres';
import { create, findByEvmAddress, findBySvmAddress } from '@dydxprotocol-indexer/postgres/build/src/stores/turnkey-users-table';
import { TurnkeyUserFromDatabase } from '@dydxprotocol-indexer/postgres/build/src/types';
import {
  route, executeRoute, setClientOptions, messages,
  RouteResponse,
} from '@skip-go/client/cjs';
import { Adapter } from '@solana/wallet-adapter-base';
import { Keypair, Transaction } from '@solana/web3.js';
import { Chain, createPublicClient, createWalletClient, encodeFunctionData, Hex, http, parseEther, PublicClient } from 'viem';
import { Turnkey } from '@turnkey/sdk-server';
import { TurnkeySigner } from '@turnkey/solana';
import { createAccount } from '@turnkey/viem';
import { decode, encode } from 'bech32';
import bs58 from 'bs58';
import express from 'express';
import {
  Controller, Post, Query, Route,
} from 'tsoa';
import nacl from 'tweetnacl';
import {
  mainnet, arbitrum,
  sepolia,
} from 'viem/chains';

import config from '../../../config';
import { handleControllerError } from '../../../lib/helpers';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { createKernelAccount, createKernelAccountClient, createZeroDevPaymasterClient, gasTokenAddresses, getUserOperationGasPrice } from '@zerodev/sdk';
import { getEntryPoint, KERNEL_V3_3, KERNEL_V3_1 } from '@zerodev/sdk/constants';
import { create7702KernelAccount, create7702KernelAccountClient, signerToEcdsaValidator } from '@zerodev/ecdsa-validator';
import { CheckBridgeSchema } from '../../../lib/validation/schemas';
// import { privateKeyToAccount } from 'viem/accounts';

const router = express.Router();
const controllerName: string = 'bridging-controller';

setClientOptions();


const publicClients: Record<string, PublicClient> = {
  [mainnet.id.toString()]: createPublicClient({
    transport: http(config.ZERODEV_RPC),
    chain: mainnet
  }),
  [arbitrum.id.toString()]: createPublicClient({
    transport: http(config.ZERODEV_RPC),
    chain: arbitrum
  }),
}

const chains: Record<string, Chain> = {
  [mainnet.id.toString()]: mainnet,
  [arbitrum.id.toString()]: arbitrum,
}

const turnkeySenderClient = new Turnkey({
  apiBaseUrl: config.TURNKEY_API_BASE_URL as string,
  apiPublicKey: config.TURNKEY_API_SENDER_PUBLIC_KEY as string,
  apiPrivateKey: config.TURNKEY_API_SENDER_PRIVATE_KEY as string,
  defaultOrganizationId: config.TURNKEY_ORGANIZATION_ID,
});

interface BridgeResponse {
  toAddress: string,
  amount: string,
  asset: string,
}

const usdcAddressByChainId: Record<string, string> = {
  [mainnet.id.toString()]: '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48', // usdc on ethereum mainnet.
  [arbitrum.id.toString()]: '0xaf88d065e77c8cC2239327C5EDb3A432268e5831', // usdc on arbitrum.
  [sepolia.id.toString()]: '0xaf88d065e77c8cC2239327C5EDb3A432268e5831', // usdc on sepolia.
  'solana': 'EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v', // usdc on solana.
  'dydx-mainnet-1': 'ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5', // usdc on dydx.
};

enum Asset {
  USDC = 'USDC',
  ETH = 'DAI',
}

const assetAddressLookerUpper: Record<Asset, Record<string, string>> = {
  [Asset.USDC]: usdcAddressByChainId,
  [Asset.ETH]: {}, // TODO: Add DAI address mappings when needed
}

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

async function getSkipRouteData(sourceAddress: string, dydxAddress: string, amount: string, chainId: string): Promise<{data: string,
toAddress: string}>{
  const routeResult = await route({
    amountIn: amount, // Desired amount in smallest denomination (e.g., uatom)
    sourceAssetDenom: usdcAddressByChainId[chainId], // USDC on mainnet. TODO: GENERALIZE
    sourceAssetChainId: chainId,
    destAssetDenom: usdcAddressByChainId['dydx-mainnet-1'],
    destAssetChainId: 'dydx-mainnet-1',
    cumulativeAffiliateFeeBps: '0',
    goFast: true,
  });
  if (!routeResult) {
    throw new Error("Failed to find a route");
  }

  console.log('usdcAddressByChainId[chainId]', usdcAddressByChainId[chainId]);
  console.log('Route Result:', routeResult);

  const userAddresses = await Promise.all(
    routeResult!.requiredChainAddresses.map(async (chainId) => ({
      chainId,
      address: await getAddress(chainId, sourceAddress, dydxAddress),
    }))
  );

  // getting the route data as shown here: https://github.com/skip-mev/skip-go/blob/a8907b389fa27fa942c42abca9181c0b0eee98e1/packages/client/src/public-functions/executeRoute.ts#L126
  let addressList: string[] = [];
  userAddresses.forEach((userAddress, index) => {
    const requiredChainAddress = routeResult.requiredChainAddresses[index];

    if (requiredChainAddress === userAddress?.chainId) {
      addressList.push(userAddress.address);
    }
  });

  if (addressList.length !== routeResult.requiredChainAddresses.length) {
    addressList = userAddresses.map((x) => x.address);
  }

  const validLength =
    addressList.length === routeResult.requiredChainAddresses.length ||
    addressList.length === routeResult.chainIds?.length;

  if (!validLength) {
    throw new Error("executeRoute error: invalid address list");
  }

  const timeoutSeconds = '60'; // Set a timeout for the messages request
  const response = await messages({
    timeoutSeconds,
    amountIn: routeResult?.amountIn,
    amountOut: routeResult.estimatedAmountOut || '0',
    sourceAssetChainId: routeResult?.sourceAssetChainId,
    sourceAssetDenom: routeResult?.sourceAssetDenom,
    destAssetChainId: routeResult?.destAssetChainId,
    destAssetDenom: routeResult?.destAssetDenom,
    operations: routeResult?.operations,
    addressList,
    slippageTolerancePercent: '1',
  });

  let data = '';
  let toAddress = '';
  response?.msgs?.forEach((msg, index) => {
    if ('evmTx' in msg) {
      console.log(`Message ${index + 1} EVM Transaction:`, msg.evmTx);
      data = msg.evmTx.data || '';
      toAddress = msg.evmTx.to || '';
    } else if ('svmTx' in msg) {
      console.log(`Message ${index + 1} SVM Transaction:`, msg.svmTx);
    } else if ('multiChainMsg' in msg) {
      console.log(`Message ${index + 1} Multi-Chain Message:`, msg.multiChainMsg);
    }
  });
  return { data, toAddress };
}

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

function getEvmSigner(suborgId: string, signWith: string) {
  return async (chainId: string) => {
    const serverClient = new Turnkey({
      apiBaseUrl: config.TURNKEY_API_BASE_URL as string,
      apiPublicKey: config.TURNKEY_API_SENDER_PUBLIC_KEY as string,
      apiPrivateKey: config.TURNKEY_API_SENDER_PRIVATE_KEY as string,
      defaultOrganizationId: suborgId,
    });

    const chain = {
      [mainnet.id.toString()]: mainnet,
      [arbitrum.id.toString()]: arbitrum,
    }[chainId];
    if (!chain) {
      throw new Error(`Unsupported chainId: ${chainId}`);
    }

    
    const turnkeyAccount = await createAccount({
      // @ts-ignore
      client: serverClient.apiClient(),
      organizationId: suborgId,
      signWith,
    });

    const rpcUrls: Record<string, string> = {
      [mainnet.id.toString()]: `https://eth-mainnet.g.alchemy.com/v2/${config.ALCHEMY_KEY}`,
      [sepolia.id.toString()]: `https://eth-sepolia.g.alchemy.com/v2/${config.ALCHEMY_KEY}`,
      [arbitrum.id.toString()]: `https://arb-mainnet.g.alchemy.com/v2/${config.ALCHEMY_KEY}`,
    };

    const rpcUrl = rpcUrls[chainId];
    if (!rpcUrl) {
      throw new Error(`No RPC URL configured for chainId: ${chainId}`);
    }

    return createWalletClient({
      account: turnkeyAccount,
      chain,
      transport: http(rpcUrl),
    });
  };
}

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
      @Query() asset: Asset,
      @Query() chainId: string,
  ): Promise<BridgeResponse> {
    if (chainId === 'solana') {
      return this.startSolanaBridge(fromAddress, amount, asset);
    }
    if (chainId === mainnet.id.toString() || chainId === arbitrum.id.toString()) {
      return this.startEvmBridge(fromAddress, amount, asset, chainId);
    }
    throw new Error(`Unsupported chainId: ${chainId}`);
  }

  async startSolanaBridge(
    fromAddress: string,
    amount: string,
    asset: Asset,
  ): Promise<BridgeResponse> {
    const chainId = 'solana';
    const addressLookUpper = assetAddressLookerUpper[asset];
    if (!addressLookUpper) {
      throw new Error(`Unsupported asset: ${asset}`);
    }
    const usdcAddress = addressLookUpper[chainId];
    if (!usdcAddress) {
      throw new Error(`Unsupported chainId: ${chainId} for asset: ${asset}`);
    }
    const path = await route({
      goFast: true,
      amountIn: amount,
      sourceAssetDenom: usdcAddress,
      sourceAssetChainId: chainId,
      destAssetChainId: 'dydx-mainnet-1',
      cumulativeAffiliateFeeBps: '0',
      destAssetDenom: usdcAddressByChainId['dydx-mainnet-1'],
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

    // find the suborgId for the user
    const record: TurnkeyUserFromDatabase | undefined = await findBySvmAddress(fromAddress);

    console.log('userAddresses is ', userAddresses);
    console.log('executing transaction...');
    // Replace with your own private key
    const solanaSponsorPrivateKey = '';
    const sponsorKeypair = Keypair.fromSecretKey(
      bs58.decode(solanaSponsorPrivateKey),
    );
    console.log('sponsorKeypair is ', sponsorKeypair.publicKey.toString());
    await executeRoute({
      route: path,
      userAddresses,
      getEvmSigner: getEvmSigner(record?.suborg_id || '', fromAddress),
      getSvmSigner: getSvmSigner(record?.suborg_id || '', fromAddress),
      svmFeePayer: {
        address: sponsorKeypair.publicKey.toString(), // Replace with the fee payer's Solana address
        signTransaction: (dataToSign: Buffer) => {
          const data = new Uint8Array(dataToSign);
          return Promise.resolve(nacl.sign.detached(data, sponsorKeypair.secretKey));
        },
      },
      // TO DO: ADD RETRY LOGIC
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

  async startEvmBridge(
    fromAddress: string,
    amount: string,
    asset: Asset,
    chainId: string,
  ): Promise<BridgeResponse> {
    const record: TurnkeyUserFromDatabase | undefined = await findByEvmAddress(fromAddress);
    if (!record || !record.dydx_address) {
      throw new Error('Failed to derive dYdX address');
    }
    const skipRoute = await getSkipRouteData(fromAddress, record.dydx_address, amount, chainId);
    if (!skipRoute) {
        console.error("Failed to get Skip route data");
        process.exit(1);
    }
    const skipRouteData = skipRoute.data;
    const skipToAddress = skipRoute.toAddress;
    console.log("Skip Route Data:", skipRouteData);

    const entryPoint = getEntryPoint("0.7");

    // Initialize a Turnkey-powered Viem Account
    const turnkeyAccount = await createAccount({
      // @ts-ignore
      client: turnkeySenderClient.apiClient(),
      organizationId: record.suborg_id,
      signWith: fromAddress,
    });

    // kernel account
    const account = await createKernelAccount(publicClients[chainId], {
        eip7702Account: turnkeyAccount,
        entryPoint,
        kernelVersion: KERNEL_V3_3,
    })
    console.log("account", account.address);

    const zerodevPaymaster = createZeroDevPaymasterClient({
      chain: chains[chainId],
      transport: http(config.ZERODEV_RPC),
    });

    const kernelClient = createKernelAccountClient({
        account,
        chain: chains[chainId],
        client: publicClients[chainId],
        bundlerTransport: http(config.ZERODEV_RPC),
        paymaster: {
            getPaymasterData: async (userOperation) => {
                return zerodevPaymaster.sponsorUserOperation({ userOperation })
            },

        },
        paymasterContext: { token: assetAddressLookerUpper[asset as Asset][chainId] },
        userOperation: {
            estimateFeesPerGas: async ({ bundlerClient }) => {
                return getUserOperationGasPrice(bundlerClient)
            }
        }
    })

    const userOpHash = await kernelClient.sendUserOperation({
        callData: await kernelClient.account.encodeCalls([
            {
                to: usdcAddressByChainId[chainId] as `0x${string}`,
                value: BigInt(0),
                data: encodeFunctionData({
                  abi: [
                      {
                          name: 'approve',
                          type: 'function',
                          stateMutability: 'nonpayable',
                          inputs: [
                              { name: 'spender', type: 'address' },
                              { name: 'amount', type: 'uint256' },
                          ],
                          outputs: [{ name: '', type: 'bool' }],
                      },
                  ],
                  functionName: 'approve',
                  args: [
                      (skipToAddress.startsWith("0x") ? skipToAddress : ("0x" + skipToAddress)) as Hex,
                      BigInt(amount) // 5.5 USDC (6 decimals)
                  ],
              }), //"0x",
            },
            {
                to: (skipToAddress.startsWith("0x") ? skipToAddress : ("0x" + skipToAddress)) as Hex,
                value: BigInt(0),
                data: skipRouteData.startsWith("0x") ? skipRouteData as Hex : ("0x" + skipRouteData) as Hex, //"0x",
            },
        ]),
    });
    console.log("UserOp sent:", userOpHash);
    console.log("Waiting for UserOp to be completed...");

    const { receipt } = await kernelClient.waitForUserOperationReceipt({
        hash: userOpHash,
    });
    console.log(
        "UserOp completed",
        `tx/${receipt.transactionHash}`
    );
    return {
      toAddress: fromAddress, 
      amount,
      asset,
    };
  }

  async startEvmBridgePre7702(
    fromAddress: string,
    amount: string,
    asset: Asset,
    chainId: string,
  ): Promise<BridgeResponse> {
    const record: TurnkeyUserFromDatabase | undefined = await findByEvmAddress(fromAddress);
    if (!record || !record.dydx_address) {
      throw new Error('Failed to derive dYdX address');
    }
    const skipRoute = await getSkipRouteData(fromAddress, record.dydx_address, amount, chainId);
    if (!skipRoute) {
        console.error("Failed to get Skip route data");
        process.exit(1);
    }
    const skipRouteData = skipRoute.data;
    const skipToAddress = skipRoute.toAddress;
    console.log("Skip Route Data:", skipRouteData);

    const entryPoint = getEntryPoint("0.7");

    // Initialize a Turnkey-powered Viem Account
    const turnkeyAccount = await createAccount({
      // @ts-ignore
      client: turnkeySenderClient.apiClient(),
      organizationId: record.suborg_id,
      signWith: fromAddress,
    });

  
    // Construct a validator
    const ecdsaValidator = await signerToEcdsaValidator(publicClients[chainId], {
      signer: turnkeyAccount,
      entryPoint,
      kernelVersion: KERNEL_V3_1,
    });

    // kernel account
    const account = await createKernelAccount(publicClients[chainId], {
        entryPoint,
        plugins: {
          sudo: ecdsaValidator,
        },
        kernelVersion: KERNEL_V3_1,
    })
    console.log("account", account.address);

    const zerodevPaymaster = createZeroDevPaymasterClient({
      chain: chains[chainId],
      transport: http(config.ZERODEV_RPC),
    });

    const kernelClient = createKernelAccountClient({
        account,
        chain: chains[chainId],
        client: publicClients[chainId],
        bundlerTransport: http(config.ZERODEV_RPC),
        paymaster: {
            getPaymasterData: async (userOperation) => {
                return zerodevPaymaster.sponsorUserOperation({ userOperation })
            },

        },
        // paymasterContext: { token: assetAddressLookerUpper[asset as Asset][chainId] },
        userOperation: {
            estimateFeesPerGas: async ({ bundlerClient }) => {
                return getUserOperationGasPrice(bundlerClient)
            }
        }
    })

    const userOpHash = await kernelClient.sendUserOperation({
        callData: await kernelClient.account.encodeCalls([
            {
                to: usdcAddressByChainId[chainId] as `0x${string}`,
                value: BigInt(0),
                data: encodeFunctionData({
                  abi: [
                      {
                          name: 'approve',
                          type: 'function',
                          stateMutability: 'nonpayable',
                          inputs: [
                              { name: 'spender', type: 'address' },
                              { name: 'amount', type: 'uint256' },
                          ],
                          outputs: [{ name: '', type: 'bool' }],
                      },
                  ],
                  functionName: 'approve',
                  args: [
                      (skipToAddress.startsWith("0x") ? skipToAddress : ("0x" + skipToAddress)) as Hex,
                      BigInt(amount)
                  ],
              }), //"0x",
            },
            {
                to: (skipToAddress.startsWith("0x") ? skipToAddress : ("0x" + skipToAddress)) as Hex,
                value: BigInt(0),
                data: skipRouteData.startsWith("0x") ? skipRouteData as Hex : ("0x" + skipRouteData) as Hex, //"0x",
            },
        ]),
    });
    console.log("UserOp sent:", userOpHash);
    console.log("Waiting for UserOp to be completed...");

    const { receipt } = await kernelClient.waitForUserOperationReceipt({
        hash: userOpHash,
    });
    console.log(
        "UserOp completed",
        `tx/${receipt.transactionHash}`
    );
    return {
      toAddress: fromAddress, 
      amount,
      asset,
    };
  }
}

router.post(
  '/startBridge',
  ...CheckBridgeSchema,
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    // alchemy don't care.
    res.status(200).send();
    try {
      const { fromAddress, amount, asset, chainId } = req.body;
      const bridgeController = new BridgeController();
      await bridgeController.startBridge(
        fromAddress,
        amount,
        asset,
        chainId,
      );
      return; 
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
