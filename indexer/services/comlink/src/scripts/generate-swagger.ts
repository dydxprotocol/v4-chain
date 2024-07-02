import { readFileSync, writeFileSync } from 'fs';

import { generateSpec } from 'tsoa';

// eslint-disable-next-line @typescript-eslint/no-floating-promises
generateSpec({
  entryFile: '../index.ts',
  outputDirectory: 'public',
  noImplicitAdditionalProperties: 'throw-on-extras',
  controllerPathGlobs: ['./src/controllers/api/v4/**/*.ts'],
  specVersion: 3,
}).then(() => {
  const filePath: string = './public/swagger.json';
  const data: string = readFileSync(filePath, 'utf8');
  const swaggerDocument = JSON.parse(data);
  swaggerDocument.info.title = 'Indexer API';
  swaggerDocument.info.version = 'v1.0.0';
  swaggerDocument.servers = [
    {
      url: '',
      description: 'Public Testnet',
    },
  ];
  writeFileSync(filePath, JSON.stringify(swaggerDocument, null, 2));
});
