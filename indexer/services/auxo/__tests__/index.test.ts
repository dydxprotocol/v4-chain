// handler.test.ts
import { handler } from '../src/index';
import * as helpers from '../src/helpers';
import { InvokeCommandOutput, LambdaClient } from '@aws-sdk/client-lambda';
import { APIGatewayEvent, Context } from 'aws-lambda';
import { AuxoEventJson } from 'src/types';

// Mock logger and startBugsnag from @dydxprotocol-indexer/base
jest.mock('@dydxprotocol-indexer/base', () => {
  const originalModule = jest.requireActual('@dydxprotocol-indexer/base');
  return {
    ...originalModule,
    logger: {
      info: jest.fn(),
      error: jest.fn(),
      crit: jest.fn(),
    },
    startBugsnag: jest.fn(),
  };
});

describe('Auxo Handler', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should return 500 when Bazooka Lambda errors', async () => {
    // Mock upgradeBazooka to do nothing
    jest.spyOn(helpers, 'upgradeBazooka').mockResolvedValue(undefined);

    // Mock LambdaClient.send to return response with FunctionError
    jest.spyOn(LambdaClient.prototype, 'send').mockImplementation(() => {
      return {
        StatusCode: 500,
        FunctionError: 'Some bazooka error',
        $metadata: {
          httpStatusCode: 200,  // api returns 200 even if lambda runtime error
          requestId: 'mock-request-id-invoke',
          extendedRequestId: 'mock-extended-request-id-invoke',
          cfId: 'mock-cf-id-invoke',
          attempts: 1,
          totalRetryDelay: 0,
        },
      } as InvokeCommandOutput;
    });

    const mockEvent: APIGatewayEvent & AuxoEventJson = {
      // APIGatewayEvent properties
      body: null,
      headers: {},
      multiValueHeaders: {},
      httpMethod: 'POST',
      isBase64Encoded: false,
      path: '/deploy',
      pathParameters: null,
      queryStringParameters: null,
      multiValueQueryStringParameters: null,
      stageVariables: null,
      resource: '',
      requestContext: {} as any,
      // AuxoEventJson properties
      upgrade_tag: 'some_tag',
      prefix: 'some_prefix',
      region: 'us-east-1',
      regionAbbrev: 'us-east-1',
      addNewKafkaTopics: false,
      onlyRunDbMigrationAndCreateKafkaTopics: false,
    };

    const mockContext: Context = {} as any;

    await expect(handler(mockEvent, mockContext)).rejects.toThrow('bazooka failure: Some bazooka error');
  });
});
