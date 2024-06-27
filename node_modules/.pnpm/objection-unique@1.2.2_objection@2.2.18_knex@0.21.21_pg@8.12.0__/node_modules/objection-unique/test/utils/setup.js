
/**
 * Module dependencies.
 */

import { Model } from 'objection';
import knex from './knex';

// Bind knex to objection model.
Model.knex(knex);
