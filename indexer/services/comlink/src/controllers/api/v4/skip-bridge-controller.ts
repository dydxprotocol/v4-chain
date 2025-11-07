import {
  createHash,
} from 'crypto';

import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  BridgeInformationTable,
  TurnkeyUsersTable,
  TurnkeyUserFromDatabase,
  BridgeInformationCreateObject,
  IsoString,
} from '@dydxprotocol-indexer/postgres';
import {
  route, executeRoute, setClientOptions, balances, TransferStatus,
} from '@skip-go/client/cjs';
import { Turnkey } from '@turnkey/sdk-server';
import { TurnkeySigner } from '@turnkey/solana';
import {
  createKernelAccountClient,
  CreateKernelAccountReturnType,
  createZeroDevPaymasterClient, getUserOperationGasPrice,
} from '@zerodev/sdk';
import express from 'express';
import { matchedData } from 'express-validator';
import {
  Controller, Route,
} from 'tsoa';
import {
  Address,
  checksumAddress, http,
} from 'viem';
import {
  type EntryPointVersion,
  type SmartAccountImplementation,
} from 'viem/account-abstraction';
import {
  avalanche,
} from 'viem/chains';

import config from '../../../config';
import {
  chains,
  getEOAAddressFromSmartAccountAddress,
  isSupportedEVMChainId,
  getZeroDevRPCEndpoint,
  publicClients,
  alchemyNetworkToChainIdMap,
} from '../../../helpers/alchemy-helpers';
import {
  getSvmSigner, getSkipCallData, getKernelAccount,
  buildUserAddresses,
  limitAmount,
} from '../../../helpers/skip-helper';
import { trackTurnkeyDepositSubmitted } from '../../../lib/amplitude-helpers';
import { handleControllerError } from '../../../lib/helpers';
import {
  dydxChainId, usdcAddressByChainId, ethDenomByChainId,
  SOLANA_USDC_QUANTUM,
  ETH_USDC_QUANTUM,
} from '../../../lib/smart-contract-constants';
import {
  CheckBridgeSchema,
  CheckGetDepositAddressSchema,
  CheckLimitAndCreatedBeforeOrAtAndOnOrAfterSchema, CheckPaginationSchema,
} from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';

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
  message?: string,
}

function formatError(error: unknown): string {
  if (error instanceof Error) {
    return error.stack || error.message;
  }
  try {
    return JSON.stringify(error);
  } catch (_) {
    return String(error);
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
    const record = await TurnkeyUsersTable.findByEvmAddress(address);
    dydxAddress = record?.dydx_address || '';
  } else if (chainId === 'solana') {
    // look up in turnkey table
    logger.info({
      at: `${controllerName}#getDydxAddress`,
      message: 'Looking up in turnkey table for svm address',
      address,
    });
    const record = await TurnkeyUsersTable.findBySvmAddress(address);
    dydxAddress = record?.dydx_address || '';
  } else {
    throw new Error(`Unsupported chainId: ${chainId}`);
  }
  return dydxAddress;
}

