import { readFileSync, writeFileSync } from 'fs';

const filePath = './public/api-documentation.md';
const data = readFileSync(filePath, 'utf8');
// Define the custom intro we want
const customIntro = `# Klyra Indexer API v1.0.0
> Scroll down for code samples, example requests and responses.
Base URLs:
* For **Testnet**, use <a href="https://klyra-testnet.imperator.co/v4">https://klyra-testnet.imperator.co/v4</a>
Note: Messages on Indexer WebSocket feeds are typically more recent than data fetched via Indexer's REST API, because the latter is backed by read replicas of the databases that feed the former. Ordinarily this difference is minimal (less than a second), but it might become prolonged under load.`;
// Find the start of the content (after the intro sections)
const contentStart = data.indexOf('# Authentication');
// Combine our custom intro with the existing content
const cleanedContent = `${customIntro}\n\n${data.slice(contentStart)}`;
writeFileSync(filePath, cleanedContent);
