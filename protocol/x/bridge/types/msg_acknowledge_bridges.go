package types

func (msg *MsgAcknowledgeBridges) ValidateBasic() error {
	// Validates that bridge event IDs are consecutive.
	for i, event := range msg.Events {
		if i > 0 && msg.Events[i-1].Id != event.Id-1 {
			return ErrBridgeIdsNotConsecutive
		}
	}
	return nil
}
