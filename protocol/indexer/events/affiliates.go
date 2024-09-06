package events

// NewRegisterAffiliateEventV1 creates a RegisterAffiliateEventV1 representing
// a referee being registered with an affiliate.
func NewRegisterAffiliateEventV1(
	referee string,
	affiliate string,
) *RegisterAffiliateEventV1 {
	return &RegisterAffiliateEventV1{
		Referee:   referee,
		Affiliate: affiliate,
	}
}
