import { CallPolicyVersion, ParamCondition, type CallPolicyParams } from '@zerodev/permissions/policies';
import type { Abi } from 'viem';
import {
  arbitrum, avalanche, base, mainnet, optimism,
} from 'viem/chains';

import { encodeToHexAndPad, getNobleForwardingAddress, nobleToHex } from '../helpers/skip-helper';
import { usdcAddressByChainId } from './smart-contract-constants';

export const abi = [
  {
    inputs: [
      {
        internalType: 'bytes32',
        name: 'sender',
        type: 'bytes32',
      },
      {
        internalType: 'bytes32',
        name: 'recipient',
        type: 'bytes32',
      },
      {
        internalType: 'uint256',
        name: 'amountIn',
        type: 'uint256',
      },
      {
        internalType: 'uint256',
        name: 'amountOut',
        type: 'uint256',
      },
      {
        internalType: 'uint32',
        name: 'destinationDomain',
        type: 'uint32',
      },
      {
        internalType: 'uint64',
        name: 'timeoutTimestamp',
        type: 'uint64',
      },
      {
        internalType: 'bytes',
        name: 'data',
        type: 'bytes',
      },
    ],
    name: 'submitOrder',
    outputs: [
      {
        internalType: 'bytes32',
        name: '',
        type: 'bytes32',
      },
    ],
    stateMutability: 'nonpayable',
    type: 'function',
  },
  {
    inputs: [
      {
        internalType: 'address',
        name: 'tokenIn',
        type: 'address',
      },
      {
        internalType: 'uint256',
        name: 'swapAmountIn',
        type: 'uint256',
      },
      {
        internalType: 'bytes',
        name: 'swapCalldata',
        type: 'bytes',
      },
      {
        internalType: 'uint256',
        name: 'executionFeeAmount',
        type: 'uint256',
      },
      {
        internalType: 'uint256',
        name: 'solverFeeBPS',
        type: 'uint256',
      },
      {
        internalType: 'bytes32',
        name: 'sender',
        type: 'bytes32',
      },
      {
        internalType: 'bytes32',
        name: 'recipient',
        type: 'bytes32',
      },
      {
        internalType: 'uint32',
        name: 'destinationDomain',
        type: 'uint32',
      },
      {
        internalType: 'uint64',
        name: 'timeoutTimestamp',
        type: 'uint64',
      },
      {
        internalType: 'bytes',
        name: 'destinationCalldata',
        type: 'bytes',
      },
    ],
    name: 'swapAndSubmitOrder',
    outputs: [
      {
        internalType: 'bytes32',
        name: '',
        type: 'bytes32',
      },
    ],
    stateMutability: 'payable',
    type: 'function',
  },
  {
    type: 'function',
    name: 'approve',
    stateMutability: 'nonpayable',
    inputs: [
      {
        name: 'spender',
        type: 'address',
      },
      {
        name: 'amount',
        type: 'uint256',
      },
    ],
    outputs: [
      {
        type: 'bool',
      },
    ],
  },
  {
    inputs: [
      {
        internalType: 'address',
        name: 'inputToken',
        type: 'address',
      },
      {
        internalType: 'uint256',
        name: 'inputAmount',
        type: 'uint256',
      },
      {
        internalType: 'bytes',
        name: 'swapCalldata',
        type: 'bytes',
      },
      {
        internalType: 'uint32',
        name: 'destinationDomain',
        type: 'uint32',
      },
      {
        internalType: 'bytes32',
        name: 'mintRecipient',
        type: 'bytes32',
      },
      {
        internalType: 'address',
        name: 'burnToken',
        type: 'address',
      },
      {
        internalType: 'uint256',
        name: 'feeAmount',
        type: 'uint256',
      },
      {
        internalType: 'bytes32',
        name: 'destinationCaller',
        type: 'bytes32',
      },
    ],
    name: 'swapAndRequestCCTPTransferWithCaller',
    outputs: [],
    stateMutability: 'payable',
    type: 'function',
  },
  {
    inputs: [
      {
        internalType: 'uint256',
        name: 'transferAmount',
        type: 'uint256',
      },
      {
        internalType: 'uint32',
        name: 'destinationDomain',
        type: 'uint32',
      },
      {
        internalType: 'bytes32',
        name: 'mintRecipient',
        type: 'bytes32',
      },
      {
        internalType: 'address',
        name: 'burnToken',
        type: 'address',
      },
      {
        internalType: 'uint256',
        name: 'feeAmount',
        type: 'uint256',
      },
      {
        internalType: 'bytes32',
        name: 'destinationCaller',
        type: 'bytes32',
      },
    ],
    name: 'requestCCTPTransferWithCaller',
    outputs: [],
    stateMutability: 'nonpayable',
    type: 'function',
  },
] as const;

