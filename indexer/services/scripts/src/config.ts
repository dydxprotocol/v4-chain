import { baseConfigSchema, parseSchema } from '@dydxprotocol-indexer/base';
import { kafkaConfigSchema } from '@dydxprotocol-indexer/kafka';
import { postgresConfigSchema } from '@dydxprotocol-indexer/postgres';

export const configSchema = {
  ...baseConfigSchema,
  ...postgresConfigSchema,
  ...kafkaConfigSchema,
};

export default parseSchema(configSchema);
