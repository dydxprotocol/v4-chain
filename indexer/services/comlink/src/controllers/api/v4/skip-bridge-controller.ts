import { logger, stats } from '@dydxprotocol-indexer/base';
import { TurnkeyUsersTable, dbHelpers } from '@dydxprotocol-indexer/postgres';
import {
  findByEvmAddress, findBySmartAccountAddress, findBySvmAddress, findByDydxAddress,
} from '@dydxprotocol-indexer/postgres/build/src/stores/turnkey-users-table';
import { TurnkeyUserFromDatabase } from '@dydxprotocol-indexer/postgres/build/src/types';
import {
  route, executeRoute, setClientOptions, balances, TransferStatus,
} from '@skip-go/client/cjs';
import { Turnkey } from '@turnkey/sdk-server';
import { TurnkeySigner } from '@turnkey/solana';
import { createAccount } from '@turnkey/viem';
import { signerToEcdsaValidator } from '@zerodev/ecdsa-validator';
import { deserializePermissionAccount } from '@zerodev/permissions';
import { toECDSASigner } from '@zerodev/permissions/signers';
import {
  createKernelAccount, createKernelAccountClient,
  CreateKernelAccountReturnType,
  createZeroDevPaymasterClient, getUserOperationGasPrice,
} from '@zerodev/sdk';
import { KERNEL_V3_3, KERNEL_V3_1 } from '@zerodev/sdk/constants';
import express from 'express';
import {
  Controller, Post, Query, Route,
} from 'tsoa';
import {
  Address,
  checksumAddress, http,
} from 'viem';
import {
  type EntryPointVersion,
  type SmartAccountImplementation,
} from 'viem/account-abstraction';
import { privateKeyToAccount } from 'viem/accounts';
import {
  arbitrum, avalanche,
} from 'viem/chains';

import config from '../../../config';
import {
  chains,
  getEOAAddressFromSmartAccountAddress,
  isSupportedEVMChainId,
  getRPCEndpoint,
  getAddress,
  publicClients,
  alchemyNetworkToChainIdMap,
  ethDenomByChainId,
} from '../../../helpers/alchemy-helpers';
import {
  getSvmSigner, getSkipCallData, suborgToApproval,
  nobleToSolana,
} from '../../../helpers/skip-helper';
import { handleControllerError } from '../../../lib/helpers';
import { entryPoint, usdcAddressByChainId } from '../../../lib/smart-contract-constants';
import { CheckBridgeSchema, CheckGetDepositAddressSchema } from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { PolicyEngine } from './policy-controller';

const router = express.Router();
const controllerName: string = 'bridging-controller';

// set the skip client options to use the skip rpc.
// for some reason, skip requires you to call this function
// as initiation even if you're not setting an rpc.
// Calling route() without this will throw error.
setClientOptions({
  endpointOptions: {
    endpoints: {
      solana: {
        rpc: 'https://go.skip.build/api/rpc/solana',
      },
    },
  },
});

// need to add this so that the address that triggered the activity is ignored.
// const processingCache: Record<string, boolean> = {};

const turnkeySenderClient = new Turnkey({
  apiBaseUrl: config.TURNKEY_API_BASE_URL as string,
  apiPublicKey: config.TURNKEY_API_SENDER_PUBLIC_KEY as string,
  apiPrivateKey: config.TURNKEY_API_SENDER_PRIVATE_KEY as string,
  defaultOrganizationId: config.TURNKEY_ORGANIZATION_ID,
});

const turnkeyGasClient = new Turnkey({
  apiBaseUrl: config.TURNKEY_API_BASE_URL as string,
  apiPublicKey: config.TURNKEY_API_PUBLIC_KEY as string,
  apiPrivateKey: config.TURNKEY_API_PRIVATE_KEY as string,
  defaultOrganizationId: config.TURNKEY_ORGANIZATION_ID,
});

const turnkeyGasFeePayer = new TurnkeySigner({
  organizationId: config.TURNKEY_ORGANIZATION_ID,
  client: turnkeyGasClient.apiClient(),
});

