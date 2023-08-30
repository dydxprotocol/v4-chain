import express from 'express';
import { ResponseWithBody } from 'src/types';

// Captures the response body and puts it in res.body
export default (_req: express.Request, res: express.Response, next: express.NextFunction) => {
  const oldWrite: Function = res.write;
  const oldEnd: Function = res.end;

  const chunks: Buffer[] = [];

  // eslint-disable-next-line  @typescript-eslint/no-explicit-any
  res.write = (...args: any[]) => {
    chunks.push(Buffer.from(args[0]));
    return oldWrite.apply(res, args);
  };

  // eslint-disable-next-line  @typescript-eslint/no-explicit-any
  res.end = (...args: any[]) => {
    if (args.length && args[0]) {
      chunks.push(Buffer.from(args[0]));
    }

    const body = Buffer.concat(chunks).toString('utf8');
    (res as ResponseWithBody).body = body;
    return oldEnd.apply(res, args);
  };

  return next();
};
