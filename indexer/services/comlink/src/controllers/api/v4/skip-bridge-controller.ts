import { logger, stats } from '@dydxprotocol-indexer/base';
import { findByEvmAddress, findBySvmAddress } from '@dydxprotocol-indexer/postgres/build/src/stores/turnkey-users-table';
import { TurnkeyUserFromDatabase } from '@dydxprotocol-indexer/postgres/build/src/types';
import {
  route, executeRoute, setClientOptions, messages,
} from '@skip-go/client/cjs';
import { Adapter } from '@solana/wallet-adapter-base';
import { Keypair, Transaction, VersionedTransaction } from '@solana/web3.js';
import { Turnkey } from '@turnkey/sdk-server';
import { TurnkeySigner } from '@turnkey/solana';
import { createAccount } from '@turnkey/viem';
import { signerToEcdsaValidator } from '@zerodev/ecdsa-validator';
import {
  createKernelAccount, createKernelAccountClient,
  createZeroDevPaymasterClient, getUserOperationGasPrice,
} from '@zerodev/sdk';
import { getEntryPoint, KERNEL_V3_3, KERNEL_V3_1 } from '@zerodev/sdk/constants';
import { Alchemy, Network } from 'alchemy-sdk';
import { decode, encode } from 'bech32';
import bs58 from 'bs58';
import express from 'express';
import {
  Controller, Post, Query, Route,
} from 'tsoa';
import nacl from 'tweetnacl';
import {
  Address,
  Chain, checksumAddress, createPublicClient, encodeFunctionData, Hex, http, PublicClient,
} from 'viem';
import {
  type SmartAccountImplementation,
} from 'viem/account-abstraction';
import {
  mainnet, arbitrum,
  avalanche, base,
  optimism,
} from 'viem/chains';

import config from '../../../config';
import { handleControllerError } from '../../../lib/helpers';
import { CheckBridgeSchema } from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';

const router = express.Router();
const controllerName: string = 'bridging-controller';

// set the skip client options to use the skip rpc.
setClientOptions({
  endpointOptions: {
    endpoints: {
      solana: {
        rpc: 'https://go.skip.build/api/rpc/solana',
      },
    },
  },
});

const chains: Record<string, Chain> = {
  [mainnet.id.toString()]: mainnet,
  [arbitrum.id.toString()]: arbitrum,
  [avalanche.id.toString()]: avalanche,
  [base.id.toString()]: base,
  [optimism.id.toString()]: optimism,
};

const chainIdToAlchemyNetworkMap: Record<string, Network> = {
  [arbitrum.id.toString()]: Network.ARB_MAINNET,
  [avalanche.id.toString()]: Network.AVAX_MAINNET,
  [base.id.toString()]: Network.BASE_MAINNET,
  [optimism.id.toString()]: Network.OPT_MAINNET,
  [mainnet.id.toString()]: Network.ETH_MAINNET,
  solana: Network.SOLANA_MAINNET,
};

const alchemyNetworkToChainIdMap: Record<string, string> = {
  ARB_MAINNET: arbitrum.id.toString(),
  AVAX_MAINNET: avalanche.id.toString(),
  BASE_MAINNET: base.id.toString(),
  OPT_MAINNET: optimism.id.toString(),
  ETH_MAINNET: mainnet.id.toString(),
  SOLANA_MAINNET: 'solana',
};

const publicClients = Object.keys(chains).reduce((acc, chainId) => {
  acc[chainId] = createPublicClient({
    transport: http(getRPCEndpoint(chainId)),
    chain: chains[chainId],
  });
  return acc;
}, {} as Record<string, PublicClient>);

function isSupportedEVMChainId(chainId: string): boolean {
  return Object.keys(chains).includes(chainId);
}

function getRPCEndpoint(chainId: string): string {
  if (!isSupportedEVMChainId(chainId)) {
    throw new Error(`Unsupported chainId: ${chainId}`);
  }
  return `${config.ZERODEV_API_BASE_URL}/${config.ZERODEV_API_KEY}/chain/${chainId}`;
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
  sourceAssetDenom: string,
}

