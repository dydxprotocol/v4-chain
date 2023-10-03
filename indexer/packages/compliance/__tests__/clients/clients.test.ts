import config from '../../src/config';
import { getComplianceClient } from '../../src/clients/clients';
import { EllipticProviderClient } from '../../src/clients/elliptic-provider';

describe('getComplianceClient', () => {
  const defaultClient: string = config.COMPLIANCE_DATA_CLIENT;

  afterAll(() => {
    config.COMPLIANCE_DATA_CLIENT = defaultClient;
  });

  it('uses Elliptic provider if invalid value in config', () => {
    config.COMPLIANCE_DATA_CLIENT = 'invalid';

    expect(getComplianceClient()).toEqual(new EllipticProviderClient());
  });
});
