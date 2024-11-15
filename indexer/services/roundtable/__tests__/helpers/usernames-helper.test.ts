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
      'dydx1q54yvrslnu0xp4drpde6f4e0k2ap9efss5hpsd',
      'dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4',
      'dydx1n88uc38xhjgxzw9nwre4ep2c8ga4fjxc565lnf',
      'dydx1n88uc38xhjgxzw9nwre4ep2c8ga4fjxc575lnf',
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
      'RudeFuel59E',
      'MacroMealFK5',
      'HappySnapWTT',
      'BumpyEdgeH5Y',
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
