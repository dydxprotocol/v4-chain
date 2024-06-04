#!/usr/bin/env node
import { prompt } from './prompt';
import { cli } from './cli';
import { readFileSync } from 'fs';
const argv = require('minimist')(process.argv.slice(2));

const question = [
  {
    _: true,
    type: 'string',
    name: 'file',
    message: 'file'
  }
];

(async () => {
  const { file } = await prompt(question, argv);
  const argvv = JSON.parse(readFileSync(file, 'utf-8'));
  await cli(argvv);
})();
