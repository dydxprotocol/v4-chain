import { IntegerPattern } from '../lib/validators';
import { GasSponsorshipStatus, IsoString } from '../types';
import BaseModel from './base-model';

export default class GasSponsorshipModel extends BaseModel {
  static get tableName() {
    return 'gas_sponsorships';
  }

  static get idColumn() {
    return 'id';
  }

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'chain_id',
        'granter_address',
        'grantee_address',
        'dydx_address',
        'spend_limit',
        'fee_amount',
        'status',
        'expires_at',
        'created_at',
      ],
      properties: {
        id: { type: 'string' },
        chain_id: { type: 'string' },
        granter_address: { type: 'string' },
        grantee_address: { type: 'string' },
        dydx_address: { type: 'string' },
        spend_limit: { type: 'string', pattern: IntegerPattern },
        fee_amount: { type: 'string', pattern: IntegerPattern },
        status: { type: 'string', enum: [...Object.values(GasSponsorshipStatus)] },
        transaction_hash: { type: ['string', 'null'] },
        expires_at: { type: 'string' },
        created_at: { type: 'string' },
      },
    };
  }

  id!: string;

  chain_id!: string;

  granter_address!: string;

  grantee_address!: string;

  dydx_address!: string;

  spend_limit!: string;

  fee_amount!: string;

  status!: GasSponsorshipStatus;

  transaction_hash?: string | null;

  expires_at!: IsoString;

  created_at!: IsoString;
}
