import { bigIntToBytes } from '@klyraprotocol-indexer/v4-proto-parser';

export function intToSerializedInt(int: number): Uint8Array {
  return bigIntToBytes(BigInt(int));
}
