import { generateUsername } from '../../src/helpers/usernames-helper';

describe('usernames-helper', () => {
  it('Check format of username', () => {
    const username: string = generateUsername();
    expect(username.match(/[A-Z]/g)).toHaveLength(2);
    expect(username.match(/\d/g)).toHaveLength(3);
    // check length is at the very minimum 7
    expect(username.length).toBeGreaterThanOrEqual(7);
  });
});
