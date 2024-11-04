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
  
  // Set the title and version
  swaggerDocument.info.title = 'Klyra Indexer API';
  swaggerDocument.info.version = 'v1.0.0';
  
  // Set the description that will appear in the intro
  swaggerDocument.info.description = `Scroll down for code samples, example requests and responses.
Base URLs:
* For **Testnet**, use <a href="https://klyra-testnet.imperator.co/v4">https://klyra-testnet.imperator.co/v4</a>
Note: Messages on Indexer WebSocket feeds are typically more recent than data fetched via Indexer's REST API, because the latter is backed by read replicas of the databases that feed the former. Ordinarily this difference is minimal (less than a second), but it might become prolonged under load.`;
  
  // Set the server
  swaggerDocument.servers = [
    {
      url: 'https://klyra-testnet.imperator.co/v4',
      description: 'Testnet',
    },
  ];
  writeFileSync(filePath, JSON.stringify(swaggerDocument, null, 2));
});