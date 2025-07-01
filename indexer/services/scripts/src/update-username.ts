import fs from 'fs';

import {
  SubaccountTable,
} from '@dydxprotocol-indexer/postgres';
import { parse } from 'csv-parse/sync';
import yargs from 'yargs';

interface AffiliateUpdate {
  address: string,
  currentName: string,
  proposedName: string,
}

function getSubaccountId(address: string, subaccountNumber: number): string {
  return SubaccountTable.uuid(
    address,
    subaccountNumber,
  );
}

function printUsernameUpdateSql(
  address: string,
  currentName: string,
  proposedName: string,
): void {
  const subaccountId: string = getSubaccountId(address, 0);

  // Print SQL statements directly to stdout
  console.log('-- Get current username:');
  console.log(`SELECT * FROM subaccount_usernames WHERE "subaccountId" = '${subaccountId}' AND username = '${currentName}';`);
  console.log('\n-- Update username:');
  console.log(`UPDATE subaccount_usernames SET username = '${proposedName}' WHERE "subaccountId" = '${subaccountId}' AND username = '${currentName}';`);
}

function processCsvFile(filePath: string): void {
  const fileContent = fs.readFileSync(filePath, 'utf-8');
  const records = parse(fileContent, {
    columns: true,
    skip_empty_lines: true,
  });

  const updates: AffiliateUpdate[] = records.map((record: any) => ({
    address: record['dYdX address'],
    currentName: record['Current Affifliate Name'],
    proposedName: record['Proposal Name Change'],
  }));

  // Create conditions for the WHERE clause
  const conditions = updates.map((update) => {
    const subaccountId = getSubaccountId(update.address, 0);
    return `("subaccountId" = '${subaccountId}' AND username = '${update.currentName}')`;
  });

  // Print single SELECT statement for all rows
  console.log('-- Verify all current usernames:');
  console.log(`SELECT * FROM subaccount_usernames WHERE ${conditions.join(' OR ')};`);

  // Print single UPDATE statement for all rows
  console.log('\n-- Update all usernames:');
  const cases = updates.map((update) => {
    const subaccountId = getSubaccountId(update.address, 0);
    return `WHEN "subaccountId" = '${subaccountId}' AND username = '${update.currentName}' THEN '${update.proposedName}'`;
  });
  
  console.log(`UPDATE subaccount_usernames SET username = CASE
${cases.join('\n')}
END
WHERE ${conditions.join(' OR ')};`);
}

const args = yargs
  .conflicts('csvFile', ['address', 'currentName', 'name'])
  .options({
    csvFile: {
      type: 'string',
      alias: 'f',
      description: 'Path to CSV file with updates',
    },
    address: {
      type: 'string',
      alias: 'a',
      description: 'dYdX address',
    },
    currentName: {
      type: 'string',
      alias: 'c',
      description: 'Current affiliate name',
    },
    name: {
      type: 'string',
      alias: 'n',
      description: 'Proposed new username',
    },
  })
  .check((argv) => {
    if (!argv.csvFile && (!argv.address || !argv.currentName || !argv.name)) {
      throw new Error('Either provide a CSV file or all three: address, currentName, and name');
    }
    return true;
  })
  .argv;

if (args.csvFile) {
  processCsvFile(args.csvFile);
} else {
  printUsernameUpdateSql(args.address!, args.currentName!, args.name!);
}
