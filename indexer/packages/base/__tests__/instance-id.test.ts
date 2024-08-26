import { setInstanceId, getInstanceId, resetForTests } from '../src/instance-id';
import { axiosRequest } from '../src/axios';
import { asMock } from '@dydxprotocol-indexer/dev';
import logger from '../src/logger';
import config from '../src/config';

jest.mock('../src/axios', () => ({
  ...(jest.requireActual('../src/axios') as object),
  axiosRequest: jest.fn(),
}));

describe('instance-id', () => {
  describe('setInstanceId', () => {
    const defaultTaskArn = 'defaultTaskArn';
    const defaultResponse = {
      TaskARN: defaultTaskArn,
    };
    const ecsUrl = config.ECS_CONTAINER_METADATA_URI_V4;

    beforeEach(() => {
      config.ECS_CONTAINER_METADATA_URI_V4 = ecsUrl;
      resetForTests();
      jest.resetAllMocks();
      jest.restoreAllMocks();
      asMock(axiosRequest).mockResolvedValue(defaultResponse);
    });

    afterAll(() => {
      jest.clearAllMocks();
      jest.restoreAllMocks();
    });

    it('should set instance id to task ARN in staging', async () => {
      jest.spyOn(config, 'isStaging').mockReturnValueOnce(true);
      config.ECS_CONTAINER_METADATA_URI_V4 = 'url';
      await setInstanceId();

      expect(getInstanceId()).toEqual(defaultTaskArn);
    });

    it('should set instance id to task ARN in production', async () => {
      jest.spyOn(config, 'isProduction').mockReturnValueOnce(true);
      config.ECS_CONTAINER_METADATA_URI_V4 = 'url';
      await setInstanceId();

      expect(getInstanceId()).toEqual(defaultTaskArn);
    });

    it('should not call metadata endpoint if not production or staging', async () => {
      config.ECS_CONTAINER_METADATA_URI_V4 = 'url';
      await setInstanceId();

      expect(getInstanceId()).not.toEqual(defaultTaskArn);
      expect(asMock(axiosRequest)).not.toHaveBeenCalled();
    });

    it('should not set instance id if already set', async () => {
      jest.spyOn(config, 'isStaging').mockReturnValue(true);
      config.ECS_CONTAINER_METADATA_URI_V4 = 'url';
      await setInstanceId();
      const instanceId = getInstanceId();
      await setInstanceId();

      expect(getInstanceId()).toEqual(instanceId);
      expect(axiosRequest).toHaveBeenCalledTimes(1);
    });

    it('should log error and set instance id to uuid if request errors', async () => {
      jest.spyOn(config, 'isStaging').mockReturnValue(true);
      config.ECS_CONTAINER_METADATA_URI_V4 = 'url';
      const loggerErrorSpy = jest.spyOn(logger, 'error');
      const emptyInstanceId = getInstanceId();
      asMock(axiosRequest).mockRejectedValueOnce(new Error());
      await setInstanceId();

      expect(loggerErrorSpy).toHaveBeenCalledTimes(1);
      expect(getInstanceId()).not.toEqual(emptyInstanceId);
    });

    it('should not call metadata endpoint if url is empty', async () => {
      jest.spyOn(config, 'isStaging').mockReturnValue(true);
      await setInstanceId();

      expect(axiosRequest).not.toHaveBeenCalled();
    });
  });
});
