import {
  generateUsernameForSubaccount,
} from '../../src/helpers/usernames-helper';

describe('usernames-helper', () => {
  it('Check result and determinism of username username', () => {
    const addresses = [
      'dydx1gf4xlnpulkyex74asxxhg9ye05r28cxdd69s9u',
      'dydx10fx7sy6ywd5senxae9dwytf8jxek3t2gcen2vs',
      'dydx1t72ww7qzdx5rjlpp6cq0cqy09qlsjj7e4kpuyt',
      'dydx1wau5mja7j7zdavtfq9lu7ejef05hm6ffenlcsn',
      'dydx168pjt8rkru35239fsqvz7rzgeclakp49zx3aum',
      'dydx1df84hz7y0dd3mrqcv3vrhw9wdttelul8edqmvp',
      'dydx16h7p7f4dysrgtzptxx2gtpt5d8t834g9dj830z',
      'dydx15u9tppy5e2pdndvlrvafxqhuurj9mnpdstzj6z',
    ];

    const expectedUsernames = [
      'CushyHandVE6',
      'AmpleCubeLKI',
      'AwareFoodHGP',
      'LoudLandXWV',
      'MossyStraw2JJ',
      'BoldGapOGY',
      'ZoomEraQE0',
      'WiryFernLEC',
    ];

    for (let i = 0; i < addresses.length; i++) {
      const address = addresses[i];
      for (let j = 0; j < 10; j++) {
        const names = new Set();
        for (let k = 0; k < 10; k++) {
          const username: string = generateUsernameForSubaccount(address, 0, k);
          if (k === 0) {
            expect(username).toEqual(expectedUsernames[i]);
          }
          names.add(username);
        }
        // for same address, difference nonce should result in different username
        expect(names.size).toEqual(10);
      }
    }
  });
});
