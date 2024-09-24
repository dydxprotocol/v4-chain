import crypto from 'crypto';

import { TooManyRequestsError, logger, stats } from '@dydxprotocol-indexer/base';
import axios, { AxiosResponse } from 'axios';
import _ from 'lodash';

import config from '../config';
import { ComplianceClientError } from '../lib/error';
import { ComplianceClientResponse } from '../types';
import { ComplianceClient } from './compliance-client';

export type EllipticPayload = object;

export interface EllipticPostArgs {
  payload: EllipticPayload,
  headers: object,
}

interface ParsedResponse {
  success: boolean,
  riskScore: number | null,
}

export const HOLISTIC: string = 'holistic';
export const API_PATH: string = '/v2/wallet/synchronous';
export const API_URI: string = `https://aml-api.elliptic.co${API_PATH}`;
export const RISK_SCORE_KEY: string = 'risk_score';
export const NO_RULES_TRIGGERED_RISK_SCORE: number = -1;
// We use different negative values of risk score to represent different elliptic response states
export const NOT_IN_BLOCKCHAIN_RISK_SCORE: number = -2;

export class EllipticProviderClient extends ComplianceClient {
  private apiKey: string;
  private apiSecret: string;

  public constructor() {
    super();
    this.apiKey = config.ELLIPTIC_API_KEY;
    this.apiSecret = config.ELLIPTIC_API_SECRET;
  }

  public async getComplianceResponse(address: string): Promise<ComplianceClientResponse> {
    const riskScore: number | null = await this.getRiskScore(address);

    if (riskScore !== null && riskScore >= config.ELLIPTIC_RISK_SCORE_THRESHOLD) {
      return {
        address,
        blocked: true,
        riskScore: riskScore.toFixed(),
      };
    }

    return {
      address,
      blocked: false,
      riskScore: riskScore === null ? undefined : riskScore.toFixed(),
    };
  }

  async getRiskScore(
    address: string,
    retries: number = 0,
  ): Promise<number> {
    const { payload, headers }: EllipticPostArgs = this.getPostArgs(address);
    const start: number = Date.now();

    try {
      const response = await axios.post(API_URI, payload, headers);
      stats.timing(`${config.SERVICE_NAME}.get_elliptic_risk_score_total_time`, Date.now() - start);

      const { success, riskScore } = this.parseApiResponse(response);
      if (!success) {
        logger.error({
          at: 'EllipticProviderClient#getRiskScore',
          message: 'Malformed response from Elliptic',
          response,
        });
        stats.increment(
          `${config.SERVICE_NAME}.get_elliptic_risk_score.status_code`,
          { status: 'malformed' },
        );
        throw new ComplianceClientError('Malformed response');
      }

      stats.increment(
        `${config.SERVICE_NAME}.get_elliptic_risk_score.status_code`,
        { status: '200' },
      );
      if (riskScore === null) {
        return NO_RULES_TRIGGERED_RISK_SCORE;
      }

      return riskScore;
    } catch (error) {
      if (
        error?.response?.status === 404 &&
        error?.response?.data?.name === 'NotInBlockchain'
      ) {
        stats.increment(
          `${config.SERVICE_NAME}.get_elliptic_risk_score.status_code`,
          { status: '404' },
        );
        return NOT_IN_BLOCKCHAIN_RISK_SCORE;
      }

      if (error?.response?.status === 429) {
        stats.increment(
          `${config.SERVICE_NAME}.get_elliptic_risk_score.status_code`,
          { status: '429' },
        );
        throw new TooManyRequestsError('Too many requests');
      }

      if (error?.response?.status === 500 && retries < config.ELLIPTIC_MAX_RETRIES) {
        stats.increment(
          `${config.SERVICE_NAME}.get_elliptic_risk_score.status_code`,
          { status: '500' },
        );
        return this.getRiskScore(address, retries + 1);
      }

      throw error;
    }
  }

  parseApiResponse(response: AxiosResponse): ParsedResponse {
    const riskScore: number | null | undefined = response.data[RISK_SCORE_KEY];

    if (riskScore === null) {
      return {
        success: true,
        riskScore: null,
      };
    }

    if (riskScore === undefined ||
      !Number.isFinite(riskScore) ||
      !_.has(response, 'data.evaluation_detail.source') ||
      !_.has(response, 'data.evaluation_detail.destination')) {
      return {
        success: false,
        riskScore: null,
      };
    }

    return {
      success: true,
      riskScore,
    };
  }

  getPostArgs(
    address: string,
  ): EllipticPostArgs {
    const payload: EllipticPayload = this.getPayload(address);
    const requestTimeMs: number = Date.now();
    const signature: string = this.getApiSignature(requestTimeMs, JSON.stringify(payload));
    const headers: object = this.getHeaders(this.apiKey, signature, requestTimeMs);

    return { payload, headers };
  }

  /*
  * Generate a signature for use when signing a request to the API
  *
  *   - timeOfRequest: current time, in milliseconds, since 1 Jan 1970 00:00:00 UTC
  *   - payload:       string encoded JSON object or '{}' if there is no request body
  *
  * Copied from Elliptic API docs (Removed some unnecessary params)
  * https://app.elliptic.co/#section/Cookbooks/Authentication
  */
  getApiSignature(timeOfRequest: number, payload: string): string {
    // create a SHA256 HMAC using the supplied secret, decoded from base64
    const secret: string = this.apiSecret;
    const hmac = crypto.createHmac('sha256', Buffer.from(secret, 'base64'));

    // concatenate the request text to be signed
    const httpMethod: string = 'POST'; // http method must be uppercase
    const httpPath: string = API_PATH;
    const requestText: string = timeOfRequest + httpMethod + httpPath.toLowerCase() + payload;

    // update the HMAC with the text to be signed
    hmac.update(requestText);

    // output the signature as a base64 encoded string
    return hmac.digest('base64');
  }

  getPayload(
    address: string,
  ): EllipticPayload {
    return {
      subject: {
        asset: HOLISTIC,
        blockchain: HOLISTIC,
        type: 'address',
        hash: address,
      },
      type: 'wallet_exposure',
      customer_reference: 'string',
    };
  }

  getHeaders(
    apiKey: string,
    signature: string,
    requestTimeMs: number,
  ): object {
    return {
      headers: {
        'x-access-key': apiKey,
        'x-access-sign': signature,
        'x-access-timestamp': requestTimeMs,
      },
    };
  }
}
