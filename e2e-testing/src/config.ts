import { baseConfigSchema, parseSchema } from '@klyraprotocol-indexer/base';
import { kafkaConfigSchema } from '@klyraprotocol-indexer/kafka';
import { postgresConfigSchema } from '@klyraprotocol-indexer/postgres';

export const configSchema = {
  ...baseConfigSchema,
  ...postgresConfigSchema,
  ...kafkaConfigSchema,
};

export default parseSchema(configSchema);