const usdcAddressByChainId: Record<string, string> = {
  [mainnet.id.toString()]: '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48', // usdc on ethereum mainnet.
  [arbitrum.id.toString()]: '0xaf88d065e77c8cC2239327C5EDb3A432268e5831', // usdc on arbitrum.
  [avalanche.id.toString()]: '0xb97ef9ef8734c71904d8002f8b6bc66dd9c48a6e', // usdc on avalanche.
  [base.id.toString()]: '0x833589fcd6edb6e08f4c7c32d4f71b54bda02913', // usdc on base.
  [optimism.id.toString()]: '0x0b2c639c533813f4aa9d7837caf62653d097ff85', // usdc on optimism.
  solana: 'EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v', // usdc on solana.
  'dydx-mainnet-1': 'ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5', // usdc on dydx.
};

enum Asset {
  USDC = 'USDC',
  ETH = 'ETH',
}

enum CosmosPrefix {
  OSMO = 'osmo',
  NEUTRON = 'neutron',
  NOBLE = 'noble',
}

// Prefix is one of osmosis, neutron, noble. This is how we convert dydx addresses
// to other chain addresses on cosmos. Address here is dydx address.
function toClientAddressWithPrefix(prefix: CosmosPrefix, address: string): string | null {
  try {
    const decoded = decode(address);
    if (decoded.prefix !== 'dydx') {
      return null;
    }
    return encode(prefix, decoded.words);
  } catch (e) {
    return null;
  }
}

// Finds the dydx address for a given evm or svm address.
async function getDydxAddress(address: string, chainId: string): Promise<string> {
  let dydxAddress = '';
  if (isSupportedEVMChainId(chainId)) {
    // look up in turnkey table
    logger.info({
      at: `${controllerName}#getDydxAddress`,
      message: 'Looking up in turnkey table for evm address',
      address,
    });
    const record = await findByEvmAddress(address);
    dydxAddress = record?.dydx_address || '';
  } else if (chainId === 'solana') {
    // look up in turnkey table
    logger.info({
      at: `${controllerName}#getDydxAddress`,
      message: 'Looking up in turnkey table for svm address',
      address,
    });
    const record = await findBySvmAddress(address);
    dydxAddress = record?.dydx_address || '';
  } else {
    throw new Error(`Unsupported chainId: ${chainId}`);
  }
  return dydxAddress;
}

// TODO: Verify that this function is 1000% correct. @RUI and @TYLER and @JARED
function getAddress(
  chainId: string,
  sourceAddress: string,
  dydxAddress: string,
): string {
  if (isSupportedEVMChainId(chainId) || chainId === 'solana') {
    return sourceAddress;
  }
  switch (chainId) {
    case 'noble-1':
      return toClientAddressWithPrefix(CosmosPrefix.NOBLE, dydxAddress) || '';
    case 'osmosis-1':
      return toClientAddressWithPrefix(CosmosPrefix.OSMO, dydxAddress) || '';
    case 'neutron':
      return toClientAddressWithPrefix(CosmosPrefix.NEUTRON, dydxAddress) || '';
    case 'dydx-mainnet-1':
      return dydxAddress;
    default:
      throw new Error(`Unsupported chain ID: ${chainId}`);
  }
}

