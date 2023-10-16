import axios, { AxiosResponse } from 'axios';
import {
  EllipticPostArgs,
  EllipticProviderClient,
  HOLISTIC,
} from '../../src/clients/elliptic-provider';
import { ComplianceClientError } from '../../src/lib/error';
import { TooManyRequestsError } from '@dydxprotocol-indexer/base';
import config from '../../src/config';
import { ComplianceClientResponse } from 'packages/compliance/src';

const defaultAddress: string = 'dydx1f9k5qldwmqrnwy8hcgp4fw6heuvszt35egvtx2';

jest.mock('axios');
describe('elliptic-provider', () => {
  const provider: EllipticProviderClient = new EllipticProviderClient();
  const axiosMock: jest.Mock = (axios.post as unknown as jest.Mock);

  beforeEach(() => {
    jest.useFakeTimers().setSystemTime(new Date(2023, 9, 25, 0, 0, 0, 0));
  });

  afterEach(() => {
    jest.clearAllMocks();
    jest.useRealTimers();
  });

  it('gets correct arguments for POST request', () => {
    const { payload, headers }: EllipticPostArgs = provider.getPostArgs(defaultAddress);
    expect(payload).toEqual({
      customer_reference: 'string',
      subject: {
        asset: HOLISTIC,
        blockchain: HOLISTIC,
        hash: defaultAddress,
        type: 'address',
      },
      type: 'wallet_exposure',
    });
    expect(headers).toEqual({
      headers: {
        'x-access-key': 'default_elliptic_api_key',
        'x-access-sign': 'ihaI6NUlY4QZ3co7fV0jLH4XYCHELKQ/nt1Q5XRN2no=',
        'x-access-timestamp': 1698192000000,
      },
    });
  });

  describe('getRiskScore', () => {
    it('throws error if Elliptic response is malformed', async () => {
      axiosMock.mockResolvedValueOnce(getMockResponse());
      await expect(provider.getRiskScore(defaultAddress))
        .rejects.toEqual(new ComplianceClientError('Malformed response'));
    });

    it('throws error if Elliptic response throws an error', async () => {
      const thrownError: Error = new Error('some error');
      axiosMock.mockRejectedValueOnce(thrownError);
      await expect(provider.getRiskScore(defaultAddress)).rejects.toEqual(thrownError);
    });

    it('throws error if Elliptic response throws TooManyRequesetsError', async () => {
      axiosMock.mockRejectedValue({ response: { status: 429 } });
      await expect(provider.getRiskScore(defaultAddress))
        .rejects.toEqual(new TooManyRequestsError('Too many requests'));
    });

    it('Successfully returns user risk score', async () => {
      axiosMock.mockResolvedValueOnce(getMockResponse(config.ELLIPTIC_RISK_SCORE_THRESHOLD + 1));
      const riskScore: number = await provider.getRiskScore(defaultAddress);
      expect(riskScore).toEqual(config.ELLIPTIC_RISK_SCORE_THRESHOLD + 1);
    });

    it('retries on internal error from Elliptic', async () => {
      axiosMock
        .mockRejectedValueOnce({ response: { status: 500 } })
        .mockResolvedValueOnce(getMockResponse(config.ELLIPTIC_RISK_SCORE_THRESHOLD));
      const riskScore: number = await provider.getRiskScore(defaultAddress);
      expect(riskScore).toEqual(config.ELLIPTIC_RISK_SCORE_THRESHOLD);
      expect(axiosMock).toHaveBeenCalledTimes(2);
    });

    it('throws error on internal error from Elliptic over retry threshold', async () => {
      const internalError: object = { response: { status: 500 } };
      axiosMock
        .mockRejectedValueOnce({ response: { status: 500 } })
        .mockRejectedValueOnce({ response: { status: 500 } })
        .mockRejectedValueOnce({ response: { status: 500 } })
        .mockRejectedValueOnce({ response: { status: 500 } })
        .mockResolvedValueOnce(getMockResponse(config.ELLIPTIC_RISK_SCORE_THRESHOLD));
      await expect(provider.getRiskScore(defaultAddress))
        .rejects.toEqual(internalError);
      expect(axiosMock).toHaveBeenCalledTimes(config.ELLIPTIC_MAX_RETRIES + 1);
    });
  });

  describe('getComplianceResponse', () => {
    it('gets compliance response for blocked user', async () => {
      axiosMock.mockResolvedValueOnce(getMockResponse(config.ELLIPTIC_RISK_SCORE_THRESHOLD));
      const complianceData: ComplianceClientResponse = await provider.getComplianceResponse(
        defaultAddress,
      );
      expect(complianceData).toEqual({
        address: defaultAddress,
        riskScore: (config.ELLIPTIC_RISK_SCORE_THRESHOLD).toFixed(),
        blocked: true,
      });
    });

    it('gets compliance response for non-blocked user', async () => {
      axiosMock.mockResolvedValueOnce(getMockResponse(config.ELLIPTIC_RISK_SCORE_THRESHOLD - 1));
      const complianceData: ComplianceClientResponse = await provider.getComplianceResponse(
        defaultAddress,
      );
      expect(complianceData).toEqual({
        address: defaultAddress,
        riskScore: (config.ELLIPTIC_RISK_SCORE_THRESHOLD - 1).toFixed(),
        blocked: false,
      });
    });

    it('throws error if Elliptic response is an error', async () => {
      const thrownError: Error = new Error('some error');
      axiosMock.mockRejectedValueOnce(thrownError);
      await expect(provider.getComplianceResponse(defaultAddress)).rejects.toEqual(thrownError);
    });
  });
});

