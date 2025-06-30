/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type ProofParamsLightAttributes struct {
	BirthDateLowerBound       string  `json:"birth_date_lower_bound"`
	BirthDateUpperBound       string  `json:"birth_date_upper_bound"`
	CallbackUrl               *string `json:"callback_url,omitempty"`
	CitizenshipMask           string  `json:"citizenship_mask"`
	EventData                 string  `json:"event_data"`
	EventId                   string  `json:"event_id"`
	ExpirationDateLowerBound  string  `json:"expiration_date_lower_bound"`
	ExpirationDateUpperBound  string  `json:"expiration_date_upper_bound"`
	IdentityCounter           int64   `json:"identity_counter"`
	IdentityCounterLowerBound int64   `json:"identity_counter_lower_bound"`
	IdentityCounterUpperBound int64   `json:"identity_counter_upper_bound"`
	Selector                  string  `json:"selector"`
	TimestampLowerBound       string  `json:"timestamp_lower_bound"`
	TimestampUpperBound       string  `json:"timestamp_upper_bound"`
}