// Grabs the raw skip route data to carry out the bridge on our own.
async function getSkipCallData(
  sourceAddress: string,
  sourceAssetDenom: string,
  dydxAddress: string,
  amount: string,
  chainId: string,
): Promise<Parameters<SmartAccountImplementation['encodeCalls']>[0]> {
  const routeResult = await route({
    amountIn: amount, // Desired amount in smallest denomination (e.g., uatom)
    sourceAssetDenom,
    sourceAssetChainId: chainId,
    destAssetDenom: usdcAddressByChainId['dydx-mainnet-1'],
    destAssetChainId: 'dydx-mainnet-1',
    cumulativeAffiliateFeeBps: '0',
    goFast: true,
  });
  if (!routeResult) {
    throw new Error('Failed to find a route');
  }

  logger.info({
    at: `${controllerName}#getSkipCallData`,
    message: 'Route result obtained',
    routeResult,
    dydxAddress,
  });

  const userAddresses = await Promise.all(
    routeResult.requiredChainAddresses.map(async (cid) => ({
      chainId: cid,
      address: await getAddress(cid, sourceAddress, dydxAddress),
    })),
  );

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

  const validLength = addressList.length === routeResult.requiredChainAddresses.length ||
    addressList.length === routeResult.chainIds?.length;

  if (!validLength) {
    throw new Error('executeRoute error: invalid address list');
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
      logger.info({
        at: `${controllerName}#getSkipCallData`,
        message: `Message ${index + 1} EVM Transaction`,
        evmTx: msg.evmTx,
      });
      data = msg.evmTx.data || '';
      toAddress = msg.evmTx.to || '';
    } else if ('svmTx' in msg) {
      logger.info({
        at: `${controllerName}#getSkipCallData`,
        message: `Message ${index + 1} SVM Transaction`,
        svmTx: msg.svmTx,
      });
    } else if ('multiChainMsg' in msg) {
      logger.info({
        at: `${controllerName}#getSkipCallData`,
        message: `Message ${index + 1} Multi-Chain Message`,
        multiChainMsg: msg.multiChainMsg,
      });
    }
  });

  return [
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
          (toAddress.startsWith('0x') ? toAddress : (`0x${toAddress}`)) as Hex,
          BigInt(amount),
        ],
      }), // "0x",
    },
    {
      to: (toAddress.startsWith('0x') ? toAddress : (`0x${toAddress}`)) as Hex,
      value: BigInt(0),
      data: data.startsWith('0x') ? data as Hex : (`0x${data}`) as Hex, // "0x",
    },
  ];
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
        // @ts-ignore
        return await turnkeySigner.signTransaction(tx, signWith);
      } catch (error) {
        throw new Error(`Failed to sign transaction with TurnkeySigner: ${error.message}`);
      }
    },
    signAllTransactions: async (txs: (Transaction | VersionedTransaction)[]) => {
      try {
        // @ts-ignore
        return await turnkeySigner.signAllTransactions(txs, signWith);
      } catch (error) {
        throw new Error(`Failed to sign transactions with TurnkeySigner: ${error.message}`);
      }
    },
  } as Adapter);
}

@Route('bridging')
class BridgeController extends Controller {
  @Post('/sweep')
  async sweep(
    @Query() fromAddress: string,
      @Query() chainId: string,
      // optionally provide the contract and amount, primarily used for solana.
      @Query() amount?: string,
  ): Promise<BridgeResponse> {
    let bridgeFn: (
      fromAddress: string,
      amount: string,
      sourceAssetDenom: string,
      chainId: string,
    ) => Promise<BridgeResponse> | undefined;
    if (chainId === 'solana') {
      bridgeFn = this.startSolanaBridge;
      if (amount) {
        try {
          await bridgeFn(fromAddress, amount, usdcAddressByChainId.solana, chainId);
        } catch (error) {
          logger.error({
            at: `${controllerName}#sweep->startSolanaBridge`,
            message: `Failed to bridge token ${usdcAddressByChainId.solana}`,
            error,
          });
        }
      } else {
        throw new Error('Amount and contract is required for solana');
      }
    } else if (isSupportedEVMChainId(chainId)) {
      bridgeFn = chainId === avalanche.id.toString()
        ? this.startEvmBridgePre7702
        : this.startEvmBridge;
      const alchemy = new Alchemy({
        apiKey: config.ALCHEMY_API_KEY,
        network: chainIdToAlchemyNetworkMap[chainId],
      });

      // search for assets that exist on this account on this chain.
      const usdcToSearch = usdcAddressByChainId[chainId];
      const assets = await alchemy.core.getTokenBalances(fromAddress);

      for (const token of assets.tokenBalances) {
        // TODO: Under what scenario will tokenBalance be undefined?
        if (
          token.contractAddress.toLowerCase() === usdcToSearch.toLowerCase() &&
          token.tokenBalance
        ) {
          // validate that the token balance is not 0.
          if (parseInt(token.tokenBalance, 16) > 0) {
            try {
              await bridgeFn(fromAddress, token.tokenBalance, token.contractAddress, chainId);
            } catch (error) {
              logger.error({
                at: `${controllerName}#sweep->startEvmBridge`,
                message: `Failed to bridge token ${token.contractAddress}`,
                error,
              });
            }
          }
        }
        // TODO: Add other assets here.
      }
    } else {
      throw new Error(`Unsupported chainId: ${chainId}`);
    }

    return {
      toAddress: fromAddress,
      amount: '0',
      sourceAssetDenom: Asset.USDC,
    };
  }

