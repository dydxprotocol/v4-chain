import path from 'path';

import Piscina from 'piscina';

export const piscina = new Piscina({
  filename: path.resolve(__dirname, 'lib/workers/from-kafka-helpers.js'),
  minThreads: 1,
  maxThreads: 2,
});
