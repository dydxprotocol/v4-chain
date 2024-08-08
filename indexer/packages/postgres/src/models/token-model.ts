import { Model } from 'objection';

import { IsoString } from '../types';
import WalletModel from './wallet-model';

class TokenModel extends Model {
  static tableName = 'tokens';

  id!: number;
  token!: string;
  address!: string;
  updatedAt!: IsoString;

  static relationMappings = {
    wallet: {
      relation: Model.BelongsToOneRelation,
      modelClass: WalletModel,
      join: {
        from: 'tokens.address',
        to: 'wallets.address',
      },
    },
  };
}

export default TokenModel;