@Route('bridging')
class BridgeController extends Controller {
  async sweep(
    fromAddress: string,
    chainId: string,
    // optionally provide the contract and amount, primarily used for solana.
    amount?: string,
  ): Promise<BridgeResponse> {
    if (chainId === 'solana') {
      if (!amount) {
        throw new Error('Amount is required for solana');
      }
      // get usdc amount and check that the usd amount is greater than the threshold.
      const usdAmt = BigInt(amount) / BigInt(SOLANA_USDC_QUANTUM);
      if (usdAmt && (usdAmt >= config.BRIDGE_THRESHOLD_USDC)) {
        try {
          await this.startSolanaBridge(fromAddress, amount, usdcAddressByChainId.solana, chainId);
        } catch (error) {
          logger.error({
            at: `${controllerName}#sweep->startSolanaBridge`,
            message: `Failed to bridge token ${usdcAddressByChainId.solana}`,
            error,
          });
        }
      } else {
        logger.info({
          at: `${controllerName}#sweep->startSolanaBridge`,
          message: 'Amount is less than threshold, skipping bridge',
          address: fromAddress,
          chainId,
          usdAmt,
          amount,
          threshold: config.BRIDGE_THRESHOLD_USDC,
        });
        return {
          success: false,
        };
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

      for (const asset of assetsToSearch) {
        const balance = assets?.chains?.[chainId]?.denoms?.[asset]?.amount;
        if (!balance) {
          logger.info({
            at: `${controllerName}#sweep->startEvmBridge`,
            message: 'Balance is not found, skipping bridge',
            address: fromAddress,
            chainId,
            asset,
          });
          continue;
        }
        let usdAmount: bigint;
        if (asset === usdcToSearch) {
          usdAmount = BigInt(balance) / BigInt(ETH_USDC_QUANTUM);
        } else {
          const valueUsd = assets?.chains?.[chainId]?.denoms?.[asset]?.valueUsd || '0';
          usdAmount = BigInt(Math.floor(parseFloat(valueUsd)));
        }
        // To sweep and asset, user needs to have at least BRIDGE_THRESHOLD_USDC in it.
        if (balance && parseInt(balance, 10) > 0 &&
          usdAmount && usdAmount >= BigInt(config.BRIDGE_THRESHOLD_USDC)) {
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
        } else {
          logger.info({
            at: `${controllerName}#sweep->startEvmBridge`,
            message: 'Amount is less than threshold, skipping bridge',
            address: fromAddress,
            chainId,
            asset,
            balance,
            usdAmount,
            threshold: config.BRIDGE_THRESHOLD_USDC,
          });
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
    const dydxAddress = await getDydxAddress(fromAddress, chainId);
    if (!dydxAddress) {
      throw new Error('Failed to derive dYdX address');
    }

    const amountIn = await limitAmount(chainId, amount, sourceAssetDenom);

    const path = await route({
      goFast: true,
      allowUnsafe: false,
      smartRelay: true, // skip recommended to enable for better routes and less faults.
      smartSwapOptions: {
        splitRoutes: true,
      },
      amountIn,
      sourceAssetDenom,
      sourceAssetChainId: chainId,
      destAssetChainId: dydxChainId,
      cumulativeAffiliateFeeBps: '0',
      destAssetDenom: usdcAddressByChainId[dydxChainId],
    });
    if (!path) {
      throw new Error('Failed to create route');
    }

    // the end user must be able to sign these intermediate addresses as
    // they are the addresses that the funds will be deposited into in case
    // of a failure.
    const userAddresses = await buildUserAddresses(
      path.requiredChainAddresses, fromAddress, dydxAddress,
    );

    // find the suborgId for the user
    const record: TurnkeyUserFromDatabase | undefined = await TurnkeyUsersTable.findBySvmAddress(
      fromAddress,
    );
    if (!record) {
      throw new Error(`Failed to find a turnkey user for svm address: ${fromAddress}`);
    }

    const solanaSponsorPublicKey = config.SOLANA_SPONSOR_PUBLIC_KEY;
    if (!solanaSponsorPublicKey) {
      throw new Error(
        'Missing required environment variable: SOLANA_SPONSOR_PUBLIC_KEY',
      );
    }
    await executeRoute({
      route: path,
      simulate: false, // turned off for now, otherwise simulation will fail due to race.
      slippageTolerancePercent: config.SKIP_SLIPPAGE_TOLERANCE_PERCENTAGE,
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
          amount: amountIn,
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
          message: 'Bridge transaction completed',
          fromAddress,
          chainId: c,
          amount: amountIn,
          sourceAssetDenom,
          transactionHash: txHash,
          status,
          completedAt: new Date().toISOString(),
        });
      },
      onTransactionTracked: async (
        { chainId: c, txHash, explorerLink }: {
          chainId: string,
          txHash: string,
          explorerLink: string,
        },
      ) => {
        try {
          const bridgeRecord = {
            from_address: fromAddress,
            chain_id: c,
            amount: amountIn,
            transaction_hash: txHash,
            created_at: new Date().toISOString(),
          };

          await BridgeInformationTable.create(bridgeRecord);
          const email = record.email?.trim().toLowerCase();
          // sha256 hash email
          const emailHash = email ? createHash('sha256').update(email).digest('hex') : record.evm_address;
          // Track TurnKey deposit confirmation event in Amplitude
          await trackTurnkeyDepositSubmitted(
            emailHash,
            c,
            amountIn,
            txHash,
            sourceAssetDenom,
          );
          logger.info({
            at: `${controllerName}#startSolanaBridge`,
            message: 'Bridge transaction tracked',
            fromAddress,
            chainId: c,
            amount: amountIn,
            sourceAssetDenom,
            transactionHash: txHash,
            explorerLink,
            trackedAt: new Date().toISOString(),
          });
        } catch (error) {
          logger.error({
            at: `${controllerName}#startSolanaBridge`,
            message: 'Failed to create bridge information record on tracked',
            fromAddress,
            chainId: c,
            amount: amountIn,
            error: error.message || error,
          });
          // Don't throw error to avoid breaking the bridge flow
        }
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
        throw new Error(
          `Cannot proceed with bridge, cannot get EOA address from smart account address : ${formatError(error)}`,
        );
      }
    }
    const record: TurnkeyUserFromDatabase | undefined = await TurnkeyUsersTable.findByEvmAddress(
      srcAddress,
    );
    if (!record || !record.dydx_address) {
      throw new Error('Failed to derive dYdX address');
    }
    // we cannot bridge more than the max amount allowed through the bridge.
    const amountToUse = await limitAmount(chainId, amount, sourceAssetDenom);
    let callData: Parameters<SmartAccountImplementation['encodeCalls']>[0] = [];
    try {
      callData = await getSkipCallData(
        srcAddress,
        sourceAssetDenom,
        record.dydx_address,
        amountToUse,
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

    // kernel account is required here. This point onwards is just calling zerodev api to carry out
    // the user op on chain.
    let account: CreateKernelAccountReturnType<EntryPointVersion>;
    try {
      account = await getKernelAccount(chainId, srcAddress, record.suborg_id);
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
      transport: http(getZeroDevRPCEndpoint(chainId)),
    });

    const kernelClient = createKernelAccountClient({
      account,
      chain: chains[chainId],
      client: publicClients[chainId],
      bundlerTransport: http(getZeroDevRPCEndpoint(chainId)),
      paymaster: zerodevPaymaster,
      userOperation: {
        estimateFeesPerGas: async ({ bundlerClient }) => {
          return getUserOperationGasPrice(bundlerClient);
        },
      },
    });

    // sending the userop to chain via zerodev kernel.
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
      // track the transaction hash in the bridge information table
      const { receipt } = await kernelClient.waitForUserOperationReceipt({
        hash: userOpHash,
      });
      logger.info({
        at: `${controllerName}#startEvmBridge`,
        message: 'UserOp completed',
        transactionHash: receipt.transactionHash,
      });

      // Update bridge information record with transaction hash
      try {
        // Create bridge information record before sending the transaction
        const bridgeRecord: BridgeInformationCreateObject = {
          from_address: fromAddress,
          chain_id: chainId,
          amount: amountToUse,
          transaction_hash: receipt.transactionHash,
          created_at: new Date().toISOString(),
        };
        await BridgeInformationTable.create(
          bridgeRecord,
        );

        // Track TurnKey deposit confirmation event in Amplitude
        if (record.email) {
          const email = record.email.trim().toLowerCase();
          // sha256 hash email
          const emailHash = createHash('sha256').update(email).digest('hex');
          await trackTurnkeyDepositSubmitted(
            emailHash,
            chainId,
            amountToUse,
            receipt.transactionHash,
            sourceAssetDenom,
          );
        } else {
          await trackTurnkeyDepositSubmitted(
            record.evm_address,
            chainId,
            amountToUse,
            receipt.transactionHash,
            sourceAssetDenom,
          );
        }

        logger.info({
          at: `${controllerName}#startEvmBridge`,
          message: 'Bridge information record created',
          transactionHash: receipt.transactionHash,
        });
      } catch (error) {
        logger.error({
          at: `${controllerName}#startEvmBridge`,
          message: 'Failed to create bridge information record',
          transactionHash: receipt.transactionHash,
          error,
        });
        // Don't throw error here as the bridge operation was successful
      }
    } catch (error) {
      logger.error({
        at: `${controllerName}#startEvmBridge`,
        message: 'Failed to send user operation',
        callData,
        error,
      });
      throw new Error(`Failed to send user operation, error: ${formatError(error)}`);
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
  // for evm parsing only. Activity will be filled by the webhook if
  // we're parsing a evm transaction.
  if (activity) {
    for (const act of activity) {
      const bridgeOriginAddress = act.toAddress;
      addressesToProcess.set(bridgeOriginAddress, '');
    }
  }
  // for solana parsing only. Transaction will be filled by the webhook if
  // we're parsing a solana transaction.
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
      const record = await TurnkeyUsersTable.findBySvmAddress(bridgeOriginAddress);
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
        record = await TurnkeyUsersTable.findBySmartAccountAddress(checkSummedFromAddress);
      } else {
        record = await TurnkeyUsersTable.findByEvmAddress(checkSummedFromAddress);
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
      const record: TurnkeyUserFromDatabase | undefined = await TurnkeyUsersTable.findByDydxAddress(
        dydxAddress,
      );

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
    const start: number = Date.now();
    try {
      const bridgeController = new BridgeController();
      const { addressesToSweep, chainId } = await parseEvent(req);
      let processedCount = 0;
      let skippedCount = 0;

      // Iterate over the set 'toProcess' and process each item
      for (const fromAddress of addressesToSweep.keys()) {
        const result = await bridgeController.sweep(
          fromAddress,
          chainId,
          addressesToSweep.get(fromAddress) === '' ? undefined : addressesToSweep.get(fromAddress), // amount
        );

        if (result.success && !result.message) {
          processedCount += 1;
        } else if (result.message && result.message.includes('already being processed')) {
          skippedCount += 1;
        }
      }

      logger.info({
        at: `${controllerName}#startBridge`,
        message: 'Bridge processing completed',
        totalAddresses: addressesToSweep.size,
        processedCount,
        skippedCount,
        chainId,
      });

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

router.get(
  '/getDeposits/:dydxAddress',
  ...CheckGetDepositAddressSchema,
  ...CheckLimitAndCreatedBeforeOrAtAndOnOrAfterSchema,
  ...CheckPaginationSchema,
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    try {
      const { dydxAddress } = req.params;
      const {
        limit, page, createdOnOrAfter,
      } = matchedData(req) as {
        limit?: number,
        page?: number,
        createdOnOrAfter?: IsoString,
      };
      const record: TurnkeyUserFromDatabase | undefined = await TurnkeyUsersTable.findByDydxAddress(
        dydxAddress,
      );

      if (!record) {
        return res.status(404).json({
          error: 'User not found',
          message: `No user found with dydx address: ${dydxAddress}`,
        });
      }
      if (!record.smart_account_address) {
        return res.status(404).json({
          error: 'User not found',
          message: `No user found with dydx address: ${dydxAddress}`,
        });
      }

      const deposits = await BridgeInformationTable.searchBridgeInformation(
        {
          from_addresses: [record.evm_address, record.smart_account_address, record.svm_address],
          sinceDate: createdOnOrAfter,
        },
        {
          limit,
          page,
        },
      );

      return res.status(200).json({
        deposits,
        ...(createdOnOrAfter && { since: createdOnOrAfter }),
        total: deposits.total || deposits.results.length,
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
export default router;
