package types

func (dm DelayedMessage) Validate() error {
	if dm.Msg == nil || len(dm.Msg) == 0 {
		return ErrMsgIsNil
	}
	return nil
}
