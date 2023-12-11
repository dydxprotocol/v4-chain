CREATE OR REPLACE FUNCTION dydx_uuid_from_transfer_parts(event_id bytea, asset_id text, sender_subaccount_id uuid, recipient_subaccount_id uuid, sender_wallet_address text, recipient_wallet_address text) RETURNS uuid AS $$
/**
  Returns a UUID using the parts of a transfer.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
DECLARE
    sender_subaccount_id_or_undefined text;
    recipient_subaccount_id_or_undefined text;
    sender_wallet_address_or_undefined text;
    recipient_wallet_address_or_undefined text;
BEGIN
    /** TODO(IND-483): Fix all uuid string substitutions to use Array.join so that we can drop the 'undefined' substitutions below. */
    IF sender_subaccount_id IS NULL THEN
        sender_subaccount_id_or_undefined = 'undefined';
    ELSE
        sender_subaccount_id_or_undefined = sender_subaccount_id;
    END IF;
    IF recipient_subaccount_id IS NULL THEN
        recipient_subaccount_id_or_undefined = 'undefined';
    ELSE
        recipient_subaccount_id_or_undefined = recipient_subaccount_id;
    END IF;
    IF sender_wallet_address IS NULL THEN
        sender_wallet_address_or_undefined = 'undefined';
    ELSE
        sender_wallet_address_or_undefined = sender_wallet_address;
    END IF;
    IF recipient_wallet_address IS NULL THEN
        recipient_wallet_address_or_undefined = 'undefined';
    ELSE
        recipient_wallet_address_or_undefined = recipient_wallet_address;
    END IF;
    return dydx_uuid(concat(sender_subaccount_id_or_undefined, '-', recipient_subaccount_id_or_undefined, '-', sender_wallet_address_or_undefined, '-', recipient_wallet_address_or_undefined, '-', encode(event_id, 'hex'), '-', asset_id));
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
