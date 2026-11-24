import {
  concat, keccak256, pad, toHex,
} from 'viem';

/**
 * Fixtures and helpers for the tests of the deposit route to Injective.
 */

// A user's dYdX address, and the Injective account they name for their deposits, in its three
// forms. The two share nothing: the 20 bytes of the Injective address are 0x01 to 0x14.
export const DYDX_ADDRESS: string = 'dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4';
export const INJECTIVE_ADDRESS: string = 'inj1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc54tm65y';
export const INJECTIVE_EVM_ADDRESS: string = '0x0102030405060708090a0B0c0d0e0f1011121314';
export const INJECTIVE_MINT_RECIPIENT: string = '0x0000000000000000000000000102030405060708090a0b0c0d0e0f1011121314';
// Another Injective account, with the address bytes 0x1234567890 four times over.
export const OTHER_INJECTIVE_ADDRESS: string = 'inj1zg69v7yszg69v7yszg69v7yszg69v7ysljkwjv';
export const OTHER_INJECTIVE_MINT_RECIPIENT: string = '0x0000000000000000000000001234567890123456789012345678901234567890';
// The Injective account with the same address bytes as the dYdX address. Deposits do not go
// here unless the user names it.
export const INJECTIVE_ADDRESS_WITH_DYDX_BYTES: string = 'inj199tqg4wdlnu4qjlxchpd7seg454937hjqx2jw6';

// The transaction Skip returned on 2026-11-23 for bridging 100 USDC from Arbitrum to
// INJECTIVE_ADDRESS.
export const SKIP_ARBITRUM_BRIDGE_TX: { to: string, data: string } = {
  to: '0x400BB58033a7763A834199190B68F66A2661aE73',
  data: [
    'd1cc447d',
    '0000000000000000000000000000000000000000000000000000000005f5df86',
    '000000000000000000000000000000000000000000000000000000000000001d',
    '0000000000000000000000000102030405060708090a0b0c0d0e0f1011121314',
    '000000000000000000000000af88d065e77c8cc2239327c5edb3a432268e5831',
    '000000000000000000000000000000000000000000000000000000000000017a',
    '0000000000000000000000000000000000000000000000000000000000000000',
    '00000000000000000000000000000000000000000000000000000000000007d0',
    '000000000000000000000000965540fd2232b7798a795bff93833b2a5feed2e5',
    '0000000000000000000000000000000000000000000000000000000000000120',
    '0000000000000000000000000000000000000000000000000000000000000000',
  ].join(''),
};

// The instructions of the transaction Skip returned on the same day for bridging 100 USDC from
// Solana to the same Injective account.
export const SKIP_SOLANA_BRIDGE_INSTRUCTIONS: { programKey: string, dataHex: string }[] = [
  {
    programKey: 'CCTPV2vPZJS2u2BBsUoscuikbYjnpFmbFsvVuJdgUMQe',
    dataHex: [
      'd73c3d2e723780b0', // discriminator
      '86dff50500000000', // amount
      '1d000000', // destination domain
      '0000000000000000000000000102030405060708090a0b0c0d0e0f1011121314', // mint recipient
      '000000000000000000000000965540fd2232b7798a795bff93833b2a5feed2e5', // destination caller
      '0000000000000000', // max fee
      'd0070000', // min finality threshold
    ].join(''),
  },
  {
    programKey: 'TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA',
    // TransferChecked of 378 units with 6 decimals: the relayer's fee.
    dataHex: '0c7a0100000000000006',
  },
];

/**
 * Evaluates the parts of a Turnkey policy condition that constrain an instruction's program and
 * data, against the instructions of a transaction. Understands only the forms the bridge policy
 * uses, joined by `&&`; any other clause is returned in `unevaluated`.
 */
