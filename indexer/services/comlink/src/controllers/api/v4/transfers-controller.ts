import { stats, cacheControlMiddleware } from '@dydxprotocol-indexer/base';
import {
  AssetColumns,
  AssetFromDatabase,
  AssetTable,
  DEFAULT_POSTGRES_OPTIONS,
  IsoString,
  Ordering,
  PaginationFromDatabase,
  QueryableField,
  SubaccountColumns,
  SubaccountFromDatabase,
  SubaccountTable,
  TransferColumns,
  TransferFromDatabase,
  TransferTable,
  USDC_ASSET_ID,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import { matchedData } from 'express-validator';
import _ from 'lodash';
import {
  Controller, Get, Query, Route,
} from 'tsoa';

import { defaultRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { complianceAndGeoCheck } from '../../../lib/compliance-and-geo-check';
import { NotFoundError } from '../../../lib/errors';
import { getChildSubaccountNums, handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import {
  CheckLimitAndCreatedBeforeOrAtSchema,
  CheckPaginationSchema,
  CheckParentSubaccountSchema,
  CheckSubaccountSchema,
  CheckTransferBetweenSchema,
} from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import {
  transferToParentSubaccountResponseObject,
  transferToResponseObject,
} from '../../../request-helpers/request-transformer';
import {
  AssetById,
  SubaccountById,
  TransferRequest,
  TransferResponse,
  ParentSubaccountTransferRequest,
  ParentSubaccountTransferResponse,
  ParentSubaccountTransferResponseObject,
  TransferBetweenRequest,
  TransferBetweenResponse,
} from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'transfers-controller';
const transfersCacheControlMiddleware = cacheControlMiddleware(
  config.CACHE_CONTROL_DIRECTIVE_TRANSFERS,
);

@Route('transfers')
class TransfersController extends Controller {
  @Get('/')
  async getTransfers(
    @Query() address: string,
      @Query() subaccountNumber: number,
      @Query() limit?: number,
      @Query() createdBeforeOrAtHeight?: number,
      @Query() createdBeforeOrAt?: IsoString,
      @Query() page?: number,
  ): Promise<TransferResponse> {
    const subaccountId: string = SubaccountTable.uuid(address, subaccountNumber);

    // TODO(DEC-656): Change to a cache in Redis similar to Librarian instead of querying DB.
    const [subaccount, {
      results: transfers, limit: pageSize, offset, total,
    }, idToAsset]: [
      SubaccountFromDatabase | undefined,
      PaginationFromDatabase<TransferFromDatabase>,
      AssetById,
    ] = await Promise.all([
      SubaccountTable.findById(
        subaccountId,
      ),
      TransferTable.findAllToOrFromSubaccountId(
        {
          subaccountId: [subaccountId],
          limit,
          createdBeforeOrAtHeight: createdBeforeOrAtHeight
            ? createdBeforeOrAtHeight.toString()
            : undefined,
          createdBeforeOrAt,
          page,
        },
        [QueryableField.LIMIT],
        {
          ...DEFAULT_POSTGRES_OPTIONS,
          orderBy: page !== undefined ? [
            [TransferColumns.eventId, Ordering.DESC],
          ]
            : [
              [TransferColumns.createdAtHeight, Ordering.DESC],
            ],
        },
      ),
      getAssetById(),
    ]);
    if (subaccount === undefined) {
      throw new NotFoundError(
        `No subaccount found with address ${address} and subaccountNumber ${subaccountNumber}`,
      );
    }
    const recipientSubaccountIds: string[] = _
      .map(transfers, TransferColumns.recipientSubaccountId)
      .filter(
        (recipientSubaccountId: string | undefined) => recipientSubaccountId !== undefined,
      ) as string[];
    const senderSubaccountIds: string[] = _
      .map(transfers, TransferColumns.senderSubaccountId)
      .filter(
        (senderSubaccountId: string | undefined) => senderSubaccountId !== undefined,
      ) as string[];

    const subaccountIds: string[] = _.uniq([
      ...recipientSubaccountIds,
      ...senderSubaccountIds,
    ]);
    const idToSubaccount: SubaccountById = await idToSubaccountFromSubaccountIds(subaccountIds);

    return {
      transfers: transfers.map((transfer: TransferFromDatabase) => {
        return transferToResponseObject(transfer, idToAsset, idToSubaccount, subaccountId);
      }),
      pageSize,
      totalResults: total,
      offset,
    };
  }

  @Get('/parentSubaccountNumber')
  async getTransfersForParentSubaccount(
    @Query() address: string,
      @Query() parentSubaccountNumber: number,
      @Query() limit?: number,
      @Query() createdBeforeOrAtHeight?: number,
      @Query() createdBeforeOrAt?: IsoString,
      @Query() page?: number,
  ): Promise<ParentSubaccountTransferResponse> {
  // get all child subaccountIds for the parent subaccount number
    const subaccountIds: string[] = getChildSubaccountNums(parentSubaccountNumber).map(
      (childSubaccountNumber: number) => SubaccountTable.uuid(address, childSubaccountNumber),
    );

    // TODO(DEC-656): Change to a cache in Redis similar to Librarian instead of querying DB.
    const [subaccounts, {
      results: transfers,
      limit: pageSize,
      offset,
      total,
    }, idToAsset]: [
      SubaccountFromDatabase[] | undefined,
      PaginationFromDatabase<TransferFromDatabase>,
      AssetById,
    ] = await Promise.all([
      SubaccountTable.findAll(
        { id: subaccountIds },
        [],
      ),
      TransferTable.findAllToOrFromParentSubaccount(
        {
          subaccountId: subaccountIds,
          limit,
          createdBeforeOrAtHeight: createdBeforeOrAtHeight
            ? createdBeforeOrAtHeight.toString()
            : undefined,
          createdBeforeOrAt,
          page,
        },
        [QueryableField.LIMIT],
        {
          ...DEFAULT_POSTGRES_OPTIONS,
          orderBy: page !== undefined ? [
            [TransferColumns.eventId, Ordering.DESC],
          ]
            : [
              [TransferColumns.createdAtHeight, Ordering.DESC],
            ],
        },
      ),
      getAssetById(),
    ]);

    if (subaccounts === undefined || subaccounts.length === 0) {
      throw new NotFoundError(
        `No subaccount found with address ${address} and parentSubaccountNumber ${parentSubaccountNumber}`,
      );
    }

    const recipientSubaccountIds: string[] = _
      .map(transfers, TransferColumns.recipientSubaccountId)
      .filter(
        (recipientSubaccountId: string | undefined) => recipientSubaccountId !== undefined,
      ) as string[];
    const senderSubaccountIds: string[] = _
      .map(transfers, TransferColumns.senderSubaccountId)
      .filter(
        (senderSubaccountId: string | undefined) => senderSubaccountId !== undefined,
      ) as string[];

    const allSubaccountIds: string[] = _.uniq([
      ...recipientSubaccountIds,
      ...senderSubaccountIds,
    ]);
    const idToSubaccount: SubaccountById = await idToSubaccountFromSubaccountIds(allSubaccountIds);

    const transfersResponse: ParentSubaccountTransferResponseObject[] = transfers.map(
      (transfer: TransferFromDatabase) => {
        return transferToParentSubaccountResponseObject(
          transfer,
          idToAsset,
          idToSubaccount,
          address,
          parentSubaccountNumber);
      });

    return {
      transfers: transfersResponse,
      pageSize,
      totalResults: total,
      offset,
    };
  }

  @Get('/between')
  async getTransferBetween(
    @Query() sourceAddress: string,
      @Query() sourceSubaccountNumber: number,
      @Query() recipientAddress: string,
      @Query() recipientSubaccountNumber: number,
      @Query() createdBeforeOrAtHeight?: number,
      @Query() createdBeforeOrAt?: IsoString,
  ): Promise<TransferBetweenResponse> {
    const sourceSubaccountId: string = SubaccountTable.uuid(sourceAddress, sourceSubaccountNumber);
    const recipientSubaccountId: string = SubaccountTable.uuid(
      recipientAddress,
      recipientSubaccountNumber,
    );

    const [
      transfers,
      idToSubaccount,
      idToAsset,
      totalNetTransfers,
    ]: [
      TransferFromDatabase[],
      SubaccountById,
      AssetById,
      string,
    ] = await Promise.all([
      TransferTable.findAll({
        limit: config.API_LIMIT_V4,
        senderSubaccountId: [sourceSubaccountId, recipientSubaccountId],
        recipientSubaccountId: [sourceSubaccountId, recipientSubaccountId],
        assetId: [USDC_ASSET_ID],
        createdBeforeOrAt,
        createdBeforeOrAtHeight: createdBeforeOrAtHeight
          ? createdBeforeOrAtHeight.toString()
          : undefined,
      }, [], { orderBy: [[TransferColumns.createdAtHeight, Ordering.DESC]] }),
      idToSubaccountFromSubaccountIds([sourceSubaccountId, recipientSubaccountId]),
      getAssetById(),
      TransferTable.getNetTransfersBetweenSubaccountIds(
        sourceSubaccountId,
        recipientSubaccountId,
        USDC_ASSET_ID,
      ),
    ]);

    return {
      transfersSubset: transfers.map((transfer: TransferFromDatabase) => {
        return transferToResponseObject(transfer, idToAsset, idToSubaccount, sourceSubaccountId);
      }),
      totalNetTransfers,
    };
  }
}

router.get(
  '/',
  rateLimiterMiddleware(defaultRateLimiter),
  transfersCacheControlMiddleware,
  ...CheckSubaccountSchema,
  ...CheckLimitAndCreatedBeforeOrAtSchema,
  ...CheckPaginationSchema,
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
      subaccountNumber,
      limit,
      createdBeforeOrAtHeight,
      createdBeforeOrAt,
      page,
    }: TransferRequest = matchedData(req) as TransferRequest;

    try {
      const controllers: TransfersController = new TransfersController();
      const response: TransferResponse = await controllers.getTransfers(
        address,
        subaccountNumber,
        limit,
        createdBeforeOrAtHeight,
        createdBeforeOrAt,
        page,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'TransfersController GET /',
        'Transfers error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_transfers.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/parentSubaccountNumber',
  rateLimiterMiddleware(defaultRateLimiter),
  transfersCacheControlMiddleware,
  ...CheckParentSubaccountSchema,
  ...CheckLimitAndCreatedBeforeOrAtSchema,
  ...CheckPaginationSchema,
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
      parentSubaccountNumber,
      limit,
      createdBeforeOrAtHeight,
      createdBeforeOrAt,
      page,
    }: ParentSubaccountTransferRequest = matchedData(req) as ParentSubaccountTransferRequest;

    // The schema checks allow subaccountNumber to be a string, but we know it's a number here.
    const parentSubaccountNum: number = +parentSubaccountNumber;

    try {
      const controllers: TransfersController = new TransfersController();
      const response: TransferResponse = await controllers.getTransfersForParentSubaccount(
        address,
        parentSubaccountNum,
        limit,
        createdBeforeOrAtHeight,
        createdBeforeOrAt,
        page,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'TransfersController GET /',
        'Transfers error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_transfers.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/between',
  rateLimiterMiddleware(defaultRateLimiter),
  transfersCacheControlMiddleware,
  ...CheckTransferBetweenSchema,
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      sourceAddress,
      sourceSubaccountNumber,
      recipientAddress,
      recipientSubaccountNumber,
      createdBeforeOrAtHeight,
      createdBeforeOrAt,
    }: TransferBetweenRequest = matchedData(req) as TransferBetweenRequest;

    try {
      const controllers: TransfersController = new TransfersController();
      const response: TransferBetweenResponse = await controllers.getTransferBetween(
        sourceAddress,
        sourceSubaccountNumber,
        recipientAddress,
        recipientSubaccountNumber,
        createdBeforeOrAtHeight,
        createdBeforeOrAt,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'TransfersController GET /between',
        'Transfers error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_transfers_between.timing`,
        Date.now() - start,
      );
    }
  },
);

async function getAssetById(): Promise<AssetById> {
  const assets: AssetFromDatabase[] = await AssetTable.findAll({}, []);
  return _.keyBy(assets, AssetColumns.id);
}

async function idToSubaccountFromSubaccountIds(
  subaccountIds: string[],
): Promise<SubaccountById> {
  const subaccounts: SubaccountFromDatabase[] = await SubaccountTable.findAll(
    {
      id: subaccountIds,
    },
    [],
  );

  const idToSubaccount: SubaccountById = _.keyBy(
    subaccounts,
    SubaccountColumns.id,
  );
  return idToSubaccount;
}

export default router;
