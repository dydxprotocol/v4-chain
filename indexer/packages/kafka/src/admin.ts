import { Admin } from 'kafkajs';

import { kafka } from './kafka';

export const admin: Admin = kafka.admin();
