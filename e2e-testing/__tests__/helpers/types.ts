export type OrderDetails = {
  mnemonic: string;
  timeInForce: number;
  orderFlags: number;
  side: number;
  clobPairId: number;
  quantums: number;
  subticks: number;
};
