import { v5 as uuidv5 } from 'uuid';

// not sure if this needs to be environment specific and be loaded in env
export const NAMESPACE = '0f9da948-a6fb-4c45-9edc-4685c3f3317d';

export function getUuid(bufferValue: Buffer) {
  return uuidv5(bufferValue, NAMESPACE);
}
