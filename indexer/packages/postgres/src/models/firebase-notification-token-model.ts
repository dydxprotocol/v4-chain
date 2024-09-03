import { Model } from 'objection';

import { IsoString } from '../types';
import WalletModel from './wallet-model';

class FirebaseNotificationTokenModel extends Model {
  static tableName = 'firebase_notification_tokens';

  id!: number;
  token!: string;
  address!: string;
  updatedAt!: IsoString;
  language!: string;

  static relationMappings = {
    wallet: {
      relation: Model.BelongsToOneRelation,
      modelClass: WalletModel,
      join: {
        from: 'firebase_notification_tokens.address',
        to: 'wallets.address',
      },
    },
  };
}

export default FirebaseNotificationTokenModel;