interface BridgeResponse {
  success: boolean,
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

@Route('bridging')
class BridgeController extends Controller {
  @Post('/sweep')
  async sweep(
    @Query() fromAddress: string,
    @Query() chainId: string,
    // optionally provide the contract and amount, primarily used for solana.
    @Query() amount?: string,
  ): Promise<BridgeResponse> {
    if (chainId === 'solana') {
      // no sweeping if amount is less than 20, so far we only support usdc on svm.
      // so no need to check for the dollar value here.
      if (amount && parseInt(amount, 10) >= config.BRIDGE_THRESHOLD_USDC) {
        try {
          await this.startSolanaBridge(fromAddress, amount, usdcAddressByChainId.solana, chainId);
        } catch (error) {
          logger.error({
            at: `${controllerName}#sweep->startSolanaBridge`,
            message: `Failed to bridge token ${usdcAddressByChainId.solana}`,
            error,
          });
        }
      } else if (!amount) {
        throw new Error('Amount is required for solana');
      } else {
        throw new Error(`Amount must be greater than ${config.BRIDGE_THRESHOLD_USDC} to start auto bridge`);
      }
    } else if (isSupportedEVMChainId(chainId)) {
      // search for assets that exist on this account on this chain.
      const usdcToSearch = usdcAddressByChainId[chainId];
      const ethToSearch = ethDenomByChainId[chainId];

      const assetsToSearch = [usdcToSearch, ethToSearch];
      const assets = await balances({
        chains: {
          [chainId]: {
            address: fromAddress,
          },
        },
      });
      logger.info({
        at: `${controllerName}#sweep->startEvmBridge`,
        message: 'Assets found',
        assets,
      });

      for (let asset of assetsToSearch) {
        asset = (asset && asset.startsWith('0x')) ? checksumAddress(asset as Address) : asset;
        const balance = assets?.chains?.[chainId]?.denoms?.[asset]?.amount;
        if (balance && parseInt(balance, 10) > 0) {
          logger.info({
            at: `${controllerName}#sweep->startEvmBridge`,
            message: 'Bridge token',
            fromAddress,
            chainId,
            asset,
            balance,
          });
          try {
            await this.startEvmBridge(fromAddress, balance, asset, chainId);
          } catch (error) {
            logger.error({
              at: `${controllerName}#sweep->startEvmBridge`,
              message: `Failed to bridge token ${asset}`,
              error,
            });
          }
        }
      }
    } else {
      throw new Error(`Unsupported chainId: ${chainId}`);
    }

    return {
      success: true,
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
      allowUnsafe: false,
      smartRelay: true, // skip recommended to enable for better routes and less faults.
      smartSwapOptions: {
        splitRoutes: true,
      },
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
    if (!record) {
      throw new Error(`Failed to find a turnkey user for svm address: ${fromAddress}`);
    }
    // Replace with your own private key
    const solanaSponsorPublicKey = config.SOLANA_SPONSOR_PUBLIC_KEY;
    if (!solanaSponsorPublicKey) {
      throw new Error(
        'Missing required environment variable: SOLANA_SPONSOR_PRIVATE_KEY',
      );
    }
    await executeRoute({
      route: path,
      simulate: false, // turned off for now, otherwise simulation will fail due to race.
      userAddresses,
      getSvmSigner: getSvmSigner(record?.suborg_id || '', fromAddress),
      svmFeePayer: {
        address: solanaSponsorPublicKey, // Replace with the fee payer's Solana address
        signTransaction: (dataToSign: Buffer) => {
          const data = new Uint8Array(dataToSign);
          return turnkeyGasFeePayer.signMessage(data, solanaSponsorPublicKey);
        },
      },
      onTransactionBroadcast: async (
        { chainId: c, txHash }: {
          chainId: string,
          txHash: string,
        },
        // eslint-disable-next-line @typescript-eslint/require-await
      ) => {
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
      onTransactionCompleted: async (
        { chainId: c, txHash, status }: {
          chainId: string,
          txHash: string,
          status?: TransferStatus,
        },
        // eslint-disable-next-line @typescript-eslint/require-await
      ) => {
        logger.info({
          at: `${controllerName}#startSolanaBridge`,
          message: 'Transaction completed',
          chainId: c,
          txHash,
          status,
        });
      },
      onTransactionTracked: async (
        { chainId: c, txHash, explorerLink }: {
          chainId: string,
          txHash: string,
          explorerLink: string,
        },
        // eslint-disable-next-line @typescript-eslint/require-await
      ) => {
        logger.info({
          at: `${controllerName}#startSolanaBridge`,
          message: 'Transaction tracked',
          chainId: c,
          txHash,
          explorerLink,
        });
      },
      onTransactionSignRequested: async (
        { chainId: c, txIndex, signerAddress }: {
          chainId: string,
          txIndex: number,
          signerAddress?: string,
        },
        // eslint-disable-next-line @typescript-eslint/require-await
      ) => {
        logger.info({
          at: `${controllerName}#startSolanaBridge`,
          message: 'Sign requested',
          chainId: c,
          signerAddress,
          txIndex,
        });
      },
    });

    return {
      success: true,
    };
  }

  async getKernelAccount(
    chainId: string,
    fromAddress: string,
    suborgId: string,
  ): Promise<CreateKernelAccountReturnType<EntryPointVersion>> {
    // Initialize a Turnkey-powered Viem Account
    // needs to sign with eoa address.
    const turnkeyAccount = await createAccount({
      // @ts-ignore
      client: turnkeySenderClient.apiClient(),
      organizationId: suborgId,
      signWith: fromAddress,
    });
    // use the permissioned master key as a signer.
    const privateKeyAccount = privateKeyToAccount(config.MASTER_SIGNER_PRIVATE as `0x${string}`);
    const sessionKeySigner = await toECDSASigner({
      signer: privateKeyAccount,
    });
    if (chainId === arbitrum.id.toString()) {
      const sessionKeyAccount = await deserializePermissionAccount(
        publicClients[chainId],
        entryPoint,
        KERNEL_V3_3,
        suborgToApproval.get(suborgId) || '',
        sessionKeySigner,
      );
      return sessionKeyAccount;
    }
    if (chainId === avalanche.id.toString()) {
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
      return account;
    }
    return createKernelAccount(publicClients[chainId], {
      entryPoint,
      eip7702Account: turnkeyAccount,
      kernelVersion: KERNEL_V3_3,
    });
  }

  /*
   * This function is used to bridge evm assets that are supported by the skip bridge.
   * This is done by using the eip7702 kernel account with zero dev to sponsor the user operation
   * for ethereum, arbitrum, base, optimism.
   *
   * Avalanche does not support eip 7702 so we need a separate smart contract account from the EOA.
   * Then this account is sponsored by zerodev for gas. Right now, this is being done by swapping
   * the smart account address for the EOA address and carrying on with the sponsorship from there.
  */
  async startEvmBridge(
    fromAddress: string,
    amount: string,
    sourceAssetDenom: string,
    chainId: string,
  ): Promise<BridgeResponse> {
    const pre7702 = chainId === avalanche.id.toString();
    let srcAddress = fromAddress;
    if (pre7702) {
      // need to swap the smart account address to find what the signing address is for
      // avalanche because avalanche does not support pectra (eip7702) yet.
      try {
        srcAddress = await getEOAAddressFromSmartAccountAddress(srcAddress as Address);
      } catch (error) {
        logger.error({
          at: `${controllerName}#startEvmBridgePre7702`,
          message: 'Failed to get EOA address from smart account address',
          error,
        });
        throw new Error(`Cannot proceed with bridge, cannot get EOA address from smart account address : ${error.message}`);
      }
    }
    const record: TurnkeyUserFromDatabase | undefined = await findByEvmAddress(srcAddress);
    if (!record || !record.dydx_address) {
      throw new Error('Failed to derive dYdX address');
    }
    let callData: Parameters<SmartAccountImplementation['encodeCalls']>[0] = [];
    try {
      callData = await getSkipCallData(
        srcAddress,
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
    let account: CreateKernelAccountReturnType<EntryPointVersion>;
    try {
      account = await this.getKernelAccount(chainId, srcAddress, record.suborg_id);
    } catch (error) {
      logger.error({
        at: `${controllerName}#startEvmBridge`,
        message: 'Failed to get kernel account',
        error,
      });
      throw error;
    }

    const zerodevPaymaster = createZeroDevPaymasterClient({
      chain: chains[chainId],
      transport: http(getRPCEndpoint(chainId)),
    });

    const kernelClient = createKernelAccountClient({
      account,
      chain: chains[chainId],
      client: publicClients[chainId],
      bundlerTransport: http(getRPCEndpoint(chainId)),
      paymaster: zerodevPaymaster,
      userOperation: {
        estimateFeesPerGas: async ({ bundlerClient }) => {
          return getUserOperationGasPrice(bundlerClient);
        },
      },
    });

    try {
      const encoded = await kernelClient.account.encodeCalls(callData);
      const userOpHash = await kernelClient.sendUserOperation({
        callData: encoded,
      });
      logger.info({
        at: `${controllerName}#startEvmBridge`,
        message: 'UserOp sent',
        userOpHash,
      });
      const { receipt } = await kernelClient.waitForUserOperationReceipt({
        hash: userOpHash,
      });
      logger.info({
        at: `${controllerName}#startEvmBridge`,
        message: 'UserOp completed',
        transactionHash: receipt.transactionHash,
      });
    } catch (error) {
      logger.error({
        at: `${controllerName}#startEvmBridge`,
        message: 'Failed to send user operation',
        callData,
        error,
      });
      throw new Error(`Failed to send user operation, error: ${error}`);
    }
    return {
      success: true,
    };
  }

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
  if (!chainId) {
    throw new Error(`Unsupported network: ${network}`);
  }
  const addressesToProcess = new Map<string, string>();
  // for evm parsing only
  if (activity) {
    for (const act of activity) {
      const bridgeOriginAddress = act.toAddress;
      addressesToProcess.set(bridgeOriginAddress, '');
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
  for (const bridgeOriginAddress of addressesToProcess.keys()) {
    // if the chain is solana, then we need to also include the token amount.
    // USDC is the only supported asset for solana so that will be the contract address.
    if (chainId === 'solana') {
      const record = await findBySvmAddress(bridgeOriginAddress);
      if (!record || !record.dydx_address) {
        logger.warning({
          at: `${controllerName}#parseEvent`,
          message: 'Failed to find a turnkey user for address',
          address: bridgeOriginAddress,
        });
        continue;
      }
      logger.info({
        at: `${controllerName}#parseEvent`,
        message: 'Found a turnkey user for address',
        address: bridgeOriginAddress,
      });
      // add the amount to the map as well.
      addressesToSweep.set(bridgeOriginAddress, addressesToProcess.get(bridgeOriginAddress) || '');
    } else {
      // evm otherwise, check to see if the chain is avalanche, in which case
      // we need to use the underlying eoa address to find the turnkey user.
      // this is also the address we need to use to kick off the bridge, but
      // this address is hot swapped on the actual bridging fn because we still need
      // the smart account address for amount validation.
      const checkSummedFromAddress = checksumAddress(bridgeOriginAddress as Address);
      let record: TurnkeyUserFromDatabase | undefined;
      if (chainId === avalanche.id.toString()) {
        record = await findBySmartAccountAddress(checkSummedFromAddress);
      } else {
        record = await findByEvmAddress(checkSummedFromAddress);
      }
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

router.get(
  '/getDepositAddress/:dydxAddress',
  ...CheckGetDepositAddressSchema,
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    try {
      const { dydxAddress } = req.params;
      const record: TurnkeyUserFromDatabase | undefined = await findByDydxAddress(dydxAddress);

      if (!record) {
        return res.status(404).json({
          error: 'User not found',
          message: `No user found with dydx address: ${dydxAddress}`,
        });
      }

      return res.status(200).json({
        evmAddress: record.evm_address,
        avalancheAddress: record.smart_account_address,
        svmAddress: record.svm_address,
      });
    } catch (error) {
      return handleControllerError(
        'BridgeController GET /getDepositAddress',
        'Get deposit address error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_deposit_address.timing`,
        Date.now() - start,
      );
    }
  },
);

router.post(
  '/startBridge',
  ...CheckBridgeSchema,
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    // await dbHelpers.clearData();
    // await TurnkeyUsersTable.create({
    //   suborg_id: 'b717b91e-61a4-4707-ba16-a6635ee14143',
    //   svm_address: 'EWHuLzjV2PMR2w64RLMyEW7tPA7sWHPXRSWbnf2pEscE',
    //   evm_address: '0x74e7A23338D294b14Fc819Be0a179ed9E2a26ca1',
    //   smart_account_address: '0xd2A6baf165CF630B39A74ad2Ef1b5A917f74ABE0',
    //   salt: '112dca5a557c8f0f103cd88ad32c178e5bc1bd5e62cbaa1b5936d01a4538bc80',
    //   dydx_address: 'dydx1sjssdnatk99j2sdkqgqv55a8zs97fcvstzreex',
    //   created_at: new Date().toISOString(),
    // })
    const start: number = Date.now();
    try {
      const bridgeController = new BridgeController();
      const { addressesToSweep, chainId } = await parseEvent(req);
      // Iterate over the set 'toProcess' and process each item
      for (const fromAddress of addressesToSweep.keys()) {
        await bridgeController.sweep(
          fromAddress,
          chainId,
          addressesToSweep.get(fromAddress) === '' ? undefined : addressesToSweep.get(fromAddress), // amount
        );
      }

      // sending a 200 tells webhook to not retry.
      return res.status(200).send();
    } catch (error) {
      // will trigger retry with exponential backoff.
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
