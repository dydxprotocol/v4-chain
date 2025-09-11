import { CallPolicyVersion, ParamCondition, type CallPolicyParams } from '@zerodev/permissions/policies';
import type { Abi } from 'viem';
import {
  arbitrum, avalanche, base, mainnet, optimism,
} from 'viem/chains';

import config from '../config';
import { encodeToHexAndPad, getNobleForwardingAddress, nobleToHex } from '../helpers/skip-helper';
import { abi } from './bridge-abi';
import { usdcAddressByChainId } from './smart-contract-constants';

const goFastHandlerProxyByChainId: Record<string, string> = {
  [arbitrum.id.toString()]: '0x4c58aE019E54D10594F1Aa26ABF385B6fb17A52d',
  [base.id.toString()]: '0x9335C0c0CBc0317291fd48c00b2f71C8b39DA6F8',
  [optimism.id.toString()]: '0x9c540EdC86613b22968Da784b2d42AC79965af91',
  [avalanche.id.toString()]: '0xb7B287F15e5edDFEfF2b05ef1BE7F7cc73197AaA',
  [mainnet.id.toString()]: '0xa11CC0eFb1B3AcD95a2B8cd316E8c132E16048b5',
};

// this value limit is set to restrict usdc transfers less than 100k.
const valueLimit = config.CALL_POLICY_VALUE_LIMIT;

/**
 * Construct the policy for the given chainId and dydxAddress. Consolidates call policy construction
 * for avax, arbitrum, base, optimism. Ethereum is not supported because it also uses a different
 * smart contract address for slow bridges. Used to give us permission to kick off the bridge.
 *
 * @param dydxAddress
 * @param chainId - one of arbitrum, base, optimism, avalanche.
 * @returns the policy for the given chainId and dydxAddress.
 */
function constructPolicy(chainId: string): (dydxAddress: string) => Promise<CallPolicyParams<Abi, `0x${string}`>> {
  const goFastHandlerProxy = goFastHandlerProxyByChainId[chainId] as `0x${string}`;
  return (dydxAddress: string) => Promise.resolve({
    policyVersion: CallPolicyVersion.V0_0_5,
    permissions: [
      {
        target: usdcAddressByChainId[chainId] as `0x${string}`, // usdc on chainId
        abi,
        valueLimit,
        functionName: 'approve',
        args: [
          {
            condition: ParamCondition.EQUAL,
            value: goFastHandlerProxy,
          },
          null,
        ],
      },
      {
        target: goFastHandlerProxy,
        abi,
        valueLimit,
        functionName: 'submitOrder',
        args: [
          null,
          null,
          null,
          null,
          null,
          null,
          {
            condition: ParamCondition.SLICE_EQUAL,
            value: encodeToHexAndPad(dydxAddress),
            start: 277,
            length: 42,
          },
        ],
      },
      {
        target: goFastHandlerProxy,
        abi,
        valueLimit,
        functionName: 'swapAndSubmitOrder',
        args: [
          null,
          null,
          null,
          null,
          null,
          null,
          null,
          null,
          null,
          {
            condition: ParamCondition.SLICE_EQUAL,
            value: encodeToHexAndPad(dydxAddress),
            start: 277,
            length: 42,
          },
        ],
      },
    ],
  });
}

export async function getEthereumCallPolicy(dydxAddress: string): Promise<CallPolicyParams<Abi, `0x${string}`>> {
  const ethCCTPRelayerProxy = '0xf33e750336e9C0D4E2f4c0D450d753030693CC71';
  const goFastHandlerProxy = goFastHandlerProxyByChainId[mainnet.id.toString()] as `0x${string}`;
  // get the noble forwarding address.
  // cctp mints to a noble forwarding address which forwards to the dydx address.
  const nobleForwardingAddress = await getNobleForwardingAddress(dydxAddress);
  const nobleForwardingAddressEvm = nobleToHex(nobleForwardingAddress);
  // for go fast
  const destinationCallataAddr = encodeToHexAndPad(dydxAddress);
  return {
    policyVersion: CallPolicyVersion.V0_0_5,
    permissions: [
      // slow deposits, via CCTP Relayer. USDC Bridge
      {
        target: ethCCTPRelayerProxy,
        abi,
        valueLimit,
        functionName: 'requestCCTPTransferWithCaller', // for usdc bridges
        args: [
          null,
          null,
          {
            // mint recipient is the noble forwarding address.
            condition: ParamCondition.EQUAL,
            value: nobleForwardingAddressEvm,
          },
          null,
          null,
          null,
        ],
      },
      // slow deposits, via CCTP Relayer. ETH Bridge
      {
        target: ethCCTPRelayerProxy,
        abi,
        valueLimit,
        functionName: 'swapAndRequestCCTPTransferWithCaller', // for eth bridges
        args: [
          null,
          null,
          null,
          null,
          {
            // mint recipient is the noble forwarding address.
            condition: ParamCondition.EQUAL,
            value: nobleForwardingAddressEvm,
          },
          null,
          null,
          null,
        ],
      },
      // allow skip.go bridge smart contract permissions.
      {
        target: usdcAddressByChainId[mainnet.id.toString()] as `0x${string}`, // usdc on ethereum
        abi,
        valueLimit,
        functionName: 'approve',
        args: [
          {
            condition: ParamCondition.ONE_OF,
            value: [
              ethCCTPRelayerProxy,
              goFastHandlerProxy,
            ],
          },
          null,
        ],
      },
      // skip go fast bridges.
      {
        target: goFastHandlerProxy,
        abi,
        valueLimit,
        functionName: 'submitOrder',
        args: [
          null,
          null,
          null,
          null,
          null,
          null,
          {
            condition: ParamCondition.SLICE_EQUAL,
            value: destinationCallataAddr,
            start: 277,
            length: 42,
          },
        ],
      },
      {
        target: goFastHandlerProxy,
        abi,
        valueLimit,
        functionName: 'swapAndSubmitOrder',
        args: [
          null,
          null,
          null,
          null,
          null,
          null,
          null,
          null,
          null,
          {
            condition: ParamCondition.SLICE_EQUAL,
            value: destinationCallataAddr,
            start: 277,
            length: 42,
          },
        ],
      },
    ],
  };
}

export const callPolicyByChainId: Record<string, (dydxAddress: string) => Promise<CallPolicyParams<Abi, `0x${string}`>>> = {
  [mainnet.id.toString()]: getEthereumCallPolicy,
  [arbitrum.id.toString()]: constructPolicy(arbitrum.id.toString()),
  [avalanche.id.toString()]: constructPolicy(avalanche.id.toString()),
  [optimism.id.toString()]: constructPolicy(optimism.id.toString()),
  [base.id.toString()]: constructPolicy(base.id.toString()),
};
