import { ServiceException as __ServiceException, } from "@smithy/smithy-client";
export { __ServiceException };
export class ECSServiceException extends __ServiceException {
    constructor(options) {
        super(options);
        Object.setPrototypeOf(this, ECSServiceException.prototype);
    }
}
