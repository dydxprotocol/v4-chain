import { Transfer, OrderFlags, ICancelOrder, IPlaceOrder } from '../clients/types';
import { UserError } from './errors';
/**
 * @describe validatePlaceOrderMessage validates that an order to place has fields that would be
 *  valid on-chain.
 */
export declare function validatePlaceOrderMessage(subaccountNumber: number, order: IPlaceOrder): UserError | undefined;
/**
 * @describe validateCancelOrderMessage validates that an order to cancel has fields that would be
 *  valid on-chain.
 */
export declare function validateCancelOrderMessage(subaccountNumber: number, order: ICancelOrder): UserError | undefined;
/**
 * @describe validateTransferMessage validates that a transfer to place has fields that would be
 *  valid on-chain.
 */
export declare function validateTransferMessage(transfer: Transfer): UserError | undefined;
export declare function verifyOrderFlags(orderFlags: OrderFlags): boolean;
export declare function isStatefulOrder(orderFlags: OrderFlags): boolean;
export declare function isValidAddress(address: string): boolean;
