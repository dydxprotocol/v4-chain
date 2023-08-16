import axios, { AxiosRequestConfig } from 'axios';
import _ from 'lodash';

import {
  AxiosServerError,
  AxiosError,
  AxiosSafeError,
  AxiosSafeServerError,
  AxiosSafeErrorObject,
} from './errors';

axios.defaults.timeout = 10000; // ms

export async function axiosRequest(options: AxiosRequestConfig): Promise<unknown> {
  try {
    const response = await axios(options);
    return response.data;
  } catch (error) {
    if (error.isAxiosError) {
      if (error.response) {
        error.response = _.pick(error.response, ['data', 'status', 'statusText']);
        throw new AxiosServerError(error.response, error);
      }
      // request error or timeout error
      throw new AxiosError(`Axios: ${error.message}`, error);
    }
    throw error;
  }
}

export async function safeAxiosRequest(options: AxiosRequestConfig): Promise<unknown> {
  try {
    const response = await axios(options);
    return response.data;
  } catch (error) {
    if (axios.isAxiosError(error)) {
      const safeError: AxiosSafeErrorObject = _.omit(error.toJSON(), 'config');
      if (error.response) {
        const smallerResponse: { data: {}, status: number, statusText: string } = _.pick(
          error.response,
          ['data', 'status', 'statusText'],
        );
        throw new AxiosSafeServerError(smallerResponse, safeError);
      }
      // request error or timeout error
      throw new AxiosSafeError(`Axios: ${error.message}`, safeError);
    }
    throw error;
  }
}
