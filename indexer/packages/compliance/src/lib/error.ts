export class ComplianceClientError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'ComplianceClientError';
  }
}
