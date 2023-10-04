export function base64StringToBinary(base64Message: string): Uint8Array {
  return new Uint8Array(
    Buffer.from(base64Message, 'base64'),
  );
}
