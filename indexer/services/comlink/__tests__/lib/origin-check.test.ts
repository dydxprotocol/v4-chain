import { originAllowlistMiddleware } from '../../src/lib/origin-check';

describe('originAllowlistMiddleware', () => {
  let req: any;
  let res: any;
  let next: jest.Mock;

  beforeEach(() => {
    req = {
      get: jest.fn(),
    };
    res = {
      status: jest.fn().mockReturnThis(),
      json: jest.fn().mockReturnThis(),
      vary: jest.fn().mockReturnThis(),
    };
    next = jest.fn();
  });

  it('is a no-op when the allowlist is empty', () => {
    const middleware = originAllowlistMiddleware('');
    req.get.mockReturnValue('https://evil.com');

    middleware(req, res, next);

    expect(next).toHaveBeenCalledTimes(1);
    expect(res.status).not.toHaveBeenCalled();
    expect(res.vary).not.toHaveBeenCalled();
  });

  it('is a no-op when the allowlist contains only whitespace entries', () => {
    const middleware = originAllowlistMiddleware(' , , ');
    req.get.mockReturnValue('https://evil.com');

    middleware(req, res, next);

    expect(next).toHaveBeenCalledTimes(1);
    expect(res.status).not.toHaveBeenCalled();
  });

  it('allows a request whose Origin is in the allowlist', () => {
    const middleware = originAllowlistMiddleware(
      'https://dydx.trade,https://staging.dydx.trade',
    );
    req.get.mockReturnValue('https://dydx.trade');

    middleware(req, res, next);

    expect(next).toHaveBeenCalledTimes(1);
    expect(res.status).not.toHaveBeenCalled();
    expect(res.vary).toHaveBeenCalledWith('Origin');
  });

  it('allows a request whose Origin matches case-insensitively', () => {
    const middleware = originAllowlistMiddleware('https://dydx.trade');
    req.get.mockReturnValue('HTTPS://DYDX.TRADE');

    middleware(req, res, next);

    expect(next).toHaveBeenCalledTimes(1);
    expect(res.status).not.toHaveBeenCalled();
  });

  it('rejects a request whose Origin is not in the allowlist', () => {
    const middleware = originAllowlistMiddleware('https://dydx.trade');
    req.get.mockReturnValue('https://evil.com');

    middleware(req, res, next);

    expect(next).not.toHaveBeenCalled();
    expect(res.status).toHaveBeenCalledWith(403);
    expect(res.json).toHaveBeenCalledWith({
      errors: [{ msg: 'Origin not allowed' }],
    });
    expect(res.vary).toHaveBeenCalledWith('Origin');
  });

  it('rejects an Origin that only differs by a suffix (no prefix matching)', () => {
    const middleware = originAllowlistMiddleware('https://dydx.trade');
    req.get.mockReturnValue('https://dydx.trade.evil.com');

    middleware(req, res, next);

    expect(next).not.toHaveBeenCalled();
    expect(res.status).toHaveBeenCalledWith(403);
  });

  it('allows requests without an Origin header (mobile / server-to-server)', () => {
    const middleware = originAllowlistMiddleware('https://dydx.trade');
    req.get.mockReturnValue(undefined);

    middleware(req, res, next);

    expect(next).toHaveBeenCalledTimes(1);
    expect(res.status).not.toHaveBeenCalled();
  });

  it('trims whitespace around allowlist entries', () => {
    const middleware = originAllowlistMiddleware(
      '  https://dydx.trade  ,  https://staging.dydx.trade  ',
    );
    req.get.mockReturnValue('https://staging.dydx.trade');

    middleware(req, res, next);

    expect(next).toHaveBeenCalledTimes(1);
    expect(res.status).not.toHaveBeenCalled();
  });
});