export function getArbitrumCallPolicy(dydxAddress: string): Promise<CallPolicyParams<Abi, `0x${string}`>> {
  const destinationCallataAddr = encodeToHexAndPad(dydxAddress);
  const goFastHandlerProxy = '0x4c58aE019E54D10594F1Aa26ABF385B6fb17A52d' as `0x${string}`;
  return Promise.resolve({
    policyVersion: CallPolicyVersion.V0_0_5,
    permissions: [
      {
        target: usdcAddressByChainId[arbitrum.id.toString()] as `0x${string}`, // usdc on arbitrum
        abi,
        valueLimit: BigInt(1000000000000000000000000000000),
        functionName: 'approve',
      },
      {
        target: goFastHandlerProxy,
        abi,
        valueLimit: BigInt(1000000000000000000000000000000),
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
        valueLimit: BigInt(1000000000000000000000000000000),
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
  });
}

export function getBaseCallPolicy(dydxAddress: string): Promise<CallPolicyParams<Abi, `0x${string}`>> {
  const destinationCallataAddr = encodeToHexAndPad(dydxAddress);
  const goFastHandlerProxy = '0x9335C0c0CBc0317291fd48c00b2f71C8b39DA6F8' as `0x${string}`;
  return Promise.resolve({
    policyVersion: CallPolicyVersion.V0_0_5,
    permissions: [
      {
        target: usdcAddressByChainId[base.id.toString()] as `0x${string}`, // usdc on avalanche
        abi,
        valueLimit: BigInt(1000000000000000000000000000000),
        functionName: 'approve',
      },
      {
        target: goFastHandlerProxy,
        abi,
        valueLimit: BigInt(1000000000000000000000000000000),
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
        valueLimit: BigInt(1000000000000000000000000000000),
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
  });
}

export function getOptimismCallPolicy(dydxAddress: string): Promise<CallPolicyParams<Abi, `0x${string}`>> {
  const destinationCallataAddr = encodeToHexAndPad(dydxAddress);
  const goFastHandlerProxy = '0x9c540EdC86613b22968Da784b2d42AC79965af91' as `0x${string}`;
  return Promise.resolve({
    policyVersion: CallPolicyVersion.V0_0_5,
    permissions: [
      {
        target: usdcAddressByChainId[optimism.id.toString()] as `0x${string}`, // usdc on avalanche
        abi,
        valueLimit: BigInt(1000000000000000000000000000000),
        functionName: 'approve',
      },
      {
        target: goFastHandlerProxy,
        abi,
        valueLimit: BigInt(1000000000000000000000000000000),
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
        valueLimit: BigInt(1000000000000000000000000000000),
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
  });
}

export function getAvalancheCallPolicy(dydxAddress: string): Promise<CallPolicyParams<Abi, `0x${string}`>> {
  const destinationCallataAddr = encodeToHexAndPad(dydxAddress);
  const goFastHandlerProxy = '0xb7B287F15e5edDFEfF2b05ef1BE7F7cc73197AaA' as `0x${string}`;
  return Promise.resolve({
    policyVersion: CallPolicyVersion.V0_0_5,
    permissions: [
      {
        target: usdcAddressByChainId[avalanche.id.toString()] as `0x${string}`, // usdc on avalanche
        abi,
        valueLimit: BigInt(1000000000000000000000000000000),
        functionName: 'approve',
      },
      {
        target: goFastHandlerProxy,
        abi,
        valueLimit: BigInt(1000000000000000000000000000000),
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
        valueLimit: BigInt(1000000000000000000000000000000),
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
  });
}

export async function getEthereumCallPolicy(dydxAddress: string): Promise<CallPolicyParams<Abi, `0x${string}`>> {
  const ethCCTPRelayerProxy = '0xf33e750336e9C0D4E2f4c0D450d753030693CC71';
  const goFastHandlerProxy = '0xa11CC0eFb1B3AcD95a2B8cd316E8c132E16048b5';
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
        valueLimit: BigInt(1000000000000000000000000000000),
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
        valueLimit: BigInt(1000000000000000000000000000000),
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
        valueLimit: BigInt(1000000000000000000000000000000),
        functionName: 'approve',
        args: [
          {
            condition: ParamCondition.ONE_OF,
            value: [
              ethCCTPRelayerProxy,
              goFastHandlerProxy,
            ],
          },
        ],
      },
      // skip go fast bridges.
      {
        target: goFastHandlerProxy,
        abi,
        valueLimit: BigInt(1000000000000000000000000000000),
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
        valueLimit: BigInt(1000000000000000000000000000000),
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
  [arbitrum.id.toString()]: getArbitrumCallPolicy,
  [mainnet.id.toString()]: getEthereumCallPolicy,
  [avalanche.id.toString()]: getAvalancheCallPolicy,
  [optimism.id.toString()]: getOptimismCallPolicy,
  [base.id.toString()]: getBaseCallPolicy,
};
