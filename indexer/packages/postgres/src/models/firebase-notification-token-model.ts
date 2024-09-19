import { Model } from 'objection';

import { IsoString } from '../types';
import WalletModel from './wallet-model';

class FirebaseNotificationTokenModel extends Model {
  static get tableName() {
    return 'firebase_notification_tokens';
  }

  static get idColumn() {
    return 'id';
  }

  static get jsonSchema() {
    return {
    };
  }

  static get sqlToJsonConversions() {
    return {
    };
  }

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

  id!: number;

  token!: string;

  address!: string;

  updatedAt!: IsoString;

  language!: string;
}

export default FirebaseNotificationTokenModel;