  async startSolanaBridge(
    fromAddress: string,
    amount: string,
    sourceAssetDenom: string,
    chainId: string,
  ): Promise<BridgeResponse> {
    const path = await route({
      goFast: true,
      amountIn: amount,
      sourceAssetDenom,
      sourceAssetChainId: chainId,
      destAssetChainId: 'dydx-mainnet-1',
      cumulativeAffiliateFeeBps: '0',
      destAssetDenom: usdcAddressByChainId['dydx-mainnet-1'],
    });

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
    // Replace with your own private key
    const solanaSponsorPrivateKey = config.SOLANA_SPONSOR_PRIVATE_KEY;
    const sponsorKeypair = Keypair.fromSecretKey(
      bs58.decode(solanaSponsorPrivateKey),
    );
    await executeRoute({
      route: path,
      simulate: false,
      userAddresses,
      getSvmSigner: getSvmSigner(record?.suborg_id || '', fromAddress),
      svmFeePayer: {
        address: sponsorKeypair.publicKey.toString(), // Replace with the fee payer's Solana address
        signTransaction: (dataToSign: Buffer) => {
          const data = new Uint8Array(dataToSign);
          return Promise.resolve(nacl.sign.detached(data, sponsorKeypair.secretKey));
        },
      },
      // ADD RETRY LOGIC???
      // eslint-disable-next-line @typescript-eslint/require-await
      onTransactionBroadcast: async ({ chainId: c, txHash }) => {
        logger.info({
          message: `Broadcasted on ${c}: ${txHash}`,
          from: fromAddress,
          amount,
          sourceAssetDenom,
          chainId,
          toAddress: fromAddress,
          at: new Date().toISOString(),
        });
      },
      // eslint-disable-next-line @typescript-eslint/require-await
      onTransactionCompleted: async ({ chainId: c, txHash, status }) => {
        logger.info({
          at: `${controllerName}#startSolanaBridge`,
          message: 'Transaction completed',
          chainId: c,
          txHash,
          status,
        });
      },
      // eslint-disable-next-line @typescript-eslint/require-await
      onTransactionTracked: async ({ chainId: c, txHash, explorerLink }) => {
        logger.info({
          at: `${controllerName}#startSolanaBridge`,
          message: 'Transaction tracked',
          chainId: c,
          txHash,
          explorerLink,
        });
      },
      // eslint-disable-next-line @typescript-eslint/require-await
      onTransactionSignRequested: async ({ chainId: c, signerAddress }) => {
        logger.info({
          at: `${controllerName}#startSolanaBridge`,
          message: 'Sign requested',
          chainId: c,
          signerAddress,
        });
      },
      // eslint-disable-next-line @typescript-eslint/require-await
      onValidateGasBalance: async ({ chainId: c, txIndex, status }) => {
        logger.info({
          at: `${controllerName}#startSolanaBridge`,
          message: 'Gas balance validation',
          chainId: c,
          txIndex,
          status,
        });
      },
    });

    return {
      toAddress: fromAddress,
      amount,
      sourceAssetDenom,
    };
  }

