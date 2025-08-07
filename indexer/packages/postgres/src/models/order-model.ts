import path from 'path';

import { Model } from 'objection';

import {
  IntegerPattern,
  NonNegativeNumericPattern,
} from '../lib/validators';
import UpsertQueryBuilder from '../query-builders/upsert';
import {
  IsoString,
  OrderSide,
  OrderStatus,
  OrderType,
  TimeInForce,
} from '../types';
import BaseModel from './base-model';

export default class OrderModel extends BaseModel {
  static get tableName() {
    return 'orders';
  }

  static get idColumn() {
    return 'id';
  }

  static relationMappings = {
    subaccounts: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'subaccount-model'),
      join: {
        from: 'orders.subaccountId',
        to: 'subaccounts.id',
      },
    },
    fills: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'fill-model'),
      join: {
        from: 'orders.id',
        to: 'fills.orderId',
      },
    },
    blocks: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'block-model'),
      join: {
        from: 'orders.createdAtHeight',
        to: 'blocks.blockHeight',
      },
    },
  };

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'id',
        'subaccountId',
        'clientId',
        'clobPairId',
        'side',
        'size',
        'totalFilled',
        'price',
        'type',
        'status',
        'timeInForce',
        'reduceOnly',
        'orderFlags',
        'createdAtHeight',
        'clientMetadata',
        'triggerPrice',
        'updatedAt',
        'updatedAtHeight',
        'orderRouterAddress',
        'duration',
        'interval',
        'priceTolerance',
      ],
      properties: {
        id: { type: 'string', format: 'uuid' },
        subaccountId: { type: 'string', format: 'uuid' },
        clientId: { type: 'string', pattern: IntegerPattern },
        clobPairId: { type: 'string', pattern: IntegerPattern },
        side: { type: 'string', enum: [...Object.values(OrderSide)] },
        size: { type: 'string', pattern: NonNegativeNumericPattern },
        totalFilled: { type: 'string', pattern: NonNegativeNumericPattern },
        price: { type: 'string', pattern: NonNegativeNumericPattern },
        type: { type: 'string', enum: [...Object.values(OrderType)] },
        status: { type: 'string', enum: [...Object.values(OrderStatus)] },
        timeInForce: { type: 'string', enum: [...Object.values(TimeInForce)] },
        reduceOnly: { type: 'boolean' },
        orderFlags: { type: 'string', pattern: IntegerPattern },
        goodTilBlock: { type: ['string', 'null'], default: null, pattern: IntegerPattern },
        goodTilBlockTime: { type: ['string', 'null'], default: null, format: 'date-time' },
        createdAtHeight: { type: ['string', 'null'], default: null, pattern: IntegerPattern },
        clientMetadata: { type: 'string', pattern: IntegerPattern },
        triggerPrice: { type: ['string', 'null'], default: null, pattern: NonNegativeNumericPattern },
        updatedAt: { type: 'string', format: 'date-time' },
        updatedAtHeight: { type: 'string', pattern: IntegerPattern },
        builderAddress: { type: ['string', 'null'], default: null },
        feePpm: { type: ['string', 'null'], default: null },
        orderRouterAddress: { type: ['string', 'null'], default: null },
        duration: { type: ['string', 'null'], default: null, pattern: NonNegativeNumericPattern },
        interval: { type: ['string', 'null'], default: null, pattern: NonNegativeNumericPattern },
        priceTolerance: { type: ['string', 'null'], default: null, pattern: NonNegativeNumericPattern },
      },
    };
  }

  /**
   * A mapping from column name to JSON conversion expected.
   * See getSqlConversionForDydxModelTypes for valid conversions.
   *
   * TODO(IND-239): Ensure that jsonSchema() / sqlToJsonConversions() / model fields match.
   */
  static get sqlToJsonConversions() {
    return {
      id: 'string',
      subaccountId: 'string',
      clientId: 'string',
      clobPairId: 'string',
      side: 'string',
      size: 'string',
      totalFilled: 'string',
      price: 'string',
      type: 'string',
      status: 'string',
      timeInForce: 'string',
      reduceOnly: 'boolean',
      orderFlags: 'string',
      goodTilBlock: 'string',
      goodTilBlockTime: 'date-time',
      createdAtHeight: 'string',
      clientMetadata: 'string',
      triggerPrice: 'string',
      updatedAt: 'date-time',
      updatedAtHeight: 'string',
      duration: 'string',
      interval: 'string',
      priceTolerance: 'string',
      builderAddress: 'string',
      feePpm: 'string',
      orderRouterAddress: 'string',
    };
  }

  id!: string;

  QueryBuilderType!: UpsertQueryBuilder<this>;

  subaccountId!: string;

  clientId!: string;

  clobPairId!: string;

  side!: OrderSide;

  size!: string;

  totalFilled!: string;

  price!: string;

  type!: OrderType;

  status!: OrderStatus;

  timeInForce!: TimeInForce;

  reduceOnly!: boolean;

  orderFlags!: string;

  goodTilBlock!: string;

  goodTilBlockTime!: string;

  createdAtHeight?: string;

  clientMetadata!: string;

  triggerPrice?: string;

  updatedAt!: IsoString;

  updatedAtHeight!: string;

  builderAddress?: string;

  feePpm?: string;

  orderRouterAddress?: string;

  duration?: string;

  interval?: string;

  priceTolerance?: string;
}
