import path from 'path';

import { Model } from 'objection';

import {
  IntegerPattern,
  NonNegativeNumericPattern,
  NumericPattern,
} from '../lib/validators';
import { IsoString, PositionSide, PerpetualPositionStatus } from '../types';

export default class PerpetualPositionModel extends Model {
  static get tableName() {
    return 'perpetual_positions';
  }

  static get idColumn() {
    return 'id';
  }

  static relationMappings = {
    subaccount: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'subaccount-model'),
      join: {
        from: 'perpetual_positions.subaccountId',
        to: 'subaccounts.id',
      },
    },
    perpetualMarket: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'perpetual-market-model'),
      join: {
        from: 'perpetual_positions.perpetualId',
        to: 'perpetual_markets.id',
      },
    },
    openTendermintEvent: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'tendermint-event-model'),
      join: {
        from: 'perpetual_positions.openEventId',
        to: 'tendermint_events.id',
      },
    },
    closeTendermintEvent: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'tendermint-event-model'),
      join: {
        from: 'perpetual_positions.closeEventId',
        to: 'tendermint_events.id',
      },
    },
    lastTendermintEvent: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'tendermint-event-model'),
      join: {
        from: 'perpetual_positions.lastEventId',
        to: 'tendermint_events.id',
      },
    },
  };

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'id',
        'subaccountId',
        'perpetualId',
        'side',
        'status',
        'size',
        'maxSize',
        'entryPrice',
        'sumOpen',
        'createdAt',
        'createdAtHeight',
        'openEventId',
        'lastEventId',
        'settledFunding',
      ],
      properties: {
        id: { type: 'string', format: 'uuid' },
        subaccountId: { type: 'string', format: 'uuid' },
        perpetualId: { type: 'string', pattern: IntegerPattern },
        side: { type: 'string', enum: [...Object.values(PositionSide)] },
        status: { type: 'string', enum: [...Object.values(PerpetualPositionStatus)] },
        size: { type: 'string', pattern: NumericPattern },
        maxSize: { type: 'string', pattern: NumericPattern },
        entryPrice: { type: 'string', pattern: NonNegativeNumericPattern },
        exitPrice: { type: ['string', 'null'], pattern: NonNegativeNumericPattern, default: null },
        sumOpen: { type: 'string', pattern: NonNegativeNumericPattern },
        sumClose: { type: 'string', pattern: NonNegativeNumericPattern },
        createdAt: { type: 'string', format: 'date-time' },
        closedAt: { type: ['string', 'null'], format: 'date-time', default: null },
        createdAtHeight: { type: 'string', pattern: IntegerPattern },
        closedAtHeight: { type: ['string', 'null'], default: null, pattern: IntegerPattern },
        settledFunding: { type: 'string', pattern: NumericPattern },
        totalRealizedPnl: { type: ['string', 'null'], default: null, pattern: NumericPattern },
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
      perpetualId: 'string',
      side: 'string',
      status: 'string',
      size: 'string',
      maxSize: 'string',
      entryPrice: 'string',
      exitPrice: 'string',
      sumOpen: 'string',
      sumClose: 'string',
      createdAt: 'date-time',
      closedAt: 'date-time',
      createdAtHeight: 'string',
      closedAtHeight: 'string',
      openEventId: 'hex-string',
      closeEventId: 'hex-string',
      lastEventId: 'hex-string',
      settledFunding: 'string',
      totalRealizedPnl: 'string',
    };
  }

  id!: string;

  subaccountId!: string;

  perpetualId!: string;

  side!: PositionSide;

  status!: PerpetualPositionStatus;

  size!: string;

  maxSize!: string;

  entryPrice!: string;

  exitPrice!: string;

  sumOpen!: string;

  sumClose!: string;

  createdAt!: IsoString;

  closedAt!: IsoString;

  createdAtHeight!: string;

  closedAtHeight!: string;

  openEventId!: Buffer;

  closeEventId?: Buffer;

  lastEventId!: Buffer;

  settledFunding!: string;

  totalRealizedPnl?: string;
}
