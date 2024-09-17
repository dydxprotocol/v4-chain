const AZ_LIST = ['apne1-az1', 'apne1-az2', 'apne1-az4'];
const AZ_ID = AZ_LIST[Math.floor(Math.random() * AZ_LIST.length)];

export function getAvailabilityZoneId(): string {
  return AZ_ID;
}