  /*
   * This function is used to bridge evm assets that are supported by the skip bridge.
   * This is done by using the eip7702 kernel account with zero dev to sponsor the user operation.
   * This function is used for all evm chains that support eip7702 including arbitrum, ethereum,
   * base, and optimism.
  */
  async startEvmBridge(
    fromAddress: string,
    amount: string,
    sourceAssetDenom: string,
    chainId: string,
  ): Promise<BridgeResponse> {
    const record: TurnkeyUserFromDatabase | undefined = await findByEvmAddress(fromAddress);
    if (!record || !record.dydx_address) {
      throw new Error('Failed to derive dYdX address');
    }
    let callData: Parameters<SmartAccountImplementation['encodeCalls']>[0] = [];
    try {
      callData = await getSkipCallData(
        fromAddress,
        sourceAssetDenom,
        record.dydx_address,
        amount,
        chainId,
      );
    } catch (error) {
      logger.error({
        at: `${controllerName}#startEvmBridge`,
        message: 'Failed to get Skip call data',
        error,
      });
      throw error;
    }

    const entryPoint = getEntryPoint('0.7');

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
    });

    const zerodevPaymaster = createZeroDevPaymasterClient({
      chain: chains[chainId],
      transport: http(getRPCEndpoint(chainId)),
    });

    const kernelClient = createKernelAccountClient({
      account,
      chain: chains[chainId],
      client: publicClients[chainId],
      bundlerTransport: http(getRPCEndpoint(chainId)),
      paymaster: {
        getPaymasterData: async (userOperation) => {
          return zerodevPaymaster.sponsorUserOperation({ userOperation });
        },

      },
      // paymasterContext: { token: assetAddressLookerUpper[asset as Asset][chainId] },
      userOperation: {
        estimateFeesPerGas: async ({ bundlerClient }) => {
          return getUserOperationGasPrice(bundlerClient);
        },
      },
    });

    const userOpHash = await kernelClient.sendUserOperation({
      callData: await kernelClient.account.encodeCalls(callData),
    });
    const { receipt } = await kernelClient.waitForUserOperationReceipt({
      hash: userOpHash,
    });
    logger.info({
      at: `${controllerName}#startEvmBridge`,
      message: 'UserOp completed',
      transactionHash: receipt.transactionHash,
    });
    return {
      toAddress: fromAddress,
      amount,
      sourceAssetDenom,
    };
  }

  /*
   * This function is only used for avalanche as they do not support the eip7702 yet.
   * Similar logic to startEvmBridge with the distinction that the smart account is a different
   * address as the EOA account.
   *
   * We will only auto bridge funds sent to the smart account address for avalance pre 7702 because
   * no gas sponsorship is possible pre 7702. We are assuming that the address provided will be a
   * smart account address and that the underlying EOA address is a valid entry in our database.
   *
   */
  async startEvmBridgePre7702(
    fromAddress: string,
    amount: string,
    sourceAssetDenom: string,
    chainId: string,
  ): Promise<BridgeResponse> {
    const eoaAddress = getEOAAddressFromSmartAccountAddress(fromAddress as Address);
    const record: TurnkeyUserFromDatabase | undefined = await findByEvmAddress(eoaAddress);
    if (!record || !record.dydx_address) {
      throw new Error('Failed to derive dYdX address');
    }
    let callData: Parameters<SmartAccountImplementation['encodeCalls']>[0] = [];
    try {
      callData = await getSkipCallData(
        eoaAddress,
        sourceAssetDenom,
        record.dydx_address,
        amount,
        chainId,
      );
    } catch (error) {
      logger.error({
        at: `${controllerName}#startEvmBridgePre7702`,
        message: 'Failed to get Skip call data',
        error,
      });
      throw error;
    }

    const entryPoint = getEntryPoint('0.7');

    // Initialize a Turnkey-powered Viem Account
    // needs to sign with eoa address.
    const turnkeyAccount = await createAccount({
      // @ts-ignore
      client: turnkeySenderClient.apiClient(),
      organizationId: record.suborg_id,
      signWith: eoaAddress,
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
    });
    logger.info({
      at: `${controllerName}#startEvmBridgePre7702`,
      message: 'Account created',
      accountAddress: account.address,
    });

    const zerodevPaymaster = createZeroDevPaymasterClient({
      chain: chains[chainId],
      transport: http(getRPCEndpoint(chainId)),
    });

    const kernelClient = createKernelAccountClient({
      account,
      chain: chains[chainId],
      client: publicClients[chainId],
      bundlerTransport: http(getRPCEndpoint(chainId)),
      paymaster: {
        getPaymasterData: async (userOperation) => {
          return zerodevPaymaster.sponsorUserOperation({ userOperation });
        },

      },
      // paymasterContext: { token: assetAddressLookerUpper[asset as Asset][chainId] },
      userOperation: {
        estimateFeesPerGas: async ({ bundlerClient }) => {
          return getUserOperationGasPrice(bundlerClient);
        },
      },
    });

    try {
      const userOpHash = await kernelClient.sendUserOperation({
        callData: await kernelClient.account.encodeCalls(callData),
      });
      logger.info({
        at: `${controllerName}#startEvmBridgePre7702`,
        message: 'UserOp sent',
        userOpHash,
      });

      const { receipt } = await kernelClient.waitForUserOperationReceipt({
        hash: userOpHash,
      });
      logger.info({
        at: `${controllerName}#startEvmBridgePre7702`,
        message: 'UserOp completed',
        transactionHash: receipt.transactionHash,
      });
    } catch (error) {
      logger.error({
        at: `${controllerName}#startEvmBridgePre7702`,
        message: 'Failed to send user operation, AVAX does not support eip7702 yet, did you remember to send to the smart account?',
        error,
      });
      throw new Error(`Failed to send user operation, AVAX does not support eip7702 yet, did you remember to send to the smart account?, error: ${error}`);
    }
    return {
      toAddress: fromAddress,
      amount,
      sourceAssetDenom,
    };
  }
}