export function evaluateInstructionClauses(
  condition: string,
  instructions: { programKey: string, dataHex: string }[],
): { result: boolean, unevaluated: string[] } {
  const unevaluated: string[] = [];
  let result: boolean = true;
  condition.split(' && ').forEach((clause: string) => {
    const count: RegExpExecArray | null = (
      /^solana\.tx\.instructions\.count\(\) == (\d+)$/.exec(clause)
    );
    const program: RegExpExecArray | null = (
      /^solana\.tx\.instructions\[(\d+)\]\.program_key == '([^']+)'$/.exec(clause)
    );
    const slice: RegExpExecArray | null = (
      /^solana\.tx\.instructions\[(\d+)\]\.instruction_data_hex\[(\d+)\.\.(\d+)\] == '([^']+)'$/
        .exec(clause)
    );
    if (count !== null) {
      result = result && instructions.length === Number(count[1]);
    } else if (program !== null) {
      result = result && instructions[Number(program[1])]?.programKey === program[2];
    } else if (slice !== null) {
      result = result && instructions[Number(slice[1])]?.dataHex.slice(
        Number(slice[2]),
        Number(slice[3]),
      ) === slice[4];
    } else {
      unevaluated.push(clause);
    }
  });
  return { result, unevaluated };
}

/**
 * A stored approval reduced to the parts that say what it permits: one call, to `target`, with
 * an argument pinned to `pinnedArgument`. For tests that do not load the ZeroDev SDK.
 */
export function makeApproval(target: string, pinnedArgument: string): string {
  return Buffer.from(JSON.stringify({
    permissionParams: {
      policies: [{
        policyParams: {
          type: 'call',
          permissions: [{
            target,
            rules: [{ condition: 0, offset: 64, params: [pinnedArgument] }],
          }],
        },
      }],
    },
  })).toString('base64');
}

export const SKIP_CCTP_V2_RELAYER: string = '0x400BB58033a7763A834199190B68F66A2661aE73';
// An approval for the earlier deposit route: a CCTP v1 relayer, pinned to a Noble address.
export const NOBLE_ROUTE_APPROVAL: string = makeApproval(
  '0x6d4d4A979D766e8b87C985a61B262cfCDf8d6446',
  '0x000000000000000000000000a8a2d0a34dc2e9a6d2e2c2e2f6a9e7f6e3c1b0a9',
);
export const INJECTIVE_ROUTE_APPROVAL: string = makeApproval(
  SKIP_CCTP_V2_RELAYER,
  INJECTIVE_MINT_RECIPIENT,
);

type Hex = `0x${string}`;
interface TypedDataField { name: string, type: string }

/**
 * Hashes typed data the way EIP-712 defines for `eth_signTypedData_v4`, using only what the
 * payload contains: the domain is hashed with the `EIP712Domain` type given in `types`, and
 * nothing is derived. Covers the field types the deposit policy uses, without nested structs.
 */
export function hashTypedDataV4(typedData: {
  domain: object,
  types: { [name: string]: TypedDataField[] },
  primaryType: string,
  message: object,
}): Hex {
  const { types } = typedData;
  if (types.EIP712Domain === undefined) {
    throw new Error('types.EIP712Domain is required');
  }
  const encodeValue = (type: string, value: unknown): Hex => {
    if (type === 'string') {
      return keccak256(toHex(value as string));
    }
    if (type === 'bytes') {
      return keccak256(value as Hex);
    }
    if (type === 'address') {
      return pad(value as Hex, { size: 32 });
    }
    if (/^uint\d+$/.test(type)) {
      return pad(toHex(BigInt(value as number | string)), { size: 32 });
    }
    if (/^bytes\d+$/.test(type)) {
      return pad(value as Hex, { size: 32, dir: 'right' });
    }
    throw new Error(`unsupported type: ${type}`);
  };
  const hashStruct = (name: string, data: object): Hex => {
    const fields: TypedDataField[] = types[name];
    const typeString: string = `${name}(${fields.map((f) => `${f.type} ${f.name}`).join(',')})`;
    return keccak256(concat([
      keccak256(toHex(typeString)),
      ...fields.map((f) => encodeValue(f.type, (data as { [key: string]: unknown })[f.name])),
    ]));
  };
  return keccak256(concat([
    '0x1901',
    hashStruct('EIP712Domain', typedData.domain),
    hashStruct(typedData.primaryType, typedData.message),
  ]));
}
