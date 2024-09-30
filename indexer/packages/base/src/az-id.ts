import { DescribeAvailabilityZonesCommand, EC2Client } from '@aws-sdk/client-ec2';

import { axiosRequest } from './axios';
import config from './config';
import logger from './logger';

export async function getAvailabilityZoneId(): Promise<string> {
  if (config.ECS_CONTAINER_METADATA_URI_V4 !== '' && config.AWS_REGION !== '') {
    const taskUrl = `${config.ECS_CONTAINER_METADATA_URI_V4}/task`;
    try {
      const response = await axiosRequest({
        method: 'GET',
        url: taskUrl,
      }) as { AvailabilityZone: string };
      const client = new EC2Client({ region: config.AWS_REGION });
      const command = new DescribeAvailabilityZonesCommand({
        ZoneNames: [response.AvailabilityZone],
      });
      try {
        const ec2Response = await client.send(command);
        const zoneId = ec2Response.AvailabilityZones![0].ZoneId!;
        logger.info({
          at: 'az-id#getAvailabilityZoneId',
          message: `Got availability zone id ${zoneId}.`,
        });
        return ec2Response.AvailabilityZones![0].ZoneId!;
      } catch (error) {
        logger.error({
          at: 'az-id#getAvailabilityZoneId',
          message: 'Failed to fetch availabilty zone id from EC2. ',
          error,
        });
        return '';
      }
    } catch (error) {
      logger.error({
        at: 'az-id#getAvailabilityZoneId',
        message: 'Failed to retrieve availability zone from metadata endpoint. No availabilty zone id found.',
        error,
        taskUrl,
      });
      return '';
    }
  } else {
    logger.error({
      at: 'az-id#getAvailabilityZoneId',
      message: 'No metadata URI or region. No availabilty zone id found.',
    });
    return '';
  }
}