function getEOAAddressFromSmartAccountAddress(_: Address): Address {
  return '0x0001';
}

/* returns the addresses to sweep and the chainId.
 * for solana sweeps, an amount is also included in the map.
 *     this amount is the usdc amount that is being swept.
 * for evm sweeps, an amount is not included in the map.
 */
async function parseEvent(e: express.Request): Promise<{
  addressesToSweep: Map<string, string>,
  chainId: string,
}> {
  const { event: { transaction, activity, network } } = e.body;
  let chainId = '';
  chainId = alchemyNetworkToChainIdMap[network];
  const addressesToProcess = new Map<string, string>();
  // for evm parsing only
  if (activity) {
    for (const act of activity) {
      const fromAddress = act.toAddress;
      addressesToProcess.set(fromAddress, '');
    }
  }
  // for solana parsing only.
  if (transaction) {
    for (const tx of transaction) {
      if (!tx.meta) {
        continue;
      }
      for (const meta of tx.meta || []) {
        for (const postTokenBalance of meta.post_token_balances || []) {
          // usdc is the only supported asset for solana so we only include the usdc amount
          // and ignore the rest and we only sweep if there is a positive amount.
          if (
            postTokenBalance.owner &&
            postTokenBalance.mint === usdcAddressByChainId.solana &&
            postTokenBalance.ui_token_amount.ui_amount > 0
          ) {
            addressesToProcess.set(postTokenBalance.owner, postTokenBalance.ui_token_amount.amount);
          }
        }
      }
    }
  }
  // validate the addressesToProcess to see if they are indeed turnkey users.
  const addressesToSweep = new Map<string, string>();
  for (const fromAddress of addressesToProcess.keys()) {
    // if the chain is solana, then we need to also include the token amount.
    // USDC is the only supported asset for solana so that will be the contract address.
    if (chainId === 'solana') {
      const record: TurnkeyUserFromDatabase | undefined = await findBySvmAddress(fromAddress);
      if (!record || !record.dydx_address) {
        logger.warning({
          at: `${controllerName}#parseEvent`,
          message: 'Failed to find a turnkey user for address',
          address: fromAddress,
        });
        continue;
      }
      logger.info({
        at: `${controllerName}#parseEvent`,
        message: 'Found a turnkey user for address',
        address: fromAddress,
      });
      // add the amount to the map as well.
      addressesToSweep.set(fromAddress, addressesToProcess.get(fromAddress) || '');
    } else {
      // need the checksummed address to find the turnkey user.
      const checkSummedFromAddress = checksumAddress(fromAddress as `0x${string}`);
      const record = await findByEvmAddress(checkSummedFromAddress);
      if (!record || !record.dydx_address) {
        logger.warning({
          at: `${controllerName}#parseEvent`,
          message: 'Failed to find a turnkey user for address',
          address: checkSummedFromAddress,
        });
        continue;
      }
      // no amount required for evm.
      addressesToSweep.set(checkSummedFromAddress, '');
    }
  }
  return {
    addressesToSweep,
    chainId,
  };
}

router.post(
  '/startBridge',
  ...CheckBridgeSchema,
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    try {
      const bridgeController = new BridgeController();
      const { addressesToSweep, chainId } = await parseEvent(req);
      // Iterate over the set 'toProcess' and process each item
      for (const fromAddress of addressesToSweep.keys()) {
        await bridgeController.sweep(
          fromAddress,
          chainId,
          addressesToSweep.get(fromAddress) === '' ? undefined : addressesToSweep.get(fromAddress),
        );
      }
      return res.status(200).send();
    } catch (error) {
      return handleControllerError(
        'BridgeController POST /startBridge',
        'Bridge start error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_addresses.timing`,
        Date.now() - start,
      );
    }
  },
);

export default router;
