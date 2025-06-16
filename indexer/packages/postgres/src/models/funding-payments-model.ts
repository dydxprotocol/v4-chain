import path from 'path';

import { Model } from 'objection';

import { IntegerPattern, NumericPattern } from '../lib/validators';
import { IsoString, PositionSide } from '../types';

export default class FundingPaymentsModel extends Model {
  static get tableName() {
    return 'funding_payments';
  }

  static get idColumn() {
    return ['subaccountId', 'createdAt', 'ticker'];
  }

  static relationMappings = {
    subaccount: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'subaccount-model'),
      join: {
        from: 'funding_payments.subaccountId',
        to: 'subaccounts.id',
      },
    },
    perpetualMarket: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'perpetual-market-model'),
      join: {
        from: 'funding_payments.perpetualId',
        to: 'perpetual_markets.id',
      },
    },
    blocks: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'block-model'),
      join: {
        from: 'funding_payments.createdAtHeight',
        to: 'blocks.blockHeight',
      },
    },
  };

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'subaccountId',
        'createdAt',
        'createdAtHeight',
        'perpetualId',
        'ticker',
        'oraclePrice',
        'size',
        'side',
        'rate',
        'payment',
        'fundingIndex',
      ],
      properties: {
        subaccountId: { type: 'string', format: 'uuid' },
        createdAt: { type: 'string', format: 'date-time' },
        createdAtHeight: { type: 'string', pattern: IntegerPattern },
        perpetualId: { type: 'string', pattern: IntegerPattern },
        ticker: { type: 'string' },
        oraclePrice: { type: 'string', pattern: NumericPattern },
        size: { type: 'string', pattern: NumericPattern },
        side: { type: 'string', enum: [...Object.values(PositionSide)] },
        rate: { type: 'string', pattern: NumericPattern },
        payment: { type: 'string', pattern: NumericPattern },
        fundingIndex: { type: 'string', pattern: NumericPattern },
      },
    };
  }

  static get sqlToJsonConversions() {
    return {
      subaccountId: 'string',
      createdAt: 'date-time',
      createdAtHeight: 'string',
      perpetualId: 'string',
      ticker: 'string',
      oraclePrice: 'string',
      size: 'string',
      side: 'string',
      rate: 'string',
      payment: 'string',
      fundingIndex: 'string',
    };
  }

  subaccountId!: string;
  createdAt!: IsoString;
  createdAtHeight!: string;
  perpetualId!: string;
  ticker!: string;
  oraclePrice!: string;
  size!: string;
  side!: PositionSide;
  rate!: string;
  payment!: string;
  fundingIndex!: string;
}
