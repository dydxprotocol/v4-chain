import { IndexerTendermintBlock, IndexerTendermintEvent } from '@dydxprotocol-indexer/v4-protos';

export interface AnnotatedIndexerTendermintEvent extends IndexerTendermintEvent {
  data: string;
}

export interface AnnotatedIndexerTendermintBlock extends IndexerTendermintBlock {
  annotatedEvents: AnnotatedIndexerTendermintEvent[];
}
