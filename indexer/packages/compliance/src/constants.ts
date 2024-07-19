export const INDEXER_GEOBLOCKED_PAYLOAD: string = 'Because you appear to be a resident of, or trading from, a jurisdiction that violates our terms of use, or have engaged in activity that violates our terms of use, you have been blocked. You may withdraw your funds from the protocol at any time.';
export const INDEXER_COMPLIANCE_BLOCKED_PAYLOAD: string = 'Because this address appears to be a resident of, or trading from, a jurisdiction that violates our terms of use, or has engaged in activity that violates our terms of use, this address has been blocked.';

// For use by other services packages, can't be used to index on the actual requests
// object as that needs to be a string literal.
export const COUNTRY_HEADER_KEY: string = 'cf-ipcountry';
