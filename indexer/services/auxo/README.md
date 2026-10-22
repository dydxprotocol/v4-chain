# Auxo

Auxo is an AWS Lambda function that orchestrates automated deployments and upgrades of the dYdX indexer services running on ECS (Elastic Container Service). It handles the complete upgrade lifecycle including database migrations, Kafka topic creation, and rolling updates of all indexer microservices.

## Responsibilities & Scope

- Upgrade the Bazooka Lambda function to a new Docker image version
- Execute database migrations via Bazooka
- Create new Kafka topics when required
- Register new ECS task definitions with updated container images
- Perform rolling updates of ECS services (Comlink, Ender, Roundtable, Socks, Vulcan)
- Coordinate service shutdown and restart to ensure safe database migrations
- Verify ECR image availability before attempting upgrades

**Out of scope:**
- Does not perform application-level health checks or validation
- Does not handle rollback logic
- Does not manage infrastructure provisioning

## Architecture & Dependencies

### Internal Structure

- `src/index.ts` – Main Lambda handler orchestrating the upgrade workflow
- `src/config.ts` – Configuration schema and environment variable parsing
- `src/constants.ts` – Service names, payloads, and static configuration
- `src/types.ts` – TypeScript type definitions for events and mappings

### Processing Flow

1. Receive upgrade event with image tag, region, and environment prefix
2. Upgrade Bazooka Lambda function with new image from ECR
3. Stop database writer services (Ender, Roundtable) to prevent write conflicts
4. Invoke Bazooka to run database migrations and optionally create Kafka topics
5. For each ECS service:
   - Fetch current task definition
   - Verify new image exists in ECR
   - Register new task definition with updated image tag
   - Wait for task definition to become available
6. Update all ECS services to use new task definitions
7. Restart database writer services with original task counts

### Internal Dependencies

- `@dydxprotocol-indexer/base` – Shared logging, configuration parsing, and utilities

### External Dependencies

- **AWS Lambda** – Runtime environment
- **AWS ECR** – Container image registry
- **AWS ECS** – Container orchestration for indexer services
- **Bazooka Lambda** – Executes database migrations and Kafka topic management

## Public Interface

### Lambda Invocation

Auxo is invoked via AWS Lambda with an API Gateway event containing:

**Event Schema:**
```json
{

upgrade_tag: string,   // Docker image tag to deploy

prefix: string,        // Environment prefix (e.g., "dev4", "prod")

region: string,        // AWS region (e.g., "ap-northeast-1")

regionAbbrev: string,  // Abbreviated region name (e.g., "apne1")

addNewKafkaTopics: boolean,       // Whether to create new Kafka topics

onlyRunDbMigrationAndCreateKafkaTopics: boolean  // Skip service upgrades

}
```
**Handler:** [`src/index.ts#handler`](./src/index.ts)
**Response:**

- HTTP 200 with `{ message: 'success' }` on successful completion

- Throws exception on failure

## Configuration
### Environment Variables
- `MAX_TASK_DEFINITION_WAIT_TIME_MS` (number, default: `60000`) – Maximum time in milliseconds to wait for ECS task definition registration

- `SLEEP_TIME_MS` (number, default: `5000`) – Polling interval in milliseconds when waiting for task definitions
Configuration is defined in [`src/config.ts`](./src/config.ts) using the `@dydxprotocol-indexer/base` configuration schema parser.

### Constants
- `BAZOOKA_LAMBDA_FUNCTION_NAME` – Name of the Bazooka Lambda function

- `ECS_SERVICE_NAMES` – List of all ECS services to upgrade: Comlink, Ender, Roundtable, Socks, Vulcan

- `ECS_DB_WRITER_SERVICE_NAMES` – Subset of services that write to the database: Ender, Roundtable

Defined in [`src/constants.ts`](./src/constants.ts).

## Running Locally
Auxo is designed to run as an AWS Lambda function and cannot be easily run locally outside of the AWS environment. For development and testing:
1. Build the service:

   ```bash
   pnpm run build
   ```
2. Run tests:
   ```bash
   pnpm test
   ```

3. Deploy to AWS Lambda:
   - Package the built code and dependencies
   - Deploy using AWS SAM, Terraform, or the AWS Console
   - Configure the Lambda with appropriate IAM permissions for ECR, ECS, and Lambda operations

## Required AWS Permissions:

 - ecr:DescribeImages
 - ecs:DescribeServices
 - ecs:DescribeTaskDefinition
 - ecs:RegisterTaskDefinition
 - ecs:UpdateService
 - lambda:GetFunction
 - lambda:UpdateFunctionCode
 - lambda:InvokeFunction

## Testing

### Unit Tests

Run the test suite:
```bash
pnpm test
```

Run tests with coverage:
```bash
pnpm coverage
```

Tests are located in `__tests__/` and use Jest as the test runner.

Test Configuration:

 - Jest config: jest.config.js
 - Global setup: jest.globalSetup.js
 - Test setup: jest.setup.js
 - Environment: .env.test loads test-specific configuration

Note: Current test coverage is minimal (placeholder test only). Integration tests would require AWS service mocking.
Observability & Operations

## Logging

All operations are logged using the @dydxprotocol-indexer/base logger with structured context:

 - at – Function/location identifier
 - message – Human-readable description
 - Additional context fields (service names, task definitions, responses, errors)

### Key Log Points

 - Bazooka upgrade start and completion
 - Database migration execution
 - ECR image verification
 - ECS task definition registration
 - Service update operations
 - Database writer service stop/start events

## Error Handling

 - Lambda throws exceptions on failure, triggering AWS Lambda error handling
 - Database writer services are restarted in a finally block to ensure recovery even on failure
 - Timeouts are enforced for task definition registration (default 60s)

## Known Limitations

 - No automatic rollback on failure
 - Task definition registration timeout is fixed (not configurable per-service)
 - Assumes all services use the same image tag convention
 - No validation of service health after upgrade

## Deployment & Runtime

### Deployment

Auxo is deployed as an AWS Lambda function. Deployment artifacts and infrastructure are managed outside this service directory.

Runtime Characteristics:

 - Execution Model: Synchronous Lambda invocation via API Gateway
 - Timeout: Must be configured to allow sufficient time for all upgrade operations (typically 5-15 minutes)
 - Memory: Minimal memory requirements (default Lambda settings sufficient)
 - Concurrency: Should be limited to 1 to prevent concurrent upgrades

### IAM Role Requirements

The Lambda execution role must have permissions to:

 - Read from ECR repositories
 - Describe and update ECS services and task definitions
 - Invoke and update the Bazooka Lambda function

## Directory Layout

 - src/ – Source code
    - index.ts – Main Lambda handler and orchestration logic
    - config.ts – Environment variable configuration schema
    - constants.ts – Service names and static configuration
    - types.ts – TypeScript type definitions
 - __tests__/ – Jest test files
 - patches/ – npm package patches (currently empty)
 - build/ – Compiled JavaScript output (generated by pnpm build)
 - package.json – npm dependencies and scripts
 - tsconfig.json – TypeScript compiler configuration
 - .eslintrc.js – ESLint configuration
 - jest.config.js – Jest test configuration

## Related Documents

 - Indexer README – Overview of the entire indexer system
 - Bazooka Service – Database migration and Kafka topic management service
 - ECS Service Configurations – Individual service directories (comlink, ender, roundtable, socks, vulcan)
