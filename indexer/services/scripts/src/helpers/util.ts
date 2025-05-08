/* eslint-disable no-console */
export function runAsyncScript(script: () => Promise<void>): void {
  script()
    .then(() => {
      console.log('Done.');
      process.exit(0);
    })
    .catch((error) => {
      if (!process.argv.includes('--verbose-errors')) {
        delete error.originalError; /* eslint-disable-line no-param-reassign */
      }
      console.error(error);
      process.exit(1);
    });
}
