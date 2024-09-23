/**
 * Helper for running an async script.
 */
export function runAsyncScript(script: () => Promise<void>): void {
  script()
    .then(() => {
      // eslint-disable-next-line no-console
      console.log('Done.');
      process.exit(0);
    })
    .catch((error) => {
      if (!process.argv.includes('--verbose-errors')) {
        delete error.originalError; /* eslint-disable-line no-param-reassign */
      }
      // eslint-disable-next-line no-console
      console.error(error);
      process.exit(1);
    });
}
