import path from 'path';

import { Model } from 'objection';

import {
  IntegerPattern,
  NonNegativeNumericPattern,
  NumericPattern,
} from '../lib/validators';
import {
  FillType,
  IsoString,
  Liquidity,
  OrderSide,
  PositionSide,
} from '../types';

export default class FillModel extends Model {
  static get tableName() {
    return 'fills';
  }

  static get idColumn() {
    return 'id';
  }

  static relationMappings = {
    subaccounts: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'subaccount-model'),
      join: {
        from: 'fills.subaccountId',
        to: 'subaccounts.id',
      },
    },
    tendermintEvents: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'tendermint-event-model'),
      join: {
        from: 'fills.eventId',
        to: 'tendermint_events.id',
      },
    },
    blocks: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'block-model'),
      join: {
        from: 'fills.createdAtHeight',
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
        'side',
        'liquidity',
        'type',
        'clobPairId',
        'orderId',
        'size',
        'price',
        'quoteAmount',
        'eventId',
        'transactionHash',
        'createdAt',
        'createdAtHeight',
        'fee',
        'affiliateRevShare',
      ],
      properties: {
        id: { type: 'string', format: 'uuid' },
        subaccountId: { type: 'string', format: 'uuid' },
        side: { type: 'string', enum: [...Object.values(OrderSide)] },
        liquidity: { type: 'string', enum: [...Object.values(Liquidity)] },
        type: { type: 'string', enum: [...Object.values(FillType)] },
        clobPairId: { type: 'string', pattern: IntegerPattern },
        orderId: { type: ['string', 'null'], default: null, format: 'uuid' },
        size: { type: 'string', pattern: NonNegativeNumericPattern },
        price: { type: 'string', pattern: NonNegativeNumericPattern },
        quoteAmount: { type: 'string', pattern: NonNegativeNumericPattern },
        transactionHash: { type: 'string' },
        createdAt: { type: 'string', format: 'date-time' },
        createdAtHeight: { type: 'string', pattern: IntegerPattern },
        clientMetadata: { type: ['string', 'null'], pattern: IntegerPattern },
        fee: { type: 'string', pattern: NumericPattern },
        affiliateRevShare: { type: 'string', pattern: NonNegativeNumericPattern },
        builderFee: { type: ['string', 'null'], default: null },
        builderAddress: { type: ['string', 'null'], default: null },
        orderRouterAddress: { type: ['string', 'null'], default: null },
        orderRouterFee: { type: ['string', 'null'], default: null },
        positionSizeBefore: { type: ['string', 'null'], pattern: NumericPattern, default: null },
        entryPriceBefore: { type: ['string', 'null'], pattern: NonNegativeNumericPattern, default: null },
        positionSideBefore: { type: ['string', 'null'], enum: [...Object.values(PositionSide), null], default: null },
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
      side: 'string',
      liquidity: 'string',
      type: 'string',
      clobPairId: 'string',
      orderId: 'string',
      size: 'string',
      price: 'string',
      quoteAmount: 'string',
      eventId: 'hex-string',
      transactionHash: 'string',
      createdAt: 'date-time',
      createdAtHeight: 'string',
      clientMetadata: 'string',
      fee: 'string',
      affiliateRevShare: 'string',
      builderFee: 'string',
      builderAddress: 'string',
      orderRouterAddress: 'string',
      orderRouterFee: 'string',
      positionSizeBefore: 'string',
      entryPriceBefore: 'string',
      positionSideBefore: 'string',
    };
  }

  id!: string;

  subaccountId!: string;

  side!: OrderSide;

  liquidity!: Liquidity;

  type!: FillType;

  clobPairId!: string;

  orderId!: string;

  size!: string;

  price!: string;

  quoteAmount!: string;

  eventId!: Buffer;

  transactionHash!: string;

  createdAt!: IsoString;

  createdAtHeight!: string;

  clientMetadata!: string;

  fee!: string;

  affiliateRevShare!: string;

  builderFee!: string;

  builderAddress!: string;

  orderRouterAddress!: string;

  orderRouterFee!: string;

  positionSizeBefore?: string;

  entryPriceBefore?: string;

  positionSideBefore?: PositionSide;
}
