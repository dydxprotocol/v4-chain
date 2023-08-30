/* eslint-disable @typescript-eslint/no-explicit-any */ // Disabled for file due to difficult anys
import express from 'express';

// Captures the response body and puts it in res.body
export default (
  _req: express.Request,
  res: any,
  next: express.NextFunction,
) => {
  const oldWrite = res.write;
  const oldEnd = res.end;

  const chunks: Buffer[] = [];

  res.write = (...args: any[]) => {
    chunks.push(Buffer.from(args[0]));
    oldWrite.apply(res, args);
  };

  res.end = (...args: any[]) => {
    if (args.length && args[0]) {
      chunks.push(Buffer.from(args[0]));
    }

    const body: string = Buffer.concat(chunks).toString('utf8');
    res.body = body;
    oldEnd.apply(res, args);
  };

  return next();
};
