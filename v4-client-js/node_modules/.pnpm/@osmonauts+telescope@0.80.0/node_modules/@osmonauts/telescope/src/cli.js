import { prompt } from './prompt';
import { Commands as commands } from './cmds';
import { Contracts as contracts } from './cmds';

export const cli = async (argv) => {
  if (argv.contract) {
    const { cmd } = await prompt(
      [
        {
          _: true,
          type: 'fuzzy',
          name: 'cmd',
          message: 'what do you want to do?',
          choices: Object.keys(contracts)
        }
      ],
      argv
    );
    if (typeof contracts[cmd] === 'function') {
      await contracts[cmd](argv);
    } else {
      console.log('command not found.');
    }
    return;
  }

  const { cmd } = await prompt(
    [
      {
        _: true,
        type: 'fuzzy',
        name: 'cmd',
        message: 'what do you want to do?',
        choices: Object.keys(commands)
      }
    ],
    argv
  );
  if (typeof commands[cmd] === 'function') {
    await commands[cmd](argv);
  } else {
    console.log('command not found.');
  }
};