function getMockResponse(
  riskScore?: number,
): AxiosResponse {
  // Copied mocked Elliptic response from Elliptic documentation:
  // https://app.elliptic.co/developers/docs#tag/Wallet-Analyses/paths/~1wallet~1synchronous/post
  return {
    status: 0,
    statusText: '',
    headers: {},
    request: null,
    config: {},
    data: {
      id: 'b7535048-76f8-4f60-bdd3-9d659298f9e7',
      type: 'wallet_exposure',
      subject: {
        asset: 'ETH',
        type: 'address',
        hash: '1MdYC22Gmjp2ejVPCxyYjFyWbQCYTGhGq8',
      },
      customer: {
        reference: 'foobar',
      },
      blockchain_info: {
        cluster: {
          inflow_value: {
            usd: 38383838,
          },
          outflow_value: {
            usd: 0,
          },
        },
      },
      created_at: '2015-05-13T10:36:21.000Z',
      updated_at: '2015-05-13T10:36:21.000Z',
      analysed_at: '2015-05-13T10:36:21.000Z',
      cluster_entities: [
        {
          name: 'Mt.Gox',
          category: 'Exchange',
          is_primary_entity: true,
          is_vasp: true,
        },
      ],
      process_status: 'running',
      team_id: 'e333694b-c7c7-4a36-bf35-ed2615865242',
      risk_score: riskScore,
      risk_score_detail: {
        source: 6,
        destination: 6,
      },
      error: {
        message: 'something went wrong',
      },
      evaluation_detail: {
        source: [
          {
            rule_id: 'b7535048-76f8-4f60-bdd3-9d659298f9e7',
            rule_name: 'Illict',
            risk_score: 9.038007,
            matched_elements: [
              {
                category: 'Dark Market',
                contribution_percentage: 0,
                contribution_value: {
                  native: 38383838,
                  native_major: 383.83838,
                  usd: 100.0,
                },
                contributions: [
                  {
                    contribution_percentage: 0,
                    entity: 'AlphaBay',
                    risk_triggers: {
                      name: 'Binance',
                      category: 'Dark Markets',
                      is_sanctioned: true,
                      country: [
                        'MM',
                      ],
                    },
                    contribution_value: {
                      native: 38383838,
                      native_major: 383.83838,
                      usd: 383.83838,
                    },
                  },
                ],
              },
            ],
          },
        ],
        destination: [
          {
            rule_id: 'b7535048-76f8-4f60-bdd3-9d659298f9e7',
            rule_name: 'Illict',
            risk_score: 9.038007,
            matched_elements: [
              {
                category: 'Dark Market',
                contribution_percentage: 0,
                contribution_value: {
                  native: 38383838,
                  native_major: 383.83838,
                  usd: 100.0,
                },
                contributions: [
                  {
                    contribution_percentage: 0,
                    entity: 'AlphaBay',
                    risk_triggers: {
                      name: 'Binance',
                      category: 'Dark Markets',
                      is_sanctioned: true,
                      country: [
                        'MM',
                      ],
                    },
                    contribution_value: {
                      native: 38383838,
                      native_major: 383.83838,
                      usd: 383.83838,
                    },
                  },
                ],
              },
            ],
          },
        ],
      },
      contributions: {
        source: [
          {
            entities: [
              {
                name: 'Alphabay',
                category: 'Dark Market',
                is_primary_entity: true,
                is_vasp: true,
              },
            ],
            contribution_percentage: 6.9883,
            contribution_value: {
              native: 0.07414304,
              native_major: 0.07414304,
              usd: 0.07414304,
            },
          },
        ],
        destination: [
          {
            entities: [
              {
                name: 'Alphabay',
                category: 'Dark Market',
                is_primary_entity: true,
                is_vasp: true,
              },
            ],
            contribution_percentage: 6.9883,
            contribution_value: {
              native: 0.07414304,
              native_major: 0.07414304,
              usd: 0.07414304,
            },
          },
        ],
      },
    },
  };
}
