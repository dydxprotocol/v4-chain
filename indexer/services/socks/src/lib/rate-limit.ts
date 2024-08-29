import config from '../config';

// General purpose, synchronous rate limiter for message handlers.
export class RateLimiter {
  private messageInfo: {
    [connectionId: string]: {
      [key: string]: {
        points: number,
        startMs: number | null,
      },
    },
  };

  private points: number;
  private durationMs: number;

  constructor({
    points,
    durationMs,
  }: {
    points: number,
    durationMs: number,
  }) {
    this.messageInfo = {};
    this.points = points;
    this.durationMs = durationMs;
  }

  public rateLimit({
    connectionId,
    key,
  }: {
    connectionId: string,
    key: string,
  }): number {
    if (!config.RATE_LIMIT_ENABLED) {
      return 0;
    }

    if (!this.messageInfo[connectionId]) {
      this.messageInfo[connectionId] = {};
    }

    if (!this.messageInfo[connectionId][key]) {
      this.messageInfo[connectionId][key] = {
        points: 0,
        startMs: null,
      };
    }

    const startMs: number | null = this.messageInfo[connectionId][key].startMs;
    const now: number = Date.now();
    if (
      startMs !== null &&
      (now - startMs) < this.durationMs
    ) {
      // additional traffic in this rate limiting session
      this.messageInfo[connectionId][key].points += 1;
    } else {
      // reset rate limit clock
      this.messageInfo[connectionId][key].startMs = now;
      this.messageInfo[connectionId][key].points = 1;
    }

    if (
      startMs !== null &&
      this.messageInfo[connectionId][key].points > this.points
    ) {
      return (this.durationMs - (now - startMs));
    }
    return 0;
  }

  public removeConnection(connectionId: string) {
    if (this.messageInfo[connectionId]) {
      delete this.messageInfo[connectionId];
    }
  }
}
