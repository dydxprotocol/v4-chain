import config from './config';
import { STATS_FUNCTION_NAME } from './constants';
import stats from './stats';

export async function runFuncWithTimingStat(
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  promise: Promise<any>,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  options: any,
  functionName?: string,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
): Promise<any> {
  const start: number = Date.now();
  let result;
  try {
    result = await promise;
  } catch (error) {
    stats.timing(
      `${config.SERVICE_NAME}.${functionName ?? STATS_FUNCTION_NAME}.timing`,
      Date.now() - start,
      options,
    );
    throw error;
  }
  stats.timing(
    `${config.SERVICE_NAME}.${functionName ?? STATS_FUNCTION_NAME}.timing`,
    Date.now() - start,
    options,
  );
  return result;
}
